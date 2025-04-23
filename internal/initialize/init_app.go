package initialize

import (
    "app/internal/task"
    "log/slog"
)

func AppInit() {

    // 初始化日志
    err := InitLog()
    if err != nil {
        panic(err)
    }

    slog.Info("应用启动，初始化中...")

    // 初始化任务进程
    task.Init()

    slog.Info("应用启动，初始化完成。")
}
