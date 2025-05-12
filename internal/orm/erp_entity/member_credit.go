package erp_entity

type MemberCredit struct {
	ErpUID string  `db:"erp_uid" gorm:"column:erp_uid" json:"erp_uid"`
	Money  float64 `db:"money" gorm:"column:money" json:"money"`
	Limit  float64 `db:"limit" gorm:"column:limit" json:"limit"`
}
