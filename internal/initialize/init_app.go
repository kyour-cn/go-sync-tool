package initialize

import (
    "app/internal/global"
    "app/internal/store"
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

    // 初始化存储库
    store.Init()

    // 初始化任务进程
    task.Init()

    global.State.Status = 1

    slog.Info("应用启动，初始化完成。")
}
