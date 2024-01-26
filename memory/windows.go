//go:build windows

package memory

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func (cleaner *cleaner) Stop() {
	if cleaner.c.MemoryCheckInterval > 0 {
		_ = os.Remove(cleaner.rammapFile)
	}
}

// rebootClean 当内存超于阈值时，重启进程以清理内存
func (cleaner *cleaner) rebootClean() {
	total, free := cleaner.getMemoryInfo()

	memoryUsage := 100.0 * (1 - float64(free)/float64(total))
	threshold := cleaner.c.MemoryUsageThreshold

	if memoryUsage > threshold {
		log.Printf("【Memory】内存占用超过【%v】, 重新启动游戏服务器\n", threshold)
		c, err := rcon.New(cleaner.c)
		if err != nil {
			log.Printf("【Memory】RCON 客户端启动失败 【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		defer c.Close()
	}
}

// clean 使用 RAMMap 清理无用内存
func (cleaner *cleaner) clean() {
	_, _ = cleaner.getMemoryInfo()
	log.Printf("【Memory】正在清理内存....\n")
	cmd := exec.Command(cleaner.rammapFile, "-Ew")
	err := cmd.Run()
	if err != nil {
		log.Printf("【Memory】运行 RAMMap 时发生错误 【%v】\n", err)
		if strings.Contains(err.Error(), "The requested operation requires elevation") {
			log.Printf("【Memory】~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("【Memory】~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("【Memory】~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
		}
	}
	log.Printf("【Memory】清理内存成功\n")
	_, _ = cleaner.getMemoryInfo()
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
		log.Printf("【Memory】获取内存信息失败【%v】\n", err)
		return
	}

	total = memStatus.ullTotalPhys / 1024 / 1024 // MB
	free = memStatus.ullAvailPhys / 1024 / 1024  // MB

	log.Printf("【Memory】空闲内存【%d MB】\n", free)

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
