package main

import (
	"app/internal/initialize"
	"context"
	"github.com/go-gourd/gourd/event"
)

func main() {

	// 初始化
	initialize.InitApp()

	event.Trigger("task.start", context.Background())

	select {}
}
