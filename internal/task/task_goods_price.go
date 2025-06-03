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
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"log/slog"
)

func NewGoodsPrice() *GoodsPrice {
	return &GoodsPrice{}
}

// GoodsPrice 同步ERP商品到商城
type GoodsPrice struct{}

func (g GoodsPrice) GetName() string {
	return "GoodsPrice"
}

func (GoodsPrice) ClearCache() error {
	return store.GoodsPriceStore.Clear()
}

func (g GoodsPrice) Run(t *Task) error {
	defer func() {
		// 缓存数据到文件
		err := store.GoodsPriceStore.Save()
		if err != nil {
			slog.Error("SaveGoodsPrice err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.GoodsPrice

	erpDb, ok := global.DbPool.Get("erp")
	if !ok {
		return errors.New("获取ERP数据库连接失败")
	}

	// 执行SQL查询
	r := erpDb.Raw(t.Config.Sql).Scan(&erpData)
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
	add, update, del := sync_tool.DiffMap[*erp_entity.GoodsPrice](store.GoodsPriceStore.Store, newMap)
	newMap = nil

	slog.Info("商品价格同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.GoodsPrice) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsPriceStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.GoodsPrice) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsPriceStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.GoodsPrice) error {
		err := g.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsPriceStore.Store.Delete(v.GoodsErpSpid)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (g GoodsPrice) addOrUpdate(item *erp_entity.GoodsPrice) error {
	shopGoods, err := shop_query.Goods.
		Where(shop_query.Goods.GoodsErpSpid.Eq(item.GoodsErpSpid)).
		Select(
			shop_query.Goods.GoodsID,
			shop_query.Goods.PriceSync,
		).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("goodsPriceSync Goods First err: " + err.Error())
		return err
	}
	if shopGoods == nil || shopGoods.PriceSync == 0 {
		return nil
	}

	// 要更新的字段
	var updateColumns []field.Expr
	var updateColumns2 []field.Expr

	if item.Price > -1 && shopGoods.Price != item.Price {
		updateColumns = append(updateColumns, shop_query.Goods.Price)
		updateColumns2 = append(updateColumns2, shop_query.GoodsSku.Price)
	}
	if item.CostPrice > -1 && shopGoods.CostPrice != item.CostPrice {
		updateColumns = append(updateColumns, shop_query.Goods.CostPrice)
		updateColumns2 = append(updateColumns2, shop_query.GoodsSku.CostPrice)
	}
	if item.MarketPrice > -1 && shopGoods.MarketPrice != item.MarketPrice {
		updateColumns = append(updateColumns, shop_query.Goods.MarketPrice)
		updateColumns2 = append(updateColumns2, shop_query.GoodsSku.MarketPrice)
	}
	if len(updateColumns) == 0 {
		return nil
	}

	slog.Debug("价格更新", "spid", item.GoodsErpSpid, "old", shopGoods.Price, "new", item.Price)

	// 更新Goods表
	goodsData := shop_model.Goods{
		Price:       item.Price,
		CostPrice:   item.CostPrice,
		MarketPrice: item.MarketPrice,
	}
	_, e := shop_query.Goods.
		Select(updateColumns...).
		Where(shop_query.Goods.GoodsID.Eq(shopGoods.GoodsID)).
		Updates(goodsData)
	if e != nil {
		slog.Error("goodsPriceSync Goods update err" + e.Error())
		return e
	}
	// 更新GoodsSku表
	skuData := shop_model.GoodsSku{
		Price:       item.Price,
		CostPrice:   item.CostPrice,
		MarketPrice: item.MarketPrice,
	}

	_, ers := shop_query.GoodsSku.
		Select(updateColumns2...).
		Where(shop_query.GoodsSku.GoodsID.Eq(shopGoods.GoodsID)).
		Updates(skuData)
	if ers != nil {
		slog.Error("goodsPriceSync GoodsSku update err: " + ers.Error())
		return ers
	}

	return nil
}

func (g GoodsPrice) delete(goods *erp_entity.GoodsPrice) error {
	// 更新价格为0
	goods.Price = 0
	goods.CostPrice = 0
	goods.MarketPrice = 0
	return g.addOrUpdate(goods)
}
