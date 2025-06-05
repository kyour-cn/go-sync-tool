package erp_entity

type NewMemberQualification struct {
	MemberID            int        `db:"member_id" gorm:"column:member_id;not null" json:"member_id"` // 用户名
	Name                UTF8String `db:"name" gorm:"column:name;size:255;not null" json:"name"`
	Identify            UTF8String `db:"identify" gorm:"column:identify;size:255;not null" json:"identify"`                                // 资质标识
	ExpirationStartDate string     `db:"expiration_start_date" gorm:"column:expiration_start_date;type:date" json:"expiration_start_date"` // 开始时间
	ExpirationEndDate   string     `db:"expiration_end_date" gorm:"column:expiration_end_date;type:date" json:"expiration_end_date"`       // 结束时间
	LongTerm            int        `db:"long_term" gorm:"column:long_term;type:tinyint(1);default:0" json:"long_term"`                     // 是否长期有效
	Image               UTF8String `db:"image" gorm:"column:image;size:255;not null;default:''" json:"image"`                              // 图片地址
	CardNo              UTF8String `db:"card_no" gorm:"column:card_no;size:255;not null;default:''" json:"card_no"`
	CustomForm          UTF8String `db:"custom_form" gorm:"column:custom_form;size:255;not null;default:''" json:"custom_form"`
	BusinessScope       UTF8String `db:"business_scope" gorm:"column:business_scope" json:"business_scope"`
}
