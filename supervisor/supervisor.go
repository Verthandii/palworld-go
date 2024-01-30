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
	ch     <-chan time.Duration
}

func New(cfg *config.Config, ch <-chan time.Duration) (Supervisor, error) {
	return &supervisor{
		config: cfg,
		ch:     ch,
	}, nil
}

func (s *supervisor) Start(ctx context.Context) {
	log.Printf("【Supervisor】启动成功，开始守护游戏进程\n")

	s.restart()

	checkDuration := time.Duration(s.config.ProcessCheckInterval) * time.Second
	ticker := time.NewTicker(checkDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("【Supervisor】成功退出\n")
			return
		case <-ticker.C:
			s.restart()
			ticker.Reset(checkDuration)
		case duration := <-s.ch:
			ticker.Reset(duration)
			log.Printf("【Supervisor】服务器即将重启, 重新设置轮询周期, 避免长时间等待服务器重启\n")
		}
	}
}

func (s *supervisor) restart() {
	if s.isAlive() {
		log.Printf("【Supervisor】ALIVE\n")
		return
	}
	log.Printf("【Supervisor】正在尝试重新启动服务器\n")

	cfg := s.config
	initCommand := filepath.Join(cfg.GamePath, cfg.ProcessName)

	cmd := exec.Command(initCommand, s.usePerfThreads()...)
	cmd.Dir = cfg.GamePath // 设置工作目录为游戏路径
	if err := cmd.Start(); err != nil {
		log.Printf("【Supervisor】服务器启动失败【%v】\n", err)
	} else {
		log.Printf("【Supervisor】服务器启动成功\n")
	}

	go func() { _ = cmd.Wait() }()
}

func (s *supervisor) usePerfThreads() []string {
	if s.config.UsePerfThreads {
		return []string{"-useperfthreads", "-NoAsyncLoadingThread", "-UseMultithreadForDS"}
	}
	return []string{}
}
