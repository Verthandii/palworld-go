package main

import (
	"context"
	"log"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/memory"
	"github.com/Verthandii/palworld-go/supervisor"
)

func main() {
	ctx := context.Background()
	cfg := config.CFG()
	spvr, err := supervisor.New()
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
