package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/memory"
	"github.com/Verthandii/palworld-go/supervisor"
)

func main() {
	fmt.Println("★《幻兽帕鲁》启动器加载成功★")

	var (
		ctx = context.Background()
		cfg = config.Init()
	)

	spvr, err := supervisor.New(cfg)
	if err != nil {
		panic(err)
	}
	go spvr.Start(ctx)

	cleaner := memory.NewCleaner(cfg)
	go cleaner.Schedule(ctx)
	defer cleaner.Stop()

	signal := listenSignal()
	log.Printf("收到信号【%v】, 退出程序\n", signal)
}
