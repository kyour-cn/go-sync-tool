package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_query"
	"errors"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

func NewOrder() *Order {
	return &Order{
		orderTableName: "sync_order",
	}
}

// Order 商城的订单同步到ERP
type Order struct {
	orderTableName string
}

func (g Order) GetName() string {
	return "Order"
}

func (g Order) Run(t *Task) error {
	// 同步订单状态
	err := g.updateStatus()
	if err != nil {
		return err
	}

	return nil
}

// 更新已同步订单的状态
func (g Order) updateStatus() error {

	// 取出ERP数据
	var erpData []erp_entity.SyncOrder

	erpDb, ok := global.DbPool.Get("erp")
	if !ok {
		return errors.New("获取ERP数据库连接失败")
	}

	result := erpDb.Table(g.orderTableName).
		Where("order_status = ?", 0).
		Select("name", "age").
		Find(erpData)
	if result.Error != nil {
		return result.Error
	}

	// 1.同步erp未支付订单
	if len(erpData) > 0 {
		unpaidNos := make([]string, 0, len(erpData))
		for _, v := range erpData {
			unpaidNos = append(unpaidNos, v.OrderNo)
		}
		//查询商城里面的订单
		timeStamp := uint(time.Now().AddDate(0, 0, -1).Unix()) //1天以内
		shopOrders, dbErr := shop_query.Order.
			Where(
				shop_query.Order.OrderStatus.In(0, -1, 1),
				shop_query.Order.CreateTime.Gt(timeStamp),
			).
			Or(
				shop_query.Order.Where(shop_query.Order.OrderNo.In(unpaidNos...)),
			).
			Find()
		if dbErr != nil {
			slog.Error("OrderSync Query Shop err: " + dbErr.Error())
			return dbErr
		}
		for _, item := range shopOrders {
			//t := time.Now().Format("2006-01-02 15:04:05.999")
			// ERP只接受 0,1,-1
			if item.OrderStatus > 1 {
				item.OrderStatus = 1
			}
			//if _, err = erp_query.ErpSyncOrder.UpdateStatusAndPayTimeByOrderNo(item.OrderNo, item.OrderStatus, t); err != nil {
			//	log.Errorf("UpdateStatusAndPayTimeByOrderNo err:%s", err)
			//	return err
			//}
			//if _, err = erp_query.ErpSyncOrderGoods.UpdateStatusByOrderNo(item.OrderNo, item.OrderStatus); err != nil {
			//	log.Errorf("UpdateStatusByOrderNo err:%s", err)
			//	return err
			//}
		}
	}
	// 2.同步ERP未退款订单
	// 查询近三天erp未退款订单
	urOrders, err := erp_query.ErpSyncOrder.GetUnRefundListByPayTime(time.Now().AddDate(0, 0, -3))
	if err != nil {
		return err
	}
	for _, uro := range urOrders {

		//查询是否有退款的商品
		shopOrder, err := shop_query.Order.
			Where(shop_query.Order.OrderNo.Eq(genShopOrderNo(uro.OrderNo.String()))).
			Select(shop_query.Order.OrderStatus).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("OrderSync Query Refund First err:%s", err)
			return err
		}
		if shopOrder == nil {
			continue
		}

		// 订单全部商品退款完成,同步状态
		if shopOrder.OrderStatus == -1 {
			if _, er := erp_query.ErpSyncOrder.UpdateStatusByOrderNo(shopOrder.OrderNo, shopOrder.OrderStatus); er != nil {
				log.Errorf("UpdateStatusByOrderNo err:%s", er)
				return er
			}
			if _, er := erp_query.ErpSyncOrderGoods.UpdateStatusByOrderNo(shopOrder.OrderNo, shopOrder.OrderStatus); er != nil {
				log.Errorf("UpdateStatusByOrderNo err:%s", er)
				return er
			}
			continue
		}
		//查询商城退款的商品
		refundedGoods, shopErr := shop_query.OrderGoods.
			Preload(shop_query.OrderGoods.Goods).
			Where(
				shop_query.OrderGoods.OrderNo.Eq(genShopOrderNo(uro.OrderNo.String())),
				shop_query.OrderGoods.RefundStatus.Eq(3), //退款完成
			).Find()
		if shopErr != nil {
			return shopErr
		}
		if len(refundedGoods) == 0 {
			continue
		}
		//查询ERP未退款的商品
		unRefundedGoods, erpErr := erp_query.ErpSyncOrderGoods.GetNotRefundedOrderGoods(uro.OrderNo.String())
		if erpErr != nil {
			return erpErr
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
				if _, er := erp_query.ErpSyncOrderGoods.UpdateStatus(-1, uro.OrderNo.String(), rog.Goods.GoodsErpSpid); er != nil {
					log.Errorf("syncOrderStatus UpdateStatus err:%s", er)
					return er
				}
			}
		}
	}

	return nil
}
