package task

import (
    "app/internal/config"
    "app/internal/global"
    "app/internal/store"
    "context"
    "errors"
    "github.com/go-gourd/gourd/event"
    "log/slog"
    "time"
)

type Task struct {
    Name        string
    Label       string
    Description string
    Type        int8 // 0=读取视图 1=写入中间表
    Config      config.TaskConfig
    Handle      Handle
    Parent      string // 父级任务，会等待父级任务完成一轮才会触发
    Ctx         context.Context

    Status      bool      // 运行状态 是否运行中
    LastRunTime time.Time // 上次运行时间
    DataCount   int       // 数据总数
    DoneCount   int       // 已完成数量
}

type Handle interface {
    // GetName 获取任务名称
    GetName() string
    // Run 执行任务入口
    Run(*Task) error
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
        Parent:      "goods",
    },
    {
        Name:        "goods_stock",
        Label:       "商品库存",
        Description: "需同步到电商平台的商品库存",
        Handle:      GoodsSyncStock{},
        Parent:      "goods",
    },
    {
        Name:        "member",
        Label:       "客户资料",
        Description: "需同步到电商平台的客户资料",
        Handle:      MemberSync{},
    },
    {
        Name:        "member_address",
        Label:       "客户地址",
        Description: "需同步到电商平台的客户地址",
        Parent:      "member",
        Handle:      MemberAddress{},
    },
    {
        Name:        "member_business_scope",
        Label:       "经营范围",
        Description: "需同步到电商平台的客户经营范围",
        Parent:      "member",
        Handle:      MemberAddress{},
    },
    {
        Name:        "salesman",
        Label:       "业务员",
        Description: "需同步到电商平台的客户地址",
    },
    {
        Name:        "order",
        Label:       "订单",
        Description: "需同步到电商平台的订单",
        Type:        1,
    },
}

// Init 初始化任务
func Init() {

    var cancelCtx context.Context
    var cancelFunc context.CancelFunc

    // 初始化存储库
    store.Init()

    err := initConfig(cancelCtx)
    if err != nil {
        slog.Error("初始化配置失败", "err", err)
        return
    }

    // 开启一个后台任务用于刷新UI
    go func() {
        for {
            if global.State.Status == 3 {
                event.Trigger("window.invalidate", cancelCtx)
            }
            time.Sleep(time.Second * 1)
        }
    }()

    // 监听事件 -启动
    event.Listen("task.start", func(ctx context.Context) {
        slog.Info("触发事件：任务启动")

        cancelCtx, cancelFunc = context.WithCancel(context.Background())
        go start(cancelCtx)
    })

    // 监听事件 -停止
    event.Listen("task.stop", func(context.Context) {
        slog.Info("触发事件：任务停止")
        if global.State.Status == 1 {
            slog.Warn("任务还未启动")
            return
        }
        if global.State.Status == 4 {
            slog.Warn("任务正在停止中")
            return
        }
        // 状态改为停止中
        global.State.Status = 4

        // 通知协程停止任务
        cancelFunc()
    })

    // 监听事件 -配置更改
    event.Listen("task.config", func(ctx context.Context) {
        _ = initConfig(ctx)
    })

}

func initConfig(ctx context.Context) error {
    // 获取配置
    taskConf, err := config.GetTaskConfigAll()
    if err != nil {
        return errors.New("获取配置失败")
    }
    if taskConf == nil {
        return errors.New("未勾选运行的任务项")
    }
    // 遍历配置的任务进行初始化
    for _, tc := range *taskConf {
        // 匹配任务名
        for i := range List {
            if List[i].Name == tc.Name {
                List[i].Config = tc
                List[i].Ctx = ctx
            }
        }
    }
    return nil
}

func startErr(msg string) {
    slog.Error(msg)
    event.Trigger("tips.show", context.WithValue(context.Background(), "tipMsg", msg))
    global.State.Status = 1
}

func start(ctx context.Context) {

    if global.State.Status != 1 {
        slog.Warn("任务未在待启动状态，稍后再试")
        return
    }

    global.State.Status = 2

    // 初始化任务配置
    err := initConfig(ctx)
    if err != nil {
        startErr(err.Error())
        return
    }

    err = global.ConnDb()
    if err != nil {
        startErr("连接数据库失败")
        slog.Error("连接数据库失败", "err", err)
        return
    }

    // 运行中
    global.State.Status = 3

    // 提示
    params := context.WithValue(context.Background(), "tipMsg", "启动成功")
    event.Trigger("tips.show", params)

    for {
        select {
        // 监测停止
        case <-ctx.Done():
            // 读取运行中的任务数
            running := 0
            for _, v := range List {
                if v.Status {
                    running++
                }
            }
            // 运行中任务为0时停止
            if running == 0 {
                stoped()
                return
            }
        default:
            // 执行任务
            for i, item := range List {
                // 运行启用的一级任务
                if !item.Status && item.Config.Status && item.Parent == "" {
                    go startOne(&List[i])
                }
            }
            // 延迟一秒
            time.Sleep(time.Second)
        }
    }
}

// 运行单个任务
func startOne(item *Task) {

    // 判断运行状态和时间差
    if item.Status || time.Since(item.LastRunTime) < time.Second*time.Duration(item.Config.IntervalTime) {
        return
    }

    // 修改状态和同步时间
    item.Status = true
    item.LastRunTime = time.Now()

    slog.Info("开始运行任务："+item.Label, "name", item.Name)

    // 运行业务代码
    err := item.Handle.Run(item)
    if err != nil {
        slog.Error("任务运行失败："+item.Label, "name", item.Name, "err", err)
    }

    slog.Info("任务运行完成："+item.Label, "name", item.Name, "耗时", time.Since(item.LastRunTime).String())

    item.Status = false
    item.DataCount = 0
    item.DoneCount = 0

    // 遍历运行子任务 -可实现递归
    for i, v := range List {
        if v.Parent == item.Name && v.Config.Status {
            go startOne(&List[i])
        }
    }
}

func stoped() {

    global.State.Status = 1

    err := global.CloseDb()
    if err != nil {
        startErr("关闭数据库失败")
        slog.Error("关闭数据库失败", "err", err)
        return
    }

    // 提示
    params := context.WithValue(context.Background(), "tipMsg", "停止成功")
    event.Trigger("tips.show", params)
}
