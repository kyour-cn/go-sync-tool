package task

import (
    "app/internal/config"
    "app/internal/global"
    "context"
    "github.com/go-gourd/gourd/event"
    "log/slog"
)

type Task struct {
    Name        string
    Label       string
    Description string
    Status      bool // 状态 0=未运行 1=运行中
    Type        int8 // 0=读取视图 1=写入中间表
    Config      config.TaskConfig
    Handle      Handle
}

type Handle interface {
    // GetName 获取任务名称
    GetName() string
    // Run 执行任务入口
    Run(*Task) error
    // Stop 停止任务
    Stop() error
}

var List = []Task{
    {
        Name:        "goods",
        Label:       "商品资料",
        Description: "需同步到电商平台的商品基础资料",
        Handle:      GoodsSync{},
    },
    {
        Name:        "goods_price",
        Label:       "商品价格",
        Description: "需同步到电商平台的商品价格",
        Handle:      GoodsSyncPrice{},
    },
    {
        Name:        "goods_stock",
        Label:       "商品库存",
        Description: "需同步到电商平台的商品库存",
        Handle:      GoodsSyncStock{},
    },
    {
        Name:        "member",
        Label:       "客户资料",
        Description: "需同步到电商平台的客户资料",
    },
    {
        Name:        "member_address",
        Label:       "客户地址",
        Description: "需同步到电商平台的客户地址",
    },
    {
        Name:        "order",
        Label:       "订单",
        Description: "需同步到电商平台的订单",
        Type:        1,
    },
}

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
            for _, item := range List {
                if item.Name == tc.Name {
                    item.Config = tc

                    slog.Info("启动任务：" + item.Name)
                    go func() {
                        err := item.Handle.Run(&item)
                        if err != nil {
                            slog.Error("任务启动失败", "name", item.Name, "err", err)
                        }
                    }()
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

    // 停止操作
    for _, tc := range *taskConf {
        if tc.Status {
            // 匹配任务
            for _, v := range List {
                if v.Name == tc.Name {
                    v.Config = tc

                    slog.Info("启动任务：" + v.Name)
                    go func() {
                        err := v.Handle.Run(&v)
                        if err != nil {
                            slog.Error("任务启动失败", "name", v.Name, "err", err)
                        }
                    }()
                }
            }
        }
    }

    err = global.CloseDb()
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
