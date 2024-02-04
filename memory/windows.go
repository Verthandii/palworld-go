//go:build windows

package memory

import (
	"embed"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

//go:embed RAMMap64.exe
var rammapFS embed.FS

var rammapFile string

func NewCleaner(c *config.Config, ch chan<- time.Duration) Cleaner {
	cleaner := &cleaner{
		c:  c,
		ch: ch,
	}

	var err error
	if c.MemoryCleanupInterval > 0 {
		rammapFile, err = extractRAMMap()
		if err != nil {
			log.Printf("无法提取 RAMMap 可执行文件【%v】\n", err)
			os.Exit(1)
		}
	}

	return cleaner
}

func (cleaner *cleaner) Stop() {
	if cleaner.c.MemoryCleanupInterval > 0 {
		_ = os.Remove(rammapFile)
	}
}

// rebootClean 当内存超于阈值时，重启进程以清理内存
func (cleaner *cleaner) rebootClean() {
	total, free := cleaner.getMemoryInfo()

	memoryUsage := 100.0 * (1 - float64(free)/float64(total))
	threshold := cleaner.c.MemoryUsageThreshold

	if memoryUsage > threshold {
		log.Printf("内存占用超过【%v】%%, 重新启动游戏服务器\n", threshold)
		c, err := rcon.New(cleaner.c)
		if err != nil {
			log.Printf("RCON 客户端启动失败【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		c.Close()
		cleaner.ch <- 70 * time.Second
	}
}

// clean 使用 RAMMap 清理无用内存
func (cleaner *cleaner) clean() {
	_, free := cleaner.getMemoryInfo()
	log.Printf("空闲内存【%d】MB, 正在清理内存....\n", free)
	cmd := exec.Command(rammapFile, "-Ew")
	err := cmd.Run()
	if err != nil {
		log.Printf("运行 RAMMap 时发生错误 【%v】\n", err)
		if strings.Contains(err.Error(), "The requested operation requires elevation") {
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
		}
	}
	_, free = cleaner.getMemoryInfo()
	log.Printf("清理内存成功, 空闲内存【%d】MB\n", free)
}

func (cleaner *cleaner) getMemoryInfo() (total, free uint64) {
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

	memStatus := MEMORYSTATUSEX{dwLength: uint32(unsafe.Sizeof(MEMORYSTATUSEX{}))}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")
	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		log.Printf("获取内存信息失败【%v】\n", err)
		return
	}

	total = memStatus.ullTotalPhys / 1024 / 1024 // MB
	free = memStatus.ullAvailPhys / 1024 / 1024  // MB

	return
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
