package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_model"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"app/ui/apptheme"
	"errors"
	"gioui.org/layout"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

func NewOrderOutbound() *OrderOutbound {
	return &OrderOutbound{}
}

// OrderOutbound 同步ERP订单出库到商城
type OrderOutbound struct{}

func (o OrderOutbound) GetName() string {
	return "orderOutbound"
}

func (OrderOutbound) ClearCache() error {
	return store.OrderOutboundStore.Clear()
}

func (o OrderOutbound) Run(t *Task) error {

	defer func() {
		// 缓存数据到文件
		err := store.OrderOutboundStore.Save()
		if err != nil {
			slog.Error("SaveGoodsStock err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.OrderOutBound

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
	newMap := safemap.New[*erp_entity.OrderOutBound]()
	for _, item := range erpData {
		newMap.Set(item.OutboundNo.String(), &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.OrderOutBound](store.OrderOutboundStore.Store, newMap)
	newMap = nil

	slog.Info("商品库存同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.OrderOutBound) error {
		err := o.add(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.OrderOutboundStore.Store.Set(v.OutboundNo.String(), v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.OrderOutBound) error {
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.OrderOutBound) error {
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (o OrderOutbound) add(v *erp_entity.OrderOutBound) error {
	// 查询商城里面是否存在该商品
	shopData, err := shop_query.ErpOrderOutbound.
		Where(shop_query.ErpOrderOutbound.OutboundNo.Eq(v.OutboundNo.String())).
		Select(shop_query.ErpOrderOutbound.OutboundNo).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("outBoundSync Outbound First err: " + err.Error())
		return err
	}
	if shopData != nil {
		return nil
	}

	ob := shop_model.ErpOrderOutbound{
		OutboundNo:       v.OutboundNo.String(),
		OrderNo:          v.OrderNo.String(),
		GoodsErpSpid:     v.GoodsErpSpid.String(),
		OutboundTime:     v.OutboundTime.String(),
		Validity:         v.Validity.String(),
		OutboundNum:      v.OutboundNum,
		OutboundPrice:    v.OutboundPrice,
		LogisticsCompany: v.LogisticsCompany.String(),
		LogisticsCode:    v.LogisticsCode.String(),
		SyncTime:         time.Now(),
	}
	if err := shop_query.ErpOrderOutbound.Create(&ob); err != nil {
		slog.Error("addOutBound Create err:" + err.Error())
		return err
	}
	return nil
}

// ConfigLayout 任务配置UI布局
func (o OrderOutbound) ConfigLayout(_ layout.Context, _ *apptheme.Theme) layout.Dimensions {
	return layout.Dimensions{}
}
