//go:build unix

package memory

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

type cleaner struct {
	c *config.Config
}

func NewCleaner(c *config.Config) Cleaner {
	return &cleaner{
		c: c,
	}
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
		log.Printf("内存占用超过 %v, 开始清理内存...\n", threshold)
		c, err := rcon.New(cfg)
		if err != nil {
			log.Printf("rcon 客户端启动失败 【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		defer c.Close()
	}
}

// clean 暂无清理内存的方法
func (cleaner *cleaner) clean() {
	return
}
