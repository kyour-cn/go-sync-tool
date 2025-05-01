package task

// GoodsSync 同步ERP商品到商城
type GoodsSync struct {
}

func (g GoodsSync) GetName() string {
    return "GoodsSync"
}

func (g GoodsSync) Run(t Task) error {
    // 取出ERP数据

    // TODO 取出缓存数据

    // TODO 比对数据差异

    // TODO 执行业务操作

    return nil
}
