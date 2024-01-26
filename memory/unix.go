//go:build linux

package memory

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

func NewCleaner(c *config.Config, ch chan<- time.Duration) Cleaner {
	return &cleaner{
		c:  c,
		ch: ch,
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
		log.Printf("【Memory】获取内存信息失败【%v】\n", err)
		return
	}

	memoryUsage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		log.Printf("【Memory】解析内存信息失败【%v】\n", err)
		return
	}

	if memoryUsage > threshold {
		log.Printf("【Memory】内存占用超过【%v】, 重新启动游戏服务器\n", threshold)
		c, err := rcon.New(cfg)
		if err != nil {
			log.Printf("【Memory】RCON 客户端启动失败【%v】\n", err)
			return
		}
		c.HandleMemoryUsage(threshold)
		c.Close()
		cleaner.ch <- 70 * time.Second
	}
}

// clean 暂无清理内存的方法
func (cleaner *cleaner) clean() {
	return
}
