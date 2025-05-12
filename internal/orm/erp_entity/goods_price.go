package erp_entity

type GoodsPrice struct {
	Price        float64 `db:"price" gorm:"column:price" json:"price"`
	CostPrice    float64 `db:"cost_price" gorm:"column:cost_price" json:"cost_price"`
	MarketPrice  float64 `db:"market_price" gorm:"column:market_price" json:"market_price"`
	GoodsErpSpid string  `db:"goods_erp_spid" gorm:"column:goods_erp_spid" json:"goods_erp_spid"`
}
