package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Verthandii/palworld-go/backup"
	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/memory"
	"github.com/Verthandii/palworld-go/supervisor"
)

func main() {
	fmt.Println("★《幻兽帕鲁》启动器加载成功★")

	var (
		ctx = context.Background()
		cfg = config.Init()
		ch  = make(chan time.Duration) // 用于通知 Supervisor 过多少秒后再次检查进程是否存活
	)

	cfg.PrintLog()

	spvr, err := supervisor.New(cfg, ch)
	if err != nil {
		panic(err)
	}
	go spvr.Start(ctx)

	cleaner := memory.NewCleaner(cfg, ch)
	go cleaner.Schedule(ctx)
	defer cleaner.Stop()

	backuper := backup.New(cfg)
	go backuper.Schedule(ctx)

	signal := listenSignal()
	log.Printf("收到信号【%v】, 退出程序\n", signal)
}
