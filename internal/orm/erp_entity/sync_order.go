package erp_entity

import "time"

type SyncOrder struct {
	OrderNo      string     `db:"order_no" gorm:"column:order_no" json:"order_no"`
	ErpUID       string     `db:"erp_uid" gorm:"column:erp_uid" json:"erp_uid"`
	PayTime      time.Time  `db:"pay_time" gorm:"column:pay_time" json:"pay_time"`
	PayMoney     float64    `db:"pay_money" gorm:"column:pay_money" json:"pay_money"`
	OrderMoney   float64    `db:"order_money" gorm:"column:order_money" json:"order_money"`
	PayType      UTF8String `db:"pay_type" gorm:"column:pay_type" json:"pay_type"`
	Name         UTF8String `db:"name" gorm:"column:name" json:"name"`
	Mobile       UTF8String `db:"mobile" gorm:"column:mobile" json:"mobile"`
	Address      UTF8String `db:"address" gorm:"column:address" json:"address"`
	FullAddress  UTF8String `db:"full_address" gorm:"column:full_address" json:"full_address"`
	BuyerMessage UTF8String `db:"buyer_message" gorm:"column:buyer_message" json:"buyer_message"`
	OrderStatus  int32      `db:"order_status" gorm:"column:order_status" json:"order_status"`
	MemberID     int64      `db:"member_id" gorm:"column:member_id" json:"member_id"`
	CouponMoney  float64    `db:"coupon_money" gorm:"column:coupon_money" json:"coupon_money"`
	Postage      float64    `db:"postage" gorm:"column:postage" json:"postage"`
	SaleerID     UTF8String `db:"saleer_id" gorm:"column:saleer_id" json:"saleer_id"`
}
