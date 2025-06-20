// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_model

import (
	"encoding/json"
)

const TableNameGoodsCategory = "ydd_goods_category"

// GoodsCategory  商品分类
type GoodsCategory struct {
	CategoryID       int32   `gorm:"column:category_id;primaryKey;autoIncrement:true" json:"category_id"`
	SiteID           int32   `gorm:"column:site_id;not null;comment:站点id" json:"site_id"`                               // 站点id
	CategoryName     string  `gorm:"column:category_name;not null;comment:分类名称" json:"category_name"`                   // 分类名称
	ShortName        string  `gorm:"column:short_name;not null;comment:简称" json:"short_name"`                           // 简称
	Pid              int32   `gorm:"column:pid;not null;comment:分类上级" json:"pid"`                                       // 分类上级
	Level            int32   `gorm:"column:level;not null;comment:层级" json:"level"`                                     // 层级
	IsShow           int32   `gorm:"column:is_show;not null;comment:是否显示（0显示  -1不显示）" json:"is_show"`                   // 是否显示（0显示  -1不显示）
	Sort             int32   `gorm:"column:sort;not null;comment:排序" json:"sort"`                                       // 排序
	Image            string  `gorm:"column:image;not null;comment:分类图片" json:"image"`                                   // 分类图片
	Keywords         string  `gorm:"column:keywords;not null;comment:分类页面关键字" json:"keywords"`                          // 分类页面关键字
	Description      string  `gorm:"column:description;not null;comment:分类介绍" json:"description"`                       // 分类介绍
	AttrClassID      int32   `gorm:"column:attr_class_id;not null;comment:关联商品类型id" json:"attr_class_id"`               // 关联商品类型id
	AttrClassName    string  `gorm:"column:attr_class_name;not null;comment:关联商品类型名称" json:"attr_class_name"`           // 关联商品类型名称
	CategoryId1      int32   `gorm:"column:category_id_1;not null;comment:一级分类id" json:"category_id_1"`                 // 一级分类id
	CategoryId2      int32   `gorm:"column:category_id_2;not null;comment:二级分类id" json:"category_id_2"`                 // 二级分类id
	CategoryId3      int32   `gorm:"column:category_id_3;not null;comment:三级分类id" json:"category_id_3"`                 // 三级分类id
	CategoryFullName string  `gorm:"column:category_full_name;not null;comment:组装名称" json:"category_full_name"`         // 组装名称
	ImageAdv         string  `gorm:"column:image_adv;not null;comment:分类广告图" json:"image_adv"`                          // 分类广告图
	CommissionRate   float64 `gorm:"column:commission_rate;not null;default:0.00;comment:佣金比率%" json:"commission_rate"` // 佣金比率%
	IsGoods          int32   `gorm:"column:is_goods;not null;comment:是否渲染商品" json:"is_goods"`                           // 是否渲染商品
}

// MarshalBinary 支持json序列化
func (m *GoodsCategory) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *GoodsCategory) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName GoodsCategory's table name
func (*GoodsCategory) TableName() string {
	return TableNameGoodsCategory
}
