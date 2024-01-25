package supervisor

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Verthandii/palworld-go/config"
)

type Supervisor interface {
	Start(ctx context.Context)
}

type supervisor struct {
	config *config.Config
}

func New(cfg *config.Config) (Supervisor, error) {
	return &supervisor{
		config: cfg,
	}, nil
}

func (s *supervisor) Start(ctx context.Context) {
	s.restart()

	checkDuration := time.Duration(s.config.CheckInterval) * time.Second
	ticker := time.NewTicker(checkDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("退出 Supervisor")
			return
		case <-ticker.C:
			if !s.isAlive() {
				s.restart()
			}
			ticker.Reset(checkDuration)
		}
	}
}

func (s *supervisor) restart() {
	cfg := s.config
	initCommand := filepath.Join(cfg.GamePath, cfg.ProcessName)

	cmd := exec.Command(initCommand, s.usePerfThreads()...)
	cmd.Dir = cfg.GamePath // 设置工作目录为游戏路径
	if err := cmd.Start(); err != nil {
		log.Printf("服务器启动失败【%v】\n", err)
	} else {
		log.Printf("服务器启动成功\n")
	}
}

func (s *supervisor) usePerfThreads() []string {
	if s.config.UsePerfThreads {
		return []string{"-useperfthreads", "-NoAsyncLoadingThread", "-UseMultithreadForDS"}
	}
	return []string{}
}
