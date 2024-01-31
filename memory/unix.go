//go:build !windows

package memory

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

//go:embed clean_memory.sh
var cleanMemoryFS embed.FS

var cleanMemoryFile string

func NewCleaner(c *config.Config, ch chan<- time.Duration) Cleaner {
	cleaner := &cleaner{
		c:  c,
		ch: ch,
	}

	var err error
	if c.MemoryCleanupInterval > 0 {
		cleanMemoryFile, err = extractCleanMemoryShell()
		if err != nil {
			log.Printf("无法提取 clean_memory.sh【%v】\n", err)
			os.Exit(1)
		}
	}

	return cleaner
}

func (cleaner *cleaner) Stop() {

}

// rebootClean 当内存超于阈值时，重启进程以清理内存
func (cleaner *cleaner) rebootClean() {
	cfg := cleaner.c
	threshold := cfg.MemoryUsageThreshold

	output, err := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'").Output()
	if err != nil {
		log.Printf("获取内存信息失败【%v】\n", err)
		return
	}

	memoryUsage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		log.Printf("解析内存信息失败【%v】\n", err)
		return
	}

	if memoryUsage > threshold {
		log.Printf("内存占用超过【%v】%%, 重新启动游戏服务器\n", threshold)
		c, err := rcon.New(cfg)
		if err != nil {
			log.Printf("RCON 客户端启动失败【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		c.Close()
		cleaner.ch <- 70 * time.Second
	}
}

func (cleaner *cleaner) clean() {
	free, err := cleaner.getMemoryInfo()
	if err != nil {
		log.Printf("获取内存信息失败【%v】\n", err)
		return
	}

	log.Printf("空闲内存【%d】MB, 正在清理内存....\n", free)
	cmd := exec.Command("sh", cleanMemoryFile)
	err = cmd.Run()
	if err != nil {
		log.Printf("运行 clean_memory.sh 时发生错误 【%v】\n", err)
		if strings.Contains(err.Error(), "The requested operation requires elevation") {
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
			log.Printf("~~~~~~~请以【管理员权限】打开终端~~~~~~~\n")
		}
	}
	free, err = cleaner.getMemoryInfo()
	if err != nil {
		log.Printf("获取内存信息失败【%v】\n", err)
		return
	}
	log.Printf("清理内存成功, 空闲内存【%d】MB\n", free)
}

func (cleaner *cleaner) getMemoryInfo() (uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "MemAvailable:":
			availableKB, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return 0, err
			}
			return availableKB / 1024, nil // KB -> MB
		}
	}

	if err = scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("available memory info not found")
}

func extractCleanMemoryShell() (string, error) {
	shellData, err := fs.ReadFile(cleanMemoryFS, "clean_memory.sh")
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "clean_memory-*.sh")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = tmpFile.Write(shellData); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
