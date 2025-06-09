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

func NewOrderInvoice() *OrderInvoice {
	return &OrderInvoice{}
}

// OrderInvoice 同步ERP订单出库到商城
type OrderInvoice struct{}

func (o OrderInvoice) GetName() string {
	return "orderInvoice"
}

func (OrderInvoice) ClearCache() error {
	return store.OrderInvoiceStore.Clear()
}

func (o OrderInvoice) Run(t *Task) error {

	defer func() {
		// 缓存数据到文件
		err := store.OrderInvoiceStore.Save()
		if err != nil {
			slog.Error("SaveGoodsStock err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.OrderInvoice

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
	newMap := safemap.New[*erp_entity.OrderInvoice]()
	for _, item := range erpData {
		newMap.Set(item.InvoiceID, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.OrderInvoice](store.OrderInvoiceStore.Store, newMap)
	newMap = nil

	slog.Info("商品库存同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.OrderInvoice) error {
		err := o.add(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.OrderInvoiceStore.Store.Set(v.InvoiceID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.OrderInvoice) error {
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.OrderInvoice) error {
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (o OrderInvoice) add(v *erp_entity.OrderInvoice) error {
	// 查询商城里面是否存在该商品
	shopData, err := shop_query.ErpInvoice.
		Where(shop_query.ErpInvoice.InvoiceID.Eq(v.InvoiceID)).
		Select(shop_query.ErpInvoice.InvoiceID).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("发票数据获取 err: " + err.Error())
		return err
	}
	// 如果已经存在则跳过
	if shopData != nil {
		return nil
	}

	invoice := shop_model.ErpInvoice{
		InvoiceID:        v.InvoiceID,
		InvoiceCode:      v.InvoiceCode.String(),
		InvoiceNo:        v.InvoiceNo.String(),
		SecurityCode:     v.SecurityCode.String(),
		InvoiceImgURL:    v.InvoiceImgURL.String(),
		InvoicePdfURL:    v.InvoicePDFURL.String(),
		OrderNo:          v.OrderNo,
		OriginCreateTime: v.CreateTime,
		CreateTime:       uint(time.Now().Unix()),
	}

	if v.OrderNo != "" {
		order, _ := shop_query.Order.Where(shop_query.Order.OrderNo.Eq(v.OrderNo)).Select(shop_query.Order.MemberID).First()
		if order != nil {
			invoice.OrderNo = order.OrderNo
			invoice.MemberID = order.MemberID
		}
	}
	if v.ErpUID != "" && invoice.MemberID == 0 {
		mb, _ := shop_query.Member.Where(shop_query.Member.ErpUID.Eq(v.ErpUID)).Select(shop_query.Member.MemberID).First()
		if mb != nil {
			invoice.MemberID = mb.MemberID
		}
	}

	if err := shop_query.ErpInvoice.Create(&invoice); err != nil {
		slog.Error("发票数据 Create err:" + err.Error())
		return err
	}
	return nil
}

// ConfigLayout 任务配置UI布局
func (o OrderInvoice) ConfigLayout(_ layout.Context, _ *apptheme.Theme, _ *Task) layout.Dimensions {
	return layout.Dimensions{}
}
