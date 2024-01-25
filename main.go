package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/memory"
	"github.com/Verthandii/palworld-go/supervisor"
)

//go:embed RAMMap64.exe
var rammapFS embed.FS

func main() {
	ctx := context.Background()
	cfg := config.CFG()
	spvr, err := supervisor.New()
	if err != nil {
		panic(err)
	}
	go spvr.Start(ctx)

	// 设置内存检查任务
	memoryCheckTask := memory.NewMemoryCheckTask()
	go memoryCheckTask.Schedule(ctx)

	if runtime.GOOS == "windows" {
		if cfg.MemoryCleanupInterval != 0 {
			log.Printf("启用 rammap 清理内存【清理过程中会导致游戏卡顿】\n")

			rammapExecutable, err := extractRAMMapExecutable()
			if err != nil {
				log.Fatalf("无法提取RAMMap可执行文件: %v", err)
			}
			defer os.Remove(rammapExecutable) // 确保程序结束时删除文件

			// 创建定时器，根据配置间隔定期运行 RAMMap
			go func(ctx context.Context) {
				duration := time.Duration(cfg.MemoryCleanupInterval) * time.Second
				ticker := time.NewTicker(duration)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						runRAMMap(rammapExecutable)
						ticker.Reset(duration)
					}
				}
			}(ctx)
		}
	}

	signal := listenSignal()
	log.Printf("收到信号【%v】, 退出程序\n", signal)
}

// extractRAMMapExecutable 从嵌入的文件系统中提取 RAMMap 并写入临时文件
func extractRAMMapExecutable() (string, error) {
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

func runRAMMap(rammapExecutable string) {
	log.Printf("正在使用 rammap 清理内存....\n")
	cmd := exec.Command(rammapExecutable, "-Ew")
	err := cmd.Run()
	if err != nil {
		log.Printf("运行 RAMMap 时发生错误 【%v】\n", err)
	}
}
