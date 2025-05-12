package erp_entity

type Invoice struct {
	InvoiceID     UTF8String `db:"invoice_id" gorm:"column:invoice_id" json:"invoice_id"`
	InvoiceCode   UTF8String `db:"invoice_code" gorm:"column:invoice_code" json:"invoice_code"`
	InvoiceNo     UTF8String `db:"invoice_no" gorm:"column:invoice_no" json:"invoice_no"`
	SecurityCode  UTF8String `db:"security_code" gorm:"column:security_code" json:"security_code"`
	InvoiceImgURL UTF8String `db:"invoice_img_url" gorm:"column:invoice_img_url" json:"invoice_img_url"`
	InvoicePDFURL UTF8String `db:"invoice_pdf_url" gorm:"column:invoice_pdf_url" json:"invoice_pdf_url"`
	ErpUID        UTF8String `db:"erp_uid" gorm:"column:erp_uid" json:"erp_uid"`
	OrderNo       UTF8String `db:"order_no" gorm:"column:order_no" json:"order_no"`
	CreateTime    UTF8String `db:"create_time" gorm:"column:create_time" json:"create_time"`
}
