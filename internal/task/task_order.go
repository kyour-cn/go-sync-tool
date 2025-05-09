package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/orm/shop_model"
    "app/internal/orm/shop_query"
    "app/internal/tools/safemap"
    "errors"
    "fmt"
    "gorm.io/gorm"
    "log/slog"
    "time"
)

func NewOrder() *Order {
    return &Order{
        orderTable:      "jxkj_sync_order",
        orderGoodsTable: "jxkj_sync_order_goods",
    }
}

// Order 商城的订单同步到ERP
type Order struct {
    orderTable      string
    orderGoodsTable string
}

func (o Order) GetName() string {
    return "Order"
}

func (o Order) Run(t *Task) error {

    // 同步新订单
    err := o.syncNewOrder(t)
    if err != nil {
        return err
    }

    // 刷新ERP订单状态
    err = o.refreshErpStatus()
    if err != nil {
        return err
    }

    return nil
}

func (o Order) syncNewOrder(t *Task) error {
    orderList, err := shop_query.Order.
        Preload(shop_query.Order.OrderGoods).
        Preload(shop_query.Order.OrderGoods.Goods).
        Preload(shop_query.Order.SettlementType).
        Preload(shop_query.Order.StaffSalesman).
        Preload(shop_query.Order.Member).
        Where(
            shop_query.Order.SyncTime.Eq(0),
            shop_query.Order.OrderStatus.In(0, 1, 3),
            shop_query.Order.PromotionType.Neq("pointexchange"), //积分订单不同步
        ).
        Find()
    if err != nil {
        slog.Error("getShopOrder orderList err: " + err.Error())
        return err
    }
    if len(orderList) == 0 {
        return nil
    }

    orderMap := safemap.New[*shop_model.Order]()
    for _, v := range orderList {
        orderMap.Set(v.OrderNo, v)
    }
    orderList = nil

    err = batchProcessor(*orderMap.GetMap(), func(item *shop_model.Order) error {
        // 添加订单到ERP
        err := o.add(item)
        if err != nil {
            // 这里忽略错误，否则将中断任务
            return nil
        }

        return nil
    }, 1, t.Ctx)
    if err != nil {
        return err
    }

    return nil
}

