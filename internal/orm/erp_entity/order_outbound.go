package erp_entity

type OrderOutBound struct {
    OutboundNo       UTF8String `db:"outbound_no" gorm:"column:outbound_no" json:"outbound_no"`
    OrderNo          UTF8String `db:"order_no" gorm:"column:order_no" json:"order_no"`
    GoodsErpSpid     UTF8String `db:"goods_erp_spid" gorm:"column:goods_erp_spid" json:"goods_erp_spid"`
    OutboundTime     UTF8String `db:"outbound_time" gorm:"column:outbound_time" json:"outbound_time"`
    Validity         UTF8String `db:"validity" gorm:"column:validity" json:"validity"`
    OutboundNum      int32      `db:"outbound_num" gorm:"column:outbound_num" json:"outbound_num"`
    OutboundPrice    float64    `db:"outbound_price" gorm:"column:outbound_price" json:"outbound_price"`
    LogisticsCompany UTF8String `db:"logistics_company" gorm:"column:logistics_company" json:"logistics_company"`
    LogisticsCode    UTF8String `db:"logistics_code" gorm:"column:logistics_code" json:"logistics_code"`
}
