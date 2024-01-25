package memory

import (
	"context"
	"time"
)

type Cleaner interface {
	Schedule(ctx context.Context)
	Stop()
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

func (cleaner *cleaner) scheduleAll(ctx context.Context) {
	rebootCleanDuration := time.Duration(cleaner.c.MemoryCheckInterval) * time.Second
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
