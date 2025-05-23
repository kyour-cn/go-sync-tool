// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_model

import (
	"encoding/json"
)

const TableNameGoodsSku = "ydd_goods_sku"

// GoodsSku 商品表
type GoodsSku struct {
	SkuID              int32   `gorm:"column:sku_id;primaryKey;autoIncrement:true;comment:商品sku_id" json:"sku_id"`                                  // 商品sku_id
	SiteID             int32   `gorm:"column:site_id;not null;comment:所属店铺id" json:"site_id"`                                                       // 所属店铺id
	GoodsID            int32   `gorm:"column:goods_id;not null;comment:商品id" json:"goods_id"`                                                       // 商品id
	SkuName            string  `gorm:"column:sku_name;not null;comment:商品sku名称" json:"sku_name"`                                                    // 商品sku名称
	SkuNo              string  `gorm:"column:sku_no;not null;comment:商品sku编码" json:"sku_no"`                                                        // 商品sku编码
	SkuSpecFormat      string  `gorm:"column:sku_spec_format;comment:sku规格格式" json:"sku_spec_format"`                                               // sku规格格式
	Price              float64 `gorm:"column:price;not null;default:0.00;comment:sku单价" json:"price"`                                               // sku单价
	MarketPrice        float64 `gorm:"column:market_price;not null;default:0.00;comment:sku划线价" json:"market_price"`                                // sku划线价
	CostPrice          float64 `gorm:"column:cost_price;not null;default:0.00;comment:sku成本价" json:"cost_price"`                                    // sku成本价
	DiscountPrice      float64 `gorm:"column:discount_price;not null;default:0.00;comment:sku折扣价（默认等于单价）" json:"discount_price"`                    // sku折扣价（默认等于单价）
	PromotionType      int32   `gorm:"column:promotion_type;not null;comment:活动类型1.限时折扣" json:"promotion_type"`                                     // 活动类型1.限时折扣
	StartTime          int32   `gorm:"column:start_time;not null;comment:活动开始时间" json:"start_time"`                                                 // 活动开始时间
	EndTime            int32   `gorm:"column:end_time;not null;comment:活动结束时间" json:"end_time"`                                                     // 活动结束时间
	Stock              int32   `gorm:"column:stock;not null;comment:商品sku库存" json:"stock"`                                                          // 商品sku库存
	Weight             float64 `gorm:"column:weight;not null;default:0.000;comment:重量（单位g）" json:"weight"`                                          // 重量（单位g）
	Volume             float64 `gorm:"column:volume;not null;default:0.000;comment:体积（单位立方米）" json:"volume"`                                        // 体积（单位立方米）
	ClickNum           int32   `gorm:"column:click_num;not null;comment:点击量" json:"click_num"`                                                      // 点击量
	SaleNum            int32   `gorm:"column:sale_num;not null;comment:销量" json:"sale_num"`                                                         // 销量
	CollectNum         int32   `gorm:"column:collect_num;not null;comment:收藏量" json:"collect_num"`                                                  // 收藏量
	SkuImage           string  `gorm:"column:sku_image;not null;comment:sku主图" json:"sku_image"`                                                    // sku主图
	SkuImages          string  `gorm:"column:sku_images;not null;comment:sku图片" json:"sku_images"`                                                  // sku图片
	GoodsClass         int32   `gorm:"column:goods_class;not null;default:1;comment:商品种类1.实物商品2.虚拟商品3.卡券商品" json:"goods_class"`                     // 商品种类1.实物商品2.虚拟商品3.卡券商品
	GoodsClassName     string  `gorm:"column:goods_class_name;not null;comment:商品种类" json:"goods_class_name"`                                       // 商品种类
	GoodsAttrClass     int32   `gorm:"column:goods_attr_class;not null;default:1;comment:商品类型id" json:"goods_attr_class"`                           // 商品类型id
	GoodsAttrName      string  `gorm:"column:goods_attr_name;not null;comment:商品类型名称" json:"goods_attr_name"`                                       // 商品类型名称
	GoodsName          string  `gorm:"column:goods_name;not null;comment:商品名称" json:"goods_name"`                                                   // 商品名称
	GoodsContent       string  `gorm:"column:goods_content;comment:商品详情" json:"goods_content"`                                                      // 商品详情
	GoodsState         int32   `gorm:"column:goods_state;not null;comment:商品状态（1.正常0下架）" json:"goods_state"`                                        // 商品状态（1.正常0下架）
	GoodsStockAlarm    int32   `gorm:"column:goods_stock_alarm;not null;comment:库存预警" json:"goods_stock_alarm"`                                     // 库存预警
	IsVirtual          int32   `gorm:"column:is_virtual;not null;comment:是否虚拟类商品（0实物1.虚拟）" json:"is_virtual"`                                       // 是否虚拟类商品（0实物1.虚拟）
	VirtualIndate      int32   `gorm:"column:virtual_indate;not null;default:1;comment:虚拟商品有效期" json:"virtual_indate"`                              // 虚拟商品有效期
	IsFreeShipping     int32   `gorm:"column:is_free_shipping;not null;comment:是否免邮" json:"is_free_shipping"`                                       // 是否免邮
	ShippingTemplate   int32   `gorm:"column:shipping_template;not null;comment:指定运费模板" json:"shipping_template"`                                   // 指定运费模板
	GoodsSpecFormat    string  `gorm:"column:goods_spec_format;comment:商品规格格式" json:"goods_spec_format"`                                            // 商品规格格式
	GoodsAttrFormat    string  `gorm:"column:goods_attr_format;comment:商品属性格式" json:"goods_attr_format"`                                            // 商品属性格式
	IsDelete           int32   `gorm:"column:is_delete;not null;comment:是否已经删除" json:"is_delete"`                                                   // 是否已经删除
	Introduction       string  `gorm:"column:introduction;not null;comment:促销语" json:"introduction"`                                                // 促销语
	Keywords           string  `gorm:"column:keywords;not null;comment:关键词" json:"keywords"`                                                        // 关键词
	Unit               string  `gorm:"column:unit;not null;comment:单位" json:"unit"`                                                                 // 单位
	Sort               int32   `gorm:"column:sort;not null;comment:排序" json:"sort"`                                                                 // 排序
	CreateTime         uint    `gorm:"column:create_time;not null;autoCreateTime;comment:创建时间" json:"create_time"`                                  // 创建时间
	ModifyTime         int32   `gorm:"column:modify_time;not null;comment:修改时间" json:"modify_time"`                                                 // 修改时间
	VideoURL           string  `gorm:"column:video_url;not null;comment:视频" json:"video_url"`                                                       // 视频
	Evaluate           int32   `gorm:"column:evaluate;not null;comment:评价数" json:"evaluate"`                                                        // 评价数
	EvaluateShaitu     int32   `gorm:"column:evaluate_shaitu;not null;comment:晒图评价数" json:"evaluate_shaitu"`                                        // 晒图评价数
	EvaluateShipin     int32   `gorm:"column:evaluate_shipin;not null;comment:视频评价数" json:"evaluate_shipin"`                                        // 视频评价数
	EvaluateZhuiping   int32   `gorm:"column:evaluate_zhuiping;not null;comment:追评数" json:"evaluate_zhuiping"`                                      // 追评数
	EvaluateHaoping    int32   `gorm:"column:evaluate_haoping;not null;comment:好评数" json:"evaluate_haoping"`                                        // 好评数
	EvaluateZhongping  int32   `gorm:"column:evaluate_zhongping;not null;comment:中评数" json:"evaluate_zhongping"`                                    // 中评数
	EvaluateChaping    int32   `gorm:"column:evaluate_chaping;not null;comment:差评数" json:"evaluate_chaping"`                                        // 差评数
	SpecName           string  `gorm:"column:spec_name;not null;comment:规格名称" json:"spec_name"`                                                     // 规格名称
	SupplierID         int32   `gorm:"column:supplier_id;not null;comment:供应商id" json:"supplier_id"`                                                // 供应商id
	IsConsumeDiscount  bool    `gorm:"column:is_consume_discount;not null;comment:是否参与会员等级折扣" json:"is_consume_discount"`                           // 是否参与会员等级折扣
	DiscountConfig     bool    `gorm:"column:discount_config;not null;comment:优惠设置（0默认 1自定义）" json:"discount_config"`                               // 优惠设置（0默认 1自定义）
	DiscountMethod     string  `gorm:"column:discount_method;not null;comment:优惠方式（discount打折 manjian 满减 fixed_price 指定价格）" json:"discount_method"` // 优惠方式（discount打折 manjian 满减 fixed_price 指定价格）
	MemberPrice        string  `gorm:"column:member_price;not null;comment:会员价" json:"member_price"`                                                // 会员价
	GoodsServiceIds    string  `gorm:"column:goods_service_ids;not null;comment:商品服务id" json:"goods_service_ids"`                                   // 商品服务id
	VirtualSale        int32   `gorm:"column:virtual_sale;not null;comment:虚拟销量" json:"virtual_sale"`                                               // 虚拟销量
	MaxBuy             int32   `gorm:"column:max_buy;not null;comment:限购" json:"max_buy"`                                                           // 限购
	MinBuy             int32   `gorm:"column:min_buy;not null" json:"min_buy"`
	RecommendWay       int32   `gorm:"column:recommend_way;not null;comment:推荐方式，1：新品，2：精品，3；推荐" json:"recommend_way"`   // 推荐方式，1：新品，2：精品，3；推荐
	FenxiaoPrice       float64 `gorm:"column:fenxiao_price;not null;default:0.00;comment:分销计算价格" json:"fenxiao_price"`   // 分销计算价格
	StockAlarm         int32   `gorm:"column:stock_alarm;not null;comment:sku库存预警" json:"stock_alarm"`                   // sku库存预警
	SaleSort           int32   `gorm:"column:sale_sort;not null;comment:销量排序字段 占位用" json:"sale_sort"`                    // 销量排序字段 占位用
	IsDefault          bool    `gorm:"column:is_default;not null;comment:是否默认" json:"is_default"`                        // 是否默认
	VerifyNum          int32   `gorm:"column:verify_num;not null;comment:核销次数" json:"verify_num"`                        // 核销次数
	IsLimit            int32   `gorm:"column:is_limit;not null;comment:是否限购(0否1是)" json:"is_limit"`                      // 是否限购(0否1是)
	LimitType          int32   `gorm:"column:limit_type;not null;default:1;comment:限购类型(1单次限购2长期限购)" json:"limit_type"`  // 限购类型(1单次限购2长期限购)
	QrID               int32   `gorm:"column:qr_id;not null;comment:社群二维码id" json:"qr_id"`                               // 社群二维码id
	TemplateID         int32   `gorm:"column:template_id;not null;comment:海报id" json:"template_id"`                      // 海报id
	SuccessEvaluateNum int32   `gorm:"column:success_evaluate_num;not null;comment:评价审核通过数" json:"success_evaluate_num"` // 评价审核通过数
	FailEvaluateNum    int32   `gorm:"column:fail_evaluate_num;not null;comment:评价审核失败数" json:"fail_evaluate_num"`       // 评价审核失败数
	WaitEvaluateNum    int32   `gorm:"column:wait_evaluate_num;not null;comment:评价待审核数" json:"wait_evaluate_num"`        // 评价待审核数
	GoodsArea          string  `gorm:"column:goods_area;comment:商品库区" json:"goods_area"`                                 // 商品库区
	IsPresale          int32   `gorm:"column:is_presale;not null;comment:是否预售" json:"is_presale"`                        // 是否预售
	TaxRate            float64 `gorm:"column:tax_rate;not null;default:0.00;comment:税率 单位:百分比" json:"tax_rate"`          // 税率 单位:百分比
	LiveGoodsID        int32   `gorm:"column:live_goods_id;not null;comment:是否同步到直播" json:"live_goods_id"`               // 是否同步到直播
}

// MarshalBinary 支持json序列化
func (m *GoodsSku) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *GoodsSku) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName GoodsSku's table name
func (*GoodsSku) TableName() string {
	return TableNameGoodsSku
}
