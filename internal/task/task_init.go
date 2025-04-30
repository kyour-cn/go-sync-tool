package task

import (
    "app/internal/config"
    "app/internal/global"
    "context"
    "github.com/go-gourd/gourd/event"
    "log/slog"
)

// Init 初始化任务 协程启动
func Init() {

    // 监听事件 -启动
    event.Listen("task.start", func(context.Context) {
        slog.Info("触发事件：任务启动")
        start()
    })

    // 监听事件 -停止
    event.Listen("task.stop", func(context.Context) {
        slog.Info("触发事件：任务停止")
        stop()
    })

}

func startErr(msg string) {
    slog.Error(msg)
    event.Trigger("tips.show", context.WithValue(context.Background(), "tipMsg", msg))
    global.State.Status = 1
}

func start() {

    if global.State.Status != 1 {
        slog.Warn("任务未在待启动状态，稍后再试")
        return
    }

    global.State.Status = 2

    // 获取配置
    taskConf, err := config.GetTaskConfigAll()
    if err != nil {
        startErr("获取配置失败")
        return
    }
    if taskConf == nil {
        startErr("未勾选运行的任务项")
        return
    }

    err = global.ConnDb()
    if err != nil {
        startErr("连接数据库失败")
        slog.Error("连接数据库失败", "err", err)
        return
    }

    // 遍历配置的任务
    for _, tc := range *taskConf {
        if tc.Status {
            // 匹配任务
            for _, v := range List {
                if v.Name == tc.Name {
                    v.Config = tc
                    startTask(v)
                }
            }
        }
    }

    // 运行中
    global.State.Status = 3

    // 提示
    params := context.WithValue(context.Background(), "tipMsg", "启动成功")
    event.Trigger("tips.show", params)
}

func stop() {

    if global.State.Status == 1 {
        slog.Warn("任务还未启动")
        return
    }
    if global.State.Status == 4 {
        slog.Warn("任务正在停止中")
        return
    }

    global.State.Status = 4

    // 停止操作

    err := global.CloseDb()
    if err != nil {
        startErr("关闭数据库失败")
        slog.Error("关闭数据库失败", "err", err)
        return
    }

    global.State.Status = 1

    // 提示
    params := context.WithValue(context.Background(), "tipMsg", "停止成功")
    event.Trigger("tips.show", params)
}

func startTask(task Task) {

}
