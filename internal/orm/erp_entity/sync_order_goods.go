package erp_entity

type SyncOrderGoods struct {
	OrderNo       UTF8String `db:"order_no" gorm:"column:order_no" json:"order_no"`
	GoodsErpSpid  UTF8String `db:"goods_erp_spid" gorm:"column:goods_erp_spid" json:"goods_erp_spid"`
	OrderGoodsId  int32      `db:"order_goods_id" gorm:"column:order_goods_id" json:"order_goods_id"`
	Price         float64    `db:"price" gorm:"column:price" json:"price"`
	OriginalPrice float64    `db:"original_price" gorm:"column:original_price" json:"original_price"`
	Num           int32      `db:"num" gorm:"column:num" json:"num"`
	GoodsMoney    float64    `db:"goods_money" gorm:"column:goods_money" json:"goods_money"`
	Status        int64      `db:"status" gorm:"column:status" json:"status"`
}
