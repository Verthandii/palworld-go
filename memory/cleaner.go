package memory

import (
	"context"
	"log"
	"time"

	"github.com/Verthandii/palworld-go/config"
)

type Cleaner interface {
	Schedule(ctx context.Context)
	Stop()
}

type cleaner struct {
	c  *config.Config
	ch chan<- time.Duration
}

func (cleaner *cleaner) Schedule(ctx context.Context) {
	log.Printf("【Memory】启动成功，定时清理服务器内存\n")

	if cleaner.c.MemoryCleanupInterval > 0 {
		cleaner.scheduleAll(ctx)
	} else {
		cleaner.scheduleRebootClean(ctx)
	}
}

func (cleaner *cleaner) scheduleAll(ctx context.Context) {
	rebootCleanDuration := 80 * time.Second
	cleanDuration := time.Duration(cleaner.c.MemoryCleanupInterval) * time.Second
	rebootCleanTicker := time.NewTicker(rebootCleanDuration)
	cleanTicker := time.NewTicker(cleanDuration)
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
	duration := 80 * time.Second
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
