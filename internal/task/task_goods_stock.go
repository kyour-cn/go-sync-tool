package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/store"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
    "golang.org/x/exp/slog"
    "time"
)

// GoodsSyncStock 同步ERP商品到商城
type GoodsSyncStock struct {
    IsRunning bool
}

func (g GoodsSyncStock) GetName() string {
    return "GoodsSyncStock"
}

func (g GoodsSyncStock) Run(t *Task) error {
    g.IsRunning = true
    for {
        slog.Info("开始同步商品库存")
        st := time.Now()
        err := g.runLoop(t)
        if err != nil {
            slog.Error("同步商品库存失败", "err", err)
        }

        // 计算耗时
        slog.Info("同步商品库存完成，耗时：" + time.Since(st).String())

        // 间隔1分钟
        time.Sleep(time.Duration(10) * time.Second)
    }
}

func (g GoodsSyncStock) Stop() error {
    g.IsRunning = false

    return nil
}

func (g GoodsSyncStock) runLoop(t *Task) error {
    // 取出ERP全量数据
    var erpData []erp_entity.GoodsStock

    erpDb, ok := global.DbPool.Get("erp")
    if !ok {
        return errors.New("获取ERP数据库连接失败")
    }

    // 执行SQL查询
    r := erpDb.Db.Raw(t.Config.Sql).Scan(&erpData)
    if r.Error != nil {
        return r.Error
    }

    // 创建新的Map
    newMap := safemap.New[*erp_entity.GoodsStock]()
    for _, item := range erpData {
        newMap.Set(item.GoodsErpSpid, &item)
    }
    erpData = nil

    // 比对数据差异
    add, update, del := sync_tool.DiffMap[*erp_entity.GoodsStock](store.GoodsStockStore, newMap)
    newMap = nil

    // 添加
    for _, v := range add.Values() {
        addOrUpdateGoodsStock(v)
        store.GoodsStockStore.Set(v.GoodsErpSpid, v)
    }

    // 更新
    for _, v := range update.Values() {
        addOrUpdateGoodsStock(v)
        store.GoodsStockStore.Set(v.GoodsErpSpid, v)
    }

    // 删除
    for _, v := range del.Values() {
        delGoodsStock(v)
        store.GoodsStockStore.Delete(v.GoodsErpSpid)
    }

    // 缓存数据到文件
    err := store.SaveGoodsStock()
    if err != nil {
        return err
    }

    return nil
}

func addOrUpdateGoodsStock(goods *erp_entity.GoodsStock) {
    // TODO 执行业务操作

}

func delGoodsStock(goods *erp_entity.GoodsStock) {
    // TODO 执行业务操作
}
