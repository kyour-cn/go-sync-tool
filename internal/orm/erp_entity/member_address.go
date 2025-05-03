package erp_entity

type MemberAddress struct {
    ID       string     `db:"id" gorm:"column:id" json:"id"`
    RealName UTF8String `db:"realname" gorm:"column:realname" json:"realname"`
    Address  UTF8String `db:"address" gorm:"column:address" json:"address"`
    Mobile   string     `db:"mobile" gorm:"column:mobile" json:"mobile"`
    ErpUID   string     `db:"erp_uid" gorm:"column:erp_uid" json:"erp_uid"`
    Province UTF8String `db:"province" gorm:"column:province" json:"province"`
    City     UTF8String `db:"city" gorm:"column:city" json:"city"`
    District UTF8String `db:"district" gorm:"column:district" json:"district"`
}
