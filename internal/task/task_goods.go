package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/store"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
)

// GoodsSync 同步ERP商品到商城
type GoodsSync struct {
}

func (g GoodsSync) GetName() string {
    return "GoodsSync"
}

func (g GoodsSync) Run(t Task) error {
    // 取出ERP全量数据
    var erpData []erp_entity.Goods

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
    newMap := safemap.New[*erp_entity.Goods]()
    for _, goods := range erpData {
        newMap.Set(goods.GoodsErpSpid, &goods)
    }
    erpData = nil

    // 比对数据差异
    add, update, del := sync_tool.DiffMap[*erp_entity.Goods](store.GoodsStore, newMap)
    newMap = nil

    // 添加
    for _, v := range add.Values() {
        addOrUpdateGoods(v)
        store.GoodsStore.Set(v.GoodsErpSpid, v)
    }

    // 更新
    for _, v := range update.Values() {
        addOrUpdateGoods(v)
        store.GoodsStore.Set(v.GoodsErpSpid, v)
    }

    // 删除
    for _, v := range del.Values() {
        delGoods(v)
        store.GoodsStore.Delete(v.GoodsErpSpid)
    }

    // 缓存数据到文件
    err := store.SaveGoods()
    if err != nil {
        return err
    }

    return nil
}

func addOrUpdateGoods(goods *erp_entity.Goods) {
    // TODO 执行业务操作
}

func delGoods(goods *erp_entity.Goods) {
    // TODO 执行业务操作
}
