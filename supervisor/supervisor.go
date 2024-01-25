package supervisor

import (
	"context"
	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

type Supervisor interface {
	Start(ctx context.Context)
}

type supervisor struct {
	c      *rcon.Client
	config *config.Config
}

func New() (Supervisor, error) {
	cfg := config.CFG()
	c, err := rcon.New(cfg.Address, cfg.AdminPassword)
	if err != nil {
		return nil, err
	}

	return &supervisor{
		c:      c,
		config: cfg,
	}, nil
}

func (s *supervisor) Start(ctx context.Context) {
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
		log.Printf("服务器重启失败【%v】\n", err)
	} else {
		log.Printf("服务器重启成功\n")
	}
}

func (s *supervisor) usePerfThreads() []string {
	if s.config.UsePerfThreads {
		return []string{"-useperfthreads", "-NoAsyncLoadingThread", "-UseMultithreadForDS"}
	}
	return []string{}
}
