package erp_entity

type NewMember struct {
	MemberID      int64      `db:"member_id" gorm:"column:member_id" json:"member_id"`
	Nickname      UTF8String `db:"nickname" gorm:"column:nickname" json:"nickname"`
	Realname      UTF8String `db:"realname" gorm:"column:realname" json:"realname"`
	Username      UTF8String `db:"username" gorm:"column:username" json:"username"`
	Mobile        UTF8String `db:"mobile" gorm:"column:mobile" json:"mobile"`
	MemberType    UTF8String `db:"member_type" gorm:"column:member_type" json:"member_type"`
	SyncTime      UTF8String `db:"sync_time" gorm:"column:sync_time" json:"sync_time"`
	Province      UTF8String `db:"province" gorm:"column:province" json:"province"`
	City          UTF8String `db:"city" gorm:"column:city" json:"city"`
	District      UTF8String `db:"district" gorm:"column:district" json:"district"`
	Address       UTF8String `db:"address" gorm:"column:address" json:"address"`
	BusinessScope UTF8String `db:"business_scope" gorm:"column:business_scope" json:"business_scope"`
}
