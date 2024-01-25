//go:build windows

package supervisor

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

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

func (s *supervisor) isAlive() bool {
	out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq "+s.config.ProcessName).Output()
	if err != nil {
		log.Printf("Supervisor 健康检查失败 【%v】\n", err)
		return false
	}
	return strings.Contains(string(out), s.config.ProcessName)
}

func (s *supervisor) restart() {
	cfg := s.config
	command := filepath.Join(cfg.GamePath, cfg.ProcessName+".exe")
	cmd := exec.Command(command)
	cmd.Dir = cfg.GamePath // 设置工作目录为游戏路径
	if err := cmd.Start(); err != nil {
		log.Printf("服务器重启失败【%v】\n", err)
	} else {
		log.Printf("服务器重启成功\n")
	}
}
