// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_model

import (
	"encoding/json"
)

const TableNameErpInvoice = "ydd_erp_invoice"

// ErpInvoice ERP同步的发票数据
type ErpInvoice struct {
	ID               int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	InvoiceID        string `gorm:"column:invoice_id;not null;comment:发票唯一id" json:"invoice_id"`                                    // 发票唯一id
	InvoiceCode      string `gorm:"column:invoice_code;not null;comment:发票代码" json:"invoice_code"`                                  // 发票代码
	InvoiceNo        string `gorm:"column:invoice_no;not null;comment:发票号码" json:"invoice_no"`                                      // 发票号码
	SecurityCode     string `gorm:"column:security_code;not null;comment:发票防伪码" json:"security_code"`                               // 发票防伪码
	InvoiceImgURL    string `gorm:"column:invoice_img_url;not null;comment:发票图片的URL，与pdf至少必传一个" json:"invoice_img_url"`             // 发票图片的URL，与pdf至少必传一个
	InvoicePdfURL    string `gorm:"column:invoice_pdf_url;not null;comment:发票PDF的URL，与图片至少必传一个" json:"invoice_pdf_url"`             // 发票PDF的URL，与图片至少必传一个
	MemberID         int32  `gorm:"column:member_id;not null;comment:客户ERPID" json:"member_id"`                                     // 客户ERPID
	OrderNo          string `gorm:"column:order_no;not null;comment:小程序订单编号" json:"order_no"`                                       // 小程序订单编号
	OriginCreateTime string `gorm:"column:origin_create_time;not null;comment:开票时间 如2006-01-02 03:04:05" json:"origin_create_time"` // 开票时间 如2006-01-02 03:04:05
	CreateTime       uint   `gorm:"column:create_time;not null;autoCreateTime" json:"create_time"`
}

// MarshalBinary 支持json序列化
func (m *ErpInvoice) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *ErpInvoice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName ErpInvoice's table name
func (*ErpInvoice) TableName() string {
	return TableNameErpInvoice
}
