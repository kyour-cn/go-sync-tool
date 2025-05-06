package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/orm/shop_model"
    "app/internal/orm/shop_query"
    "app/internal/store"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
    "gorm.io/gorm"
    "log/slog"
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

    slog.Info("商品价格同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

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

func addOrUpdateGoodsPrice(item *erp_entity.GoodsPrice) {
    shopGoods, err := shop_query.Goods.
        Where(shop_query.Goods.GoodsErpSpid.Eq(item.GoodsErpSpid)).
        Select(
            shop_query.Goods.GoodsID,
            shop_query.Goods.PriceSync,
        ).
        First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        slog.Error("goodsPriceSync Goods First err: " + err.Error())
        return
    }
    if shopGoods == nil || shopGoods.PriceSync == 0 {
        return
    }

    // 更新Goods表
    goodsData := shop_model.Goods{
        Price:       item.Price,
        CostPrice:   item.CostPrice,
        MarketPrice: item.MarketPrice,
    }
    _, e := shop_query.Goods.
        Select(
            shop_query.Goods.Price,
            shop_query.Goods.CostPrice,
            shop_query.Goods.MarketPrice,
        ).
        Where(shop_query.Goods.GoodsID.Eq(shopGoods.GoodsID)).
        Updates(goodsData)
    if e != nil {
        slog.Error("goodsPriceSync Goods update err" + e.Error())
        return
    }
    // 更新GoodsSku表
    skuData := shop_model.GoodsSku{
        Price:       item.Price,
        CostPrice:   item.CostPrice,
        MarketPrice: item.MarketPrice,
    }
    _, ers := shop_query.GoodsSku.
        Select(
            shop_query.GoodsSku.Price,
            shop_query.GoodsSku.CostPrice,
            shop_query.GoodsSku.MarketPrice,
        ).
        Where(shop_query.GoodsSku.GoodsID.Eq(shopGoods.GoodsID)).
        Updates(skuData)
    if ers != nil {
        slog.Error("goodsPriceSync GoodsSku update err: " + ers.Error())
        return
    }

}

func delGoodsPrice(goods *erp_entity.GoodsPrice) {

    // 更新价格为0
    goods.Price = 0
    goods.CostPrice = 0
    goods.MarketPrice = 0
    addOrUpdateGoodsPrice(goods)
}
