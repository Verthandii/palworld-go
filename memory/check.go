package memory

import (
	"context"
	"time"

	"github.com/Verthandii/palworld-go/config"
)

type MemoryCheckTask struct{}

func NewMemoryCheckTask() *MemoryCheckTask {
	return &MemoryCheckTask{}
}

func (task *MemoryCheckTask) Schedule(ctx context.Context) {
	cfg := config.CFG()
	duration := time.Duration(cfg.MemoryCheckInterval) * time.Second
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkMemory()
			ticker.Reset(duration)
		}
	}
}