func (o Order) add(order *shop_model.Order) error {

    var receiptTypeMap = map[int32]string{
        1: "无",
        2: "普票",
        3: "专票",
    }

    if order.Member.MemberID == 0 {
        slog.Error("会员信息为空", "order_id", order.OrderID)
        return nil
    }

    orderStatus := order.OrderStatus
    if order.OrderStatus > 1 {
        orderStatus = 1
    }

    // 如无ID则传入名称
    if order.StaffSalesman.ErpSaleerID == "" {
        order.StaffSalesman.ErpSaleerID = order.StaffSalesman.SalesmanName
    }

    slog.Info("同步新订单", "order_id", order.OrderID, "order_no", order.OrderNo)

    type u8 erp_entity.UTF8String

    orderData := map[string]any{
        "order_id":      order.OrderID,
        "order_no":      u8(order.OrderNo),
        "erp_uid":       u8(order.Member.ErpUID),
        "pay_time":      time.Unix(int64(order.CreateTime), 0).Format("2006-01-02 15:04:05"),
        "order_money":   order.OrderMoney,
        "pay_money":     order.PayMoney,
        "pay_type":      u8(order.PayTypeName),
        "settle_type":   u8(order.SettlementType.Name),
        "name":          u8(order.Name),
        "mobile":        u8(order.Mobile),
        "address":       u8(order.Address),
        "full_address":  u8(order.FullAddress),
        "buyer_message": u8(order.BuyerMessage),
        "order_status":  orderStatus,
        "member_id":     order.MemberID,
        "coupon_money":  order.CouponMoney,
        "postage":       order.DeliveryMoney,
        "saleer_id":     u8(order.StaffSalesman.ErpSaleerID),
        "invoice_type":  receiptTypeMap[order.ReceiptType],
    }

    // 订单商品表数据
    allOrderGoodsData := make([]map[string]any, 0, len(order.OrderGoods))

    for _, vs := range order.OrderGoods {
        goods, er := shop_query.Goods.
            Where(shop_query.Goods.GoodsID.Eq(vs.GoodsID)).
            First()
        if er != nil {
            slog.Error(fmt.Sprintf("erpSyncOrder goods first err:%s,goods:%+v", er, vs))
            continue
        }
        if goods == nil || vs.Num <= 0 || goods.GoodsErpSpid == "" {
            slog.Debug(fmt.Sprintf("订单商品数据不正确:%+v", goods))
            continue
        }
        price := vs.RealGoodsMoney / float64(vs.Num)

        allOrderGoodsData = append(allOrderGoodsData, map[string]any{
            "order_goods_id": vs.OrderGoodsID,
            "order_no":       erp_entity.UTF8String(vs.OrderNo),
            "goods_erp_spid": erp_entity.UTF8String(goods.GoodsErpSpid),
            "price":          price,
            "original_price": vs.OriginalPrice,
            "num":            vs.Num,
            "goods_money":    vs.RealGoodsMoney,
            "status":         orderStatus,
        })
    }

    // 获取ERP数据库连接
    erpDb, ok := global.DbPool.Get("erp")
    if !ok {
        return errors.New("获取ERP数据库连接失败")
    }

    // 创建订单
    result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
        return tx.Table(o.orderTable).Create(orderData)
    }))
    if result.Error != nil {
        slog.Error(fmt.Sprintf("ERP CreateOrder err:%s,args:%+v", result.Error, orderData))
        return nil
    }

    //创建订单商品
    if len(allOrderGoodsData) > 0 {
        for _, og := range allOrderGoodsData {
            result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
                return tx.Table(o.orderGoodsTable).Create(og)
            }))
            if result.Error != nil {
                slog.Error(fmt.Sprintf("ERP CreateOrderGoods err:%s,args:%+v", result.Error, og))
                return nil
            }
        }
    }

    // 修改商城数据为已同步
    if _, er := shop_query.Order.
        Where(shop_query.Order.OrderID.Eq(order.OrderID)).
        Update(shop_query.Order.SyncTime, time.Now().Unix()); er != nil {
        slog.Error("更新订单同步状态失败: " + er.Error())
        return nil
    }

    return nil
}

// 返回该分组的订单总金额
func (o Order) getOrderPrice(orderGoods map[string][]shop_model.OrderGoods, area string) float64 {
    orderPrice := 0.00
    for k, goods := range orderGoods {
        for _, good := range goods {
            if k == area {
                orderPrice += good.RealGoodsMoney
            }
        }

    }
    return orderPrice
}

