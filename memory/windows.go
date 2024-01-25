package memory

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

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

func checkMemory() {
	cfg := config.CFG()
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

	if 100.0*(1-float64(free)/float64(total)) > threshold {
		log.Printf("内存占用超过 %v, 开始清理内存...\n", threshold)
		// 初始化RCON客户端
		c, err := rcon.New(cfg.Address, cfg.AdminPassword)
		if err != nil {
			log.Printf("rcon 客户端启动失败 【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		defer c.Close()
	}
}
