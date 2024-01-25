//go:build unix

package memory

import (
	"context"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
		c, err := rcon.New(cfg.Address, cfg.AdminPassword)
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
