package erp_entity

type GoodsStock struct {
	GoodsStock   float64 `db:"goods_stock" gorm:"column:goods_stock" json:"goods_stock"`
	GoodsErpSpid string  `db:"goods_erp_spid" gorm:"column:goods_erp_spid" json:"goods_erp_spid"`
}
