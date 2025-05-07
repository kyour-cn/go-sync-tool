package erp_entity

type MemberBusinessScope struct {
	ID             UTF8String `gorm:"column:id;not null" db:"id" json:"id"`
	ErpUID         UTF8String `gorm:"column:erp_uid;not null" db:"erp_uid" json:"erp_uid"`
	UserBusinessID UTF8String `gorm:"column:user_business_id" db:"user_business_id" json:"user_business_id"`
	UserBusiness   UTF8String `gorm:"column:user_business;not null" db:"user_business" json:"user_business"`
}
