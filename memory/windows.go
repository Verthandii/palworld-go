//go:build windows

package memory

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

//go:embed RAMMap64.exe
var rammapFS embed.FS

type cleaner struct {
	c          *config.Config
	rammapFile string
}

func NewCleaner(c *config.Config) Cleaner {
	cleaner := &cleaner{
		c: c,
	}

	if c.MemoryCheckInterval > 0 {
		rammapFile, err := extractRAMMap()
		if err != nil {
			log.Fatalf("无法提取RAMMap可执行文件: %v", err)
		}
		cleaner.rammapFile = rammapFile
	}

	return cleaner
}

func (cleaner *cleaner) Schedule(ctx context.Context) {
	if cleaner.c.MemoryCheckInterval > 0 && cleaner.c.MemoryCleanupInterval > 0 {
		cleaner.scheduleAll(ctx)
	} else if cleaner.c.MemoryCheckInterval > 0 {
		cleaner.scheduleRebootClean(ctx)
	} else if cleaner.c.MemoryCleanupInterval > 0 {
		cleaner.scheduleClean(ctx)
	}
}

func (cleaner *cleaner) Stop() {
	if cleaner.c.MemoryCheckInterval > 0 {
		_ = os.Remove(cleaner.rammapFile)
	}
}

func (cleaner *cleaner) scheduleAll(ctx context.Context) {
	rebootCleanDuration := time.Duration(cleaner.c.MemoryCheckInterval) * time.Second
	cleanDuration := time.Duration(cleaner.c.MemoryCleanupInterval) * time.Second
	rebootCleanTicker := time.NewTicker(rebootCleanDuration)
	cleanTicker := time.NewTicker(rebootCleanDuration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-rebootCleanTicker.C:
			cleaner.rebootClean()
			rebootCleanTicker.Reset(rebootCleanDuration)
		case <-cleanTicker.C:
			cleaner.clean()
			cleanTicker.Reset(cleanDuration)
		}
	}
}

func (cleaner *cleaner) scheduleRebootClean(ctx context.Context) {
	duration := time.Duration(cleaner.c.MemoryCheckInterval) * time.Second
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleaner.rebootClean()
			ticker.Reset(duration)
		}
	}
}

func (cleaner *cleaner) scheduleClean(ctx context.Context) {
	duration := time.Duration(cleaner.c.MemoryCleanupInterval) * time.Second
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleaner.clean()
			ticker.Reset(duration)
		}
	}
}

// rebootClean 当内存超于阈值时，重启进程以清理内存
func (cleaner *cleaner) rebootClean() {
	// MEMORYSTATUSEX 结构体用于接收全局内存状态信息
	type MEMORYSTATUSEX struct {
		dwLength                uint32
		dwMemoryLoad            uint32
		ullTotalPhys            uint64 // 总物理内存大小
		ullAvailPhys            uint64 // 可用物理内存大小
		ullTotalPageFile        uint64
		ullAvailPageFile        uint64
		ullTotalVirtual         uint64
		ullAvailVirtual         uint64
		ullAvailExtendedVirtual uint64
	}

	cfg := cleaner.c
	threshold := cfg.MemoryUsageThreshold
	memStatus := MEMORYSTATUSEX{dwLength: uint32(unsafe.Sizeof(MEMORYSTATUSEX{}))}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")
	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		log.Println("调用 GlobalMemoryStatusEx 失败")
		return
	}

	total := memStatus.ullTotalPhys / 1024 / 1024 // MB
	free := memStatus.ullAvailPhys / 1024 / 1024  // MB
	memoryUsage := 100.0 * (1 - float64(free)/float64(total))

	if memoryUsage > threshold {
		log.Printf("内存占用超过 %v, 开始清理内存...\n", threshold)
		c, err := rcon.New(cfg.Address, cfg.AdminPassword)
		if err != nil {
			log.Printf("rcon 客户端启动失败 【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		defer c.Close()
	}
}

// clean 使用 RAMMap 清理无用内存
func (cleaner *cleaner) clean() {
	// TODO 打印清理前后的内存
	log.Printf("正在使用 rammap 清理内存....\n")
	cmd := exec.Command(cleaner.rammapFile, "-Ew")
	err := cmd.Run()
	if err != nil {
		log.Printf("运行 RAMMap 时发生错误 【%v】\n", err)
	}
}

// extractRAMMap 从嵌入的文件系统中提取 RAMMap 并写入临时文件
func extractRAMMap() (string, error) {
	rammapData, err := fs.ReadFile(rammapFS, "RAMMap64.exe")
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "RAMMap64-*.exe")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = tmpFile.Write(rammapData); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
