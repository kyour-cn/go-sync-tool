package initialize

import (
	"app/internal/config"
	"app/internal/global"
	"app/internal/task"
	"log/slog"
)

func InitApp() {

	// 初始化日志
	err := InitLog()
	if err != nil {
		panic(err)
	}

	slog.Info("应用启动，初始化中...")

	// 获取应用配置
	appConf, err := config.GetAppConfig()
	if err != nil {
		panic(err)
	}

	// 设置ERP编码
	global.State.ErpEncoding = appConf.ErpEncoding

	// 初始化任务进程
	task.Init()

	global.State.Status = 1

	slog.Info("应用启动，初始化完成。")
}
