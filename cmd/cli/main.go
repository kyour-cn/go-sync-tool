package main

import (
	"app/internal/initialize"
	"context"
	"github.com/go-gourd/gourd/event"
)

// 命令行无UI版本，可编译至linux系统
func main() {

	// 初始化
	initialize.InitApp()

	// 触发任务开始
	event.Trigger("task.start", context.Background())

	select {}
}
