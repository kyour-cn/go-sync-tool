package erp_entity

type Salesman struct {
	Realname UTF8String `db:"realname" gorm:"column:realname" json:"realname"`
	Mobile   UTF8String `db:"mobile" gorm:"column:mobile" json:"mobile"`
	SaleID   string     `db:"saleer_id" gorm:"column:saleer_id" json:"sale_id"`
}
