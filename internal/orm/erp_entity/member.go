package erp_entity

type Member struct {
	Username     UTF8String `db:"username" gorm:"column:username" json:"username"`
	Nickname     UTF8String `db:"nickname" gorm:"column:nickname" json:"nickname"`
	Mobile       string     `db:"mobile" gorm:"column:mobile" json:"mobile"`
	ErpUID       string     `db:"erp_uid" gorm:"column:erp_uid" json:"erp_uid"`
	Province     UTF8String `db:"province" gorm:"column:province" json:"province"`
	City         UTF8String `db:"city" gorm:"column:city" json:"city"`
	District     UTF8String `db:"district" gorm:"column:district" json:"district"`
	FullAddress  UTF8String `db:"full_address" gorm:"column:full_address" json:"full_address"`
	SaleerID     string     `db:"saleer_id" gorm:"column:saleer_id" json:"saleer_id"`
	MemberType   UTF8String `db:"member_type" gorm:"column:member_type" json:"member_type"`
	MemberID     string     `db:"member_id" gorm:"column:member_id" json:"member_id"`
	ScopeControl int32      `db:"scope_control" gorm:"column:scope_control" json:"scope_control"`
}
