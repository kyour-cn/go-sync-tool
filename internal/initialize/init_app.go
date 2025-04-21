package initialize

import "log/slog"

func AppInit() {

    // 初始化日志
    err := InitLog()
    if err != nil {
        panic(err)
    }

    slog.Info("初始化完成")
}
