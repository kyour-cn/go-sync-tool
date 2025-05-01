package task

import (
    "app/internal/config"
    "app/internal/task/sync"
)

type Task struct {
    Name        string
    Label       string
    Description string
    Status      bool
    Type        int8 // 0=读取 1=写入
    Config      config.TaskConfig
    Handle      Handle
}

type Handle interface {
    // GetName 获取任务名称
    GetName() string
    // Run 执行任务入口
    Run() error
}

var List = []Task{
    {
        Name:        "goods",
        Label:       "商品资料",
        Description: "需同步到电商平台的商品基础资料",
        Handle:      sync.GoodsSync{},
    },
    {
        Name:        "goods_price",
        Label:       "商品价格",
        Description: "需同步到电商平台的商品价格",
    },
    {
        Name:        "goods_stock",
        Label:       "商品库存",
        Description: "需同步到电商平台的商品库存",
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
