package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/store"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
)

// GoodsSyncPrice 同步ERP商品到商城
type GoodsSyncPrice struct {
}

func (g GoodsSyncPrice) GetName() string {
    return "GoodsSyncPrice"
}

func (g GoodsSyncPrice) Run(t *Task) error {
    // 取出ERP全量数据
    var erpData []erp_entity.GoodsPrice

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
    newMap := safemap.New[*erp_entity.GoodsPrice]()
    for _, item := range erpData {
        newMap.Set(item.GoodsErpSpid, &item)
    }
    erpData = nil

    // 比对数据差异
    add, update, del := sync_tool.DiffMap[*erp_entity.GoodsPrice](store.GoodsPriceStore, newMap)
    newMap = nil

    // 添加
    for _, v := range add.Values() {
        addOrUpdateGoodsPrice(v)
        store.GoodsPriceStore.Set(v.GoodsErpSpid, v)
    }

    // 更新
    for _, v := range update.Values() {
        addOrUpdateGoodsPrice(v)
        store.GoodsPriceStore.Set(v.GoodsErpSpid, v)
    }

    // 删除
    for _, v := range del.Values() {
        delGoodsPrice(v)
        store.GoodsPriceStore.Delete(v.GoodsErpSpid)
    }

    // 缓存数据到文件
    err := store.SaveGoodsPrice()
    if err != nil {
        return err
    }

    return nil
}

func addOrUpdateGoodsPrice(goods *erp_entity.GoodsPrice) {
    // TODO 执行业务操作

}

func delGoodsPrice(goods *erp_entity.GoodsPrice) {
    // TODO 执行业务操作
}
