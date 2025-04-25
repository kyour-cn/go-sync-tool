package task

import (
    "app/internal/config"
    "context"
    "github.com/go-gourd/gourd/event"
    "log/slog"
)

// Init 初始化任务 协程启动
func Init() {
    // 捕获异常
    defer func() {
        if err := recover(); err != nil {
            slog.Error("任务启动异常", "err", err)
        }
    }()

    // 监听事件
    event.Listen("task.start", func(context.Context) {
        slog.Info("触发事件：任务启动")
        start()

        params := context.WithValue(context.Background(), "tipMsg", "启动成功")
        event.Trigger("tips.show", params)
    })

    event.Listen("task.stop", func(context.Context) {
        slog.Info("触发事件：任务停止")
        stop()
    })

}

func start() {

    // 获取配置
    taskConf, err := config.GetSqlConfigAll()
    if err != nil {
        slog.Error("获取配置失败", "err", err)
        return
    }
    if taskConf == nil {
        slog.Warn("暂无", "err", err)
        return
    }

    for _, tc := range *taskConf {
        if tc.Status == true {

            // 匹配任务
            for _, v := range List {
                if v.Name == tc.Name {
                    startTask(tc, v)
                }
            }

        }
    }

}

func stop() {

}

func startTask(conf config.TaskConfig, task Task) {

}
