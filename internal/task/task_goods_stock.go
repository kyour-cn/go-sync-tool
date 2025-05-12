package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

func NewGoodsStock() *GoodsStock {
	return &GoodsStock{}
}

// GoodsStock 同步ERP商品到商城
type GoodsStock struct {
	IsRunning bool
}

func (g GoodsStock) GetName() string {
	return "GoodsStock"
}

func (g GoodsStock) Run(t *Task) error {

	defer func() {
		// 缓存数据到文件
		err := store.GoodsStockStore.Save()
		if err != nil {
			slog.Error("SaveGoodsStock err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.GoodsStock

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
	newMap := safemap.New[*erp_entity.GoodsStock]()
	for _, item := range erpData {
		newMap.Set(item.GoodsErpSpid, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.GoodsStock](store.GoodsStockStore.Store, newMap)
	newMap = nil

	slog.Info("商品库存同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.GoodsStock) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStockStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.GoodsStock) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStockStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.GoodsStock) error {
		err := g.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStockStore.Store.Delete(v.GoodsErpSpid)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (g GoodsStock) addOrUpdate(item *erp_entity.GoodsStock) error {

	shopGoods, err := shop_query.Goods.
		Where(shop_query.Goods.GoodsErpSpid.Eq(item.GoodsErpSpid)).
		Select(
			shop_query.Goods.GoodsID,
			shop_query.Goods.StockSync,
			shop_query.Goods.GoodsStock,
		).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("goodsPriceSync Goods First err: " + err.Error())
		return err
	}
	if shopGoods == nil {
		return nil
	}

	newStock := int32(item.GoodsStock)
	if shopGoods.StockSync == 0 || shopGoods.GoodsStock == newStock {
		return nil
	}

	slog.Debug("库存更新", "spid", item.GoodsErpSpid, "old", shopGoods.GoodsStock, "new", newStock)

	// 更新Goods表
	_, e := shop_query.Goods.
		Where(shop_query.Goods.GoodsID.Eq(shopGoods.GoodsID)).
		Update(shop_query.Goods.GoodsStock, newStock)
	if e != nil {
		slog.Error("goodsPriceSync Goods update err" + e.Error())
		return e
	}
	// 更新GoodsSku表
	_, ers := shop_query.GoodsSku.
		Where(shop_query.GoodsSku.GoodsID.Eq(shopGoods.GoodsID)).
		Update(shop_query.GoodsSku.Stock, newStock)
	if ers != nil {
		slog.Error("goodsPriceSync GoodsSku update err: " + ers.Error())
		return ers
	}

	return nil
}

func (g GoodsStock) delete(goods *erp_entity.GoodsStock) error {
	// 更新价格为0
	goods.GoodsStock = 0
	return g.addOrUpdate(goods)
}
