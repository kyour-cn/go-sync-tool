package erp_entity

type Goods struct {
    GoodsName          UTF8String `db:"goods_name" gorm:"column:goods_name" json:"goods_name"`
    AttrValidity       UTF8String `db:"attr_validity" gorm:"column:attr_validity" json:"attr_validity"`
    AttrSpecs          UTF8String `db:"attr_specs" gorm:"column:attr_specs" json:"attr_specs"`
    GoodsSpecs         UTF8String `db:"goods_specs" gorm:"column:goods_specs" json:"goods_specs"`
    AttrShelfLife      UTF8String `db:"attr_shelf_life" gorm:"column:attr_shelf_life" json:"attr_shelf_life"`
    GoodsBatch         UTF8String `db:"goods_batch" gorm:"column:goods_batch" json:"goods_batch"`
    Unit               UTF8String `db:"unit" gorm:"column:unit" json:"unit"`
    GoodsErpSpid       string     `db:"goods_erp_spid" gorm:"column:goods_erp_spid" json:"goods_erp_spid"`
    AttrFactory        UTF8String `db:"attr_factory" gorm:"column:attr_factory" json:"attr_factory"`
    Place              UTF8String `db:"place" gorm:"column:place" json:"place"`
    AttrApprovalNumber UTF8String `db:"attr_approval_number" gorm:"column:attr_approval_number" json:"attr_approval_number"`
    AttrDosageForm     UTF8String `db:"attr_dosage_form" gorm:"column:attr_dosage_form" json:"attr_dosage_form"`
    AttrCountryCode    UTF8String `db:"attr_country_code" gorm:"column:attr_country_code" json:"attr_country_code"`
    GoodsArea          UTF8String `db:"goods_area" gorm:"column:goods_area" json:"goods_area"`
    BarCode            string     `db:"bar_code" gorm:"column:bar_code" json:"bar_code"`
    MediumPackageNum   int        `db:"medium_package_num" gorm:"column:medium_package_num" json:"medium_package_num"`
    BusinessTypeID     UTF8String `db:"business_type_id" gorm:"column:business_type_id" json:"business_type_id"`
    BusinessTypeName   UTF8String `db:"business_type_name" gorm:"column:business_type_name" json:"business_type_name"`
    FactoryDate        UTF8String `db:"factory_date" gorm:"column:factory_date" json:"factory_date"`
    GoodsNickname      UTF8String `db:"goods_nickname" gorm:"column:goods_nickname" json:"goods_nickname"`
    BuyMinNum          int        `db:"buy_min_num" gorm:"column:buy_min_num" json:"buy_min_num"`
    BuyMaxNum          int        `db:"buy_max_num" gorm:"column:buy_max_num" json:"buy_max_num"`
    IsMedicinal        int        `db:"is_medicinal" gorm:"column:is_medicinal" json:"is_medicinal"`
    GoodsNo            string     `db:"goods_no" gorm:"column:goods_no" json:"goods_no"`
    IsPrescription     UTF8String `db:"is_prescription" gorm:"column:is_prescription" json:"is_prescription"`
    YiBaoType          UTF8String `db:"yibao_type" gorm:"column:yibao_type" json:"yibao_type"`
    YiBaoNo            UTF8String `db:"yibao_no" gorm:"column:yibao_no" json:"yibao_no"`
    IsJc               int32      `gorm:"column:is_jc;not null;comment:是否集采，0=否 1=是" json:"is_jc"`                             // 是否中药，0=否 1=是
    DzjgCode           int32      `gorm:"column:dzjg_code;not null;comment:是否中药，0=否 1=是" json:"dzjg_code"`                     // 是否中药，0=否 1=是
    TraceabilityCode   int32      `gorm:"column:traceability_code;not null;comment:是否有追溯码，0=否 1=是" json:"traceability_code"` // 是否中药，0=否 1=是
}