// 刷新已同步订单的状态
func (o Order) refreshErpStatus() error {

    erpDb, ok := global.DbPool.Get("erp")
    if !ok {
        return errors.New("获取ERP数据库连接失败")
    }

    // 将3天前的待提取[1]订单修改为异常超时[3]（兼容处理，部分erp方未更新状态）
    lastThreeDays := time.Now().AddDate(0, 0, -3).Format("2006-01-02 15:04:05")
    result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
        return tx.Table(o.orderTable).
            Where("order_status = ? AND pay_time < ?", 1, lastThreeDays).
            Update("order_status", 3)
    }))
    if result.Error != nil {
        slog.Error("更新ERP订单提取超时状态异常: " + result.Error.Error())
    }

    sqO := shop_query.Order

    // 取出ERP未支付数据
    var erpData []erp_entity.SyncOrder
    result = erpDb.Table(o.orderTable).
        Where("order_status = 0").
        Select("order_no", "order_status").
        Find(&erpData)
    if result.Error != nil {
        return result.Error
    }
    // 同步erp未支付订单
    if len(erpData) > 0 {
        orderNoList := make([]string, 0, len(erpData))
        for _, v := range erpData {
            orderNoList = append(orderNoList, v.OrderNo)
        }
        //查询商城里面的订单
        shopOrders, err := sqO.
            Or(
                sqO.Where(sqO.OrderNo.In(orderNoList...)),
                sqO.Where(sqO.OrderStatus.Neq(0)),
            ).
            Select(sqO.OrderNo, sqO.OrderStatus).
            Find()
        if err != nil {
            slog.Error("查询商城订单信息异常: " + err.Error())
            return err
        }
        for _, item := range shopOrders {
            // ERP只接受 0,1,-1
            if item.OrderStatus > 1 {
                item.OrderStatus = 1
            }
            // 更新ERP状态
            slog.Debug("更新ERP订单状态: " + item.OrderNo)
            _ = o.updateStatus(item, item.OrderStatus)
        }
    }

    // 同步ERP未退款订单 只取3天内的待提取订单
    erpData = []erp_entity.SyncOrder{}
    result = erpDb.Table(o.orderTable).
        Where("order_status = 1").
        Select("order_no", "order_status").
        Find(&erpData)
    if result.Error != nil {
        return result.Error
    }
    if len(erpData) > 0 {
        orderNoList := make([]string, 0, len(erpData))
        for _, v := range erpData {
            orderNoList = append(orderNoList, v.OrderNo)
        }
        //查询商城里面的订单
        shopOrders, err := sqO.
            Or(
                sqO.Where(sqO.OrderNo.In(orderNoList...)),
            ).
            Select(sqO.OrderNo, sqO.OrderStatus).
            Find()
        if err != nil {
            slog.Error("查询商城订单信息异常: " + err.Error())
            return err
        }
        for _, item := range shopOrders {

            // 已全部退款 - 关闭
            if item.OrderStatus == -1 {
                _ = o.updateStatus(item, item.OrderStatus)
                continue
            }

            // 查询部分退款
            refundedGoods, err := shop_query.OrderGoods.
                Preload(shop_query.OrderGoods.Goods).
                Where(
                    shop_query.OrderGoods.OrderNo.Eq(item.OrderNo),
                    shop_query.OrderGoods.RefundStatus.Eq(3), //退款完成
                ).Find()
            if err != nil {
                return err
            }
            if len(refundedGoods) == 0 {
                continue
            }
            // 查询ERP未退款的商品
            var unRefundedGoods []erp_entity.SyncOrderGoods
            result = erpDb.Table(o.orderGoodsTable).
                Where(fmt.Sprintf("status != -1 AND AND order_no = '%s'", item.OrderNo)).
                Select("order_no", "status").
                Find(&erpData)
            if result.Error != nil {
                return result.Error
            }
            if len(unRefundedGoods) == 0 {
                continue
            }

            for _, rog := range refundedGoods {
                if rog.Goods.GoodsErpSpid == "" {
                    //订单商品未关联到商品信息
                    continue
                }
                //更新商城已退款erp未退款的订单单个商品
                for _, urGoods := range unRefundedGoods {
                    if rog.Goods.GoodsErpSpid != urGoods.GoodsErpSpid.String() {
                        continue
                    }
                    _ = o.updateOGStatus(urGoods.OrderNo.String(), -1)
                }
            }
        }
    }

    return nil
}

func (o Order) updateStatus(order *shop_model.Order, status int32) error {
    erpDb, ok := global.DbPool.Get("erp")
    if !ok {
        return errors.New("获取ERP数据库连接失败")
    }
    // 更新ERP状态
    result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
        return tx.Table(o.orderTable).
            Where("order_no = ? AND order_status <> ?", order.OrderNo, 2).
            Update("order_status", status)
    }))
    if result.Error != nil {
        slog.Error("更新ERP订单状态异常: " + result.Error.Error())
        return result.Error
    }

    // 更新明细状态
    _ = o.updateOGStatus(order.OrderNo, status)

    return nil
}

func (o Order) updateOGStatus(orderNo string, status int32) error {
    erpDb, ok := global.DbPool.Get("erp")
    if !ok {
        return errors.New("获取ERP数据库连接失败")
    }
    // 更新ERP状态
    result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
        return tx.Table(o.orderGoodsTable).
            Where("order_no = ?", orderNo).
            Update("order_status", status)
    }))
    if result.Error != nil {
        slog.Error("更新ERP订单明细状态异常: " + result.Error.Error())
        return result.Error
    }

    return nil
}
