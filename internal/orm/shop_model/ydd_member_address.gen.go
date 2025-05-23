// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_model

import (
	"encoding/json"
)

const TableNameMemberAddress = "ydd_member_address"

// MemberAddress 用户地址管理
type MemberAddress struct {
	ID          int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MemberID    int32  `gorm:"column:member_id;not null;comment:会员id" json:"member_id"`                 // 会员id
	SiteID      int32  `gorm:"column:site_id;not null;comment:站点id" json:"site_id"`                     // 站点id
	Name        string `gorm:"column:name;not null;comment:用户姓名" json:"name"`                           // 用户姓名
	Mobile      string `gorm:"column:mobile;not null;comment:手机" json:"mobile"`                         // 手机
	Telephone   string `gorm:"column:telephone;not null;comment:联系电话" json:"telephone"`                 // 联系电话
	ProvinceID  int32  `gorm:"column:province_id;not null;comment:省id" json:"province_id"`              // 省id
	CityID      int32  `gorm:"column:city_id;not null;comment:市id" json:"city_id"`                      // 市id
	DistrictID  int32  `gorm:"column:district_id;not null;comment:区县id" json:"district_id"`             // 区县id
	CommunityID int32  `gorm:"column:community_id;not null;comment:社区id" json:"community_id"`           // 社区id
	Address     string `gorm:"column:address;not null;comment:地址信息" json:"address"`                     // 地址信息
	FullAddress string `gorm:"column:full_address;not null;comment:详细地址信息" json:"full_address"`         // 详细地址信息
	Longitude   string `gorm:"column:longitude;not null;comment:经度" json:"longitude"`                   // 经度
	Latitude    string `gorm:"column:latitude;not null;comment:纬度" json:"latitude"`                     // 纬度
	IsDefault   int32  `gorm:"column:is_default;not null;comment:是否是默认地址" json:"is_default"`            // 是否是默认地址
	Type        int32  `gorm:"column:type;not null;default:1;comment:地址类型  1 普通地址  2 定位地址" json:"type"` // 地址类型  1 普通地址  2 定位地址
}

// MarshalBinary 支持json序列化
func (m *MemberAddress) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *MemberAddress) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName MemberAddress's table name
func (*MemberAddress) TableName() string {
	return TableNameMemberAddress
}
