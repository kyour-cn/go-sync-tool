// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"app/internal/orm/shop_model"
)

func newMember(db *gorm.DB, opts ...gen.DOOption) member {
	_member := member{}

	_member.memberDo.UseDB(db, opts...)
	_member.memberDo.UseModel(&shop_model.Member{})

	tableName := _member.memberDo.TableName()
	_member.ALL = field.NewAsterisk(tableName)
	_member.MemberID = field.NewInt32(tableName, "member_id")
	_member.SiteID = field.NewInt32(tableName, "site_id")
	_member.SourceMember = field.NewInt32(tableName, "source_member")
	_member.FenxiaoID = field.NewInt32(tableName, "fenxiao_id")
	_member.IsFenxiao = field.NewBool(tableName, "is_fenxiao")
	_member.Username = field.NewString(tableName, "username")
	_member.Nickname = field.NewString(tableName, "nickname")
	_member.Mobile = field.NewString(tableName, "mobile")
	_member.Email = field.NewString(tableName, "email")
	_member.Password = field.NewString(tableName, "password")
	_member.Status = field.NewInt32(tableName, "status")
	_member.Headimg = field.NewString(tableName, "headimg")
	_member.MemberLevel = field.NewInt32(tableName, "member_level")
	_member.MemberLevelName = field.NewString(tableName, "member_level_name")
	_member.MemberLabel = field.NewInt32(tableName, "member_label")
	_member.MemberLabelName = field.NewString(tableName, "member_label_name")
	_member.Qq = field.NewString(tableName, "qq")
	_member.QqOpenid = field.NewString(tableName, "qq_openid")
	_member.WxOpenid = field.NewString(tableName, "wx_openid")
	_member.WeappOpenid = field.NewString(tableName, "weapp_openid")
	_member.WxUnionid = field.NewString(tableName, "wx_unionid")
	_member.AliOpenid = field.NewString(tableName, "ali_openid")
	_member.BaiduOpenid = field.NewString(tableName, "baidu_openid")
	_member.ToutiaoOpenid = field.NewString(tableName, "toutiao_openid")
	_member.DouyinOpenid = field.NewString(tableName, "douyin_openid")
	_member.LoginIP = field.NewString(tableName, "login_ip")
	_member.LoginType = field.NewString(tableName, "login_type")
	_member.LoginTime = field.NewInt32(tableName, "login_time")
	_member.LastLoginIP = field.NewString(tableName, "last_login_ip")
	_member.LastLoginType = field.NewString(tableName, "last_login_type")
	_member.LastLoginTime = field.NewInt32(tableName, "last_login_time")
	_member.LoginNum = field.NewInt32(tableName, "login_num")
	_member.Realname = field.NewString(tableName, "realname")
	_member.Sex = field.NewInt32(tableName, "sex")
	_member.Location = field.NewString(tableName, "location")
	_member.Birthday = field.NewInt32(tableName, "birthday")
	_member.RegTime = field.NewInt32(tableName, "reg_time")
	_member.Point = field.NewFloat64(tableName, "point")
	_member.Balance = field.NewFloat64(tableName, "balance")
	_member.Growth = field.NewFloat64(tableName, "growth")
	_member.BalanceMoney = field.NewFloat64(tableName, "balance_money")
	_member.Account5 = field.NewFloat64(tableName, "account5")
	_member.IsAuth = field.NewInt32(tableName, "is_auth")
	_member.SignTime = field.NewInt32(tableName, "sign_time")
	_member.SignDaysSeries = field.NewInt32(tableName, "sign_days_series")
	_member.PayPassword = field.NewString(tableName, "pay_password")
	_member.OrderMoney = field.NewFloat64(tableName, "order_money")
	_member.OrderCompleteMoney = field.NewFloat64(tableName, "order_complete_money")
	_member.OrderNum = field.NewInt32(tableName, "order_num")
	_member.OrderCompleteNum = field.NewInt32(tableName, "order_complete_num")
	_member.BalanceWithdrawApply = field.NewFloat64(tableName, "balance_withdraw_apply")
	_member.BalanceWithdraw = field.NewFloat64(tableName, "balance_withdraw")
	_member.IsDelete = field.NewBool(tableName, "is_delete")
	_member.MemberLevelType = field.NewInt32(tableName, "member_level_type")
	_member.LevelExpireTime = field.NewInt32(tableName, "level_expire_time")
	_member.IsEditUsername = field.NewInt32(tableName, "is_edit_username")
	_member.LoginTypeName = field.NewString(tableName, "login_type_name")
	_member.CanReceiveRegistergift = field.NewInt32(tableName, "can_receive_registergift")
	_member.ProvinceID = field.NewInt32(tableName, "province_id")
	_member.CityID = field.NewInt32(tableName, "city_id")
	_member.DistrictID = field.NewInt32(tableName, "district_id")
	_member.FullAddress = field.NewString(tableName, "full_address")
	_member.Address = field.NewString(tableName, "address")
	_member.ErpUID = field.NewString(tableName, "erp_uid")
	_member.ErpCode = field.NewString(tableName, "erp_code")
	_member.ErpName = field.NewString(tableName, "erp_name")
	_member.Gid = field.NewInt32(tableName, "gid")
	_member.QualificationsImages = field.NewString(tableName, "qualifications_images")
	_member.Inputmanid = field.NewInt32(tableName, "inputmanid")
	_member.Deptid = field.NewInt32(tableName, "deptid")
	_member.Employeename = field.NewString(tableName, "employeename")
	_member.ErpAgentid = field.NewInt32(tableName, "erp_agentid")
	_member.Customopcode = field.NewString(tableName, "customopcode")
	_member.SyncTime = field.NewInt32(tableName, "sync_time")
	_member.MemberRemark = field.NewString(tableName, "member_remark")
	_member.Employeeid = field.NewString(tableName, "employeeid")
	_member.InviterID = field.NewInt32(tableName, "inviter_id")
	_member.InviteCode = field.NewString(tableName, "invite_code")
	_member.QrInviteCode = field.NewString(tableName, "qr_invite_code")
	_member.IsNatureofbusiness = field.NewInt32(tableName, "is_natureofbusiness")
	_member.SalesmanID = field.NewInt32(tableName, "salesman_id")
	_member.StaffType = field.NewInt32(tableName, "staff_type")
	_member.FristLoginTime = field.NewInt32(tableName, "frist_login_time")
	_member.RegMobile = field.NewString(tableName, "reg_mobile")
	_member.RegPassword = field.NewString(tableName, "reg_password")
	_member.AppointmenSalesmanID = field.NewInt32(tableName, "appointmen_salesman_id")
	_member.CreditBalance = field.NewFloat64(tableName, "credit_balance")
	_member.CreditLimit = field.NewFloat64(tableName, "credit_limit")
	_member.SessionKey = field.NewString(tableName, "session_key")
	_member.MemberQualification = memberHasManyMemberQualification{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("MemberQualification", "shop_model.MemberQualification"),
	}

	_member.MemberAddress = memberHasManyMemberAddress{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("MemberAddress", "shop_model.MemberAddress"),
	}

	_member.fillFieldMap()

	return _member
}

// member 系统用户表
type member struct {
	memberDo

	ALL                    field.Asterisk
	MemberID               field.Int32   // 主键
	SiteID                 field.Int32   // 站点id
	SourceMember           field.Int32   // 推荐人
	FenxiaoID              field.Int32   // 分销商（分销有效）
	IsFenxiao              field.Bool    // 是否是分销商
	Username               field.String  // 用户名
	Nickname               field.String  // 用户昵称
	Mobile                 field.String  // 手机号
	Email                  field.String  // 邮箱
	Password               field.String  // 用户密码（MD5）
	Status                 field.Int32   // 用户状态  用户状态默认为1
	Headimg                field.String  // 用户头像
	MemberLevel            field.Int32   // 用户等级
	MemberLevelName        field.String  // 会员等级名称
	MemberLabel            field.Int32   // 用户标签
	MemberLabelName        field.String  // 会员标签名称
	Qq                     field.String  // qq号
	QqOpenid               field.String  // qq互联id (已改为客服系统回话token)
	WxOpenid               field.String  // 微信用户openid
	WeappOpenid            field.String  // 微信小程序openid
	WxUnionid              field.String  // 微信unionid
	AliOpenid              field.String  // 支付宝账户id
	BaiduOpenid            field.String  // 百度账户id
	ToutiaoOpenid          field.String  // 头条账号
	DouyinOpenid           field.String  // 抖音小程序openid
	LoginIP                field.String  // 当前登录ip
	LoginType              field.String  // 当前登录的操作终端类型
	LoginTime              field.Int32   // 当前登录时间
	LastLoginIP            field.String  // 上次登录ip
	LastLoginType          field.String  // 上次登录的操作终端类型
	LastLoginTime          field.Int32   // 上次登录时间
	LoginNum               field.Int32   // 登录次数
	Realname               field.String  // 真实姓名
	Sex                    field.Int32   // 性别 0保密 1男 2女
	Location               field.String  // 定位地址
	Birthday               field.Int32   // 出生日期
	RegTime                field.Int32   // 注册时间
	Point                  field.Float64 // 积分
	Balance                field.Float64 // 余额
	Growth                 field.Float64 // 成长值
	BalanceMoney           field.Float64 // 现金余额(可提现)
	Account5               field.Float64 // 账户5（改为下单时间）
	IsAuth                 field.Int32   // 是否认证
	SignTime               field.Int32   // 最后一次签到时间
	SignDaysSeries         field.Int32   // 持续签到天数
	PayPassword            field.String  // 交易密码
	OrderMoney             field.Float64 // 付款后-消费金额
	OrderCompleteMoney     field.Float64 // 订单完成-消费金额
	OrderNum               field.Int32   // 付款后-消费次数
	OrderCompleteNum       field.Int32   // 订单完成-消费次数
	BalanceWithdrawApply   field.Float64 // 提现中余额
	BalanceWithdraw        field.Float64 // 已提现余额
	IsDelete               field.Bool    // 0正常  1已删除
	MemberLevelType        field.Int32   // 会员卡类型 0免费卡 1付费卡
	LevelExpireTime        field.Int32   // 会员卡过期时间
	IsEditUsername         field.Int32   // 是否可修改用户名
	LoginTypeName          field.String  // 登陆类型名称
	CanReceiveRegistergift field.Int32   // 是否可以领取新人礼(只针对后台注册的用户 1可以 0不可以)
	ProvinceID             field.Int32   // 省id（增）
	CityID                 field.Int32   // 市id（增）
	DistrictID             field.Int32   // 区id（增）
	FullAddress            field.String  // 全地址（增）
	Address                field.String  // 详情地址（增）
	ErpUID                 field.String  // erpid（增）
	ErpCode                field.String  // erpcode（增）
	ErpName                field.String  // 客户店名（增）
	Gid                    field.Int32   // 关联账号member_id
	QualificationsImages   field.String  // 资质图片（已废弃）
	Inputmanid             field.Int32   // 维护人ID
	Deptid                 field.Int32   // 维护人部门ID
	Employeename           field.String  // 维护人姓名
	ErpAgentid             field.Int32   // 代理人ID
	Customopcode           field.String  // 客户名称首字母
	SyncTime               field.Int32   // 同步到erp的时间戳 0=未同步
	MemberRemark           field.String  // 后台显示备注（增）
	Employeeid             field.String  // 代理人ID（增）
	InviterID              field.Int32   // 推荐人id（增）
	InviteCode             field.String  // 用户邀请码（增）
	QrInviteCode           field.String  // 邀请二维码地址（增）
	IsNatureofbusiness     field.Int32   // 是否管控 0否 1是
	SalesmanID             field.Int32   // 业务员id
	StaffType              field.Int32   // 员工类型（增）1业务员
	FristLoginTime         field.Int32   // 首次登录时间
	RegMobile              field.String  // 老客户注册时保存（未登录过）
	RegPassword            field.String  // 老客户注册时保存（未登录过）
	AppointmenSalesmanID   field.Int32   // 公海预约业务员id
	CreditBalance          field.Float64 // 资信余额
	CreditLimit            field.Float64 // 资信额度
	SessionKey             field.String  // 微信session_key
	MemberQualification    memberHasManyMemberQualification

	MemberAddress memberHasManyMemberAddress

	fieldMap map[string]field.Expr
}

func (m member) Table(newTableName string) *member {
	m.memberDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m member) As(alias string) *member {
	m.memberDo.DO = *(m.memberDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *member) updateTableName(table string) *member {
	m.ALL = field.NewAsterisk(table)
	m.MemberID = field.NewInt32(table, "member_id")
	m.SiteID = field.NewInt32(table, "site_id")
	m.SourceMember = field.NewInt32(table, "source_member")
	m.FenxiaoID = field.NewInt32(table, "fenxiao_id")
	m.IsFenxiao = field.NewBool(table, "is_fenxiao")
	m.Username = field.NewString(table, "username")
	m.Nickname = field.NewString(table, "nickname")
	m.Mobile = field.NewString(table, "mobile")
	m.Email = field.NewString(table, "email")
	m.Password = field.NewString(table, "password")
	m.Status = field.NewInt32(table, "status")
	m.Headimg = field.NewString(table, "headimg")
	m.MemberLevel = field.NewInt32(table, "member_level")
	m.MemberLevelName = field.NewString(table, "member_level_name")
	m.MemberLabel = field.NewInt32(table, "member_label")
	m.MemberLabelName = field.NewString(table, "member_label_name")
	m.Qq = field.NewString(table, "qq")
	m.QqOpenid = field.NewString(table, "qq_openid")
	m.WxOpenid = field.NewString(table, "wx_openid")
	m.WeappOpenid = field.NewString(table, "weapp_openid")
	m.WxUnionid = field.NewString(table, "wx_unionid")
	m.AliOpenid = field.NewString(table, "ali_openid")
	m.BaiduOpenid = field.NewString(table, "baidu_openid")
	m.ToutiaoOpenid = field.NewString(table, "toutiao_openid")
	m.DouyinOpenid = field.NewString(table, "douyin_openid")
	m.LoginIP = field.NewString(table, "login_ip")
	m.LoginType = field.NewString(table, "login_type")
	m.LoginTime = field.NewInt32(table, "login_time")
	m.LastLoginIP = field.NewString(table, "last_login_ip")
	m.LastLoginType = field.NewString(table, "last_login_type")
	m.LastLoginTime = field.NewInt32(table, "last_login_time")
	m.LoginNum = field.NewInt32(table, "login_num")
	m.Realname = field.NewString(table, "realname")
	m.Sex = field.NewInt32(table, "sex")
	m.Location = field.NewString(table, "location")
	m.Birthday = field.NewInt32(table, "birthday")
	m.RegTime = field.NewInt32(table, "reg_time")
	m.Point = field.NewFloat64(table, "point")
	m.Balance = field.NewFloat64(table, "balance")
	m.Growth = field.NewFloat64(table, "growth")
	m.BalanceMoney = field.NewFloat64(table, "balance_money")
	m.Account5 = field.NewFloat64(table, "account5")
	m.IsAuth = field.NewInt32(table, "is_auth")
	m.SignTime = field.NewInt32(table, "sign_time")
	m.SignDaysSeries = field.NewInt32(table, "sign_days_series")
	m.PayPassword = field.NewString(table, "pay_password")
	m.OrderMoney = field.NewFloat64(table, "order_money")
	m.OrderCompleteMoney = field.NewFloat64(table, "order_complete_money")
	m.OrderNum = field.NewInt32(table, "order_num")
	m.OrderCompleteNum = field.NewInt32(table, "order_complete_num")
	m.BalanceWithdrawApply = field.NewFloat64(table, "balance_withdraw_apply")
	m.BalanceWithdraw = field.NewFloat64(table, "balance_withdraw")
	m.IsDelete = field.NewBool(table, "is_delete")
	m.MemberLevelType = field.NewInt32(table, "member_level_type")
	m.LevelExpireTime = field.NewInt32(table, "level_expire_time")
	m.IsEditUsername = field.NewInt32(table, "is_edit_username")
	m.LoginTypeName = field.NewString(table, "login_type_name")
	m.CanReceiveRegistergift = field.NewInt32(table, "can_receive_registergift")
	m.ProvinceID = field.NewInt32(table, "province_id")
	m.CityID = field.NewInt32(table, "city_id")
	m.DistrictID = field.NewInt32(table, "district_id")
	m.FullAddress = field.NewString(table, "full_address")
	m.Address = field.NewString(table, "address")
	m.ErpUID = field.NewString(table, "erp_uid")
	m.ErpCode = field.NewString(table, "erp_code")
	m.ErpName = field.NewString(table, "erp_name")
	m.Gid = field.NewInt32(table, "gid")
	m.QualificationsImages = field.NewString(table, "qualifications_images")
	m.Inputmanid = field.NewInt32(table, "inputmanid")
	m.Deptid = field.NewInt32(table, "deptid")
	m.Employeename = field.NewString(table, "employeename")
	m.ErpAgentid = field.NewInt32(table, "erp_agentid")
	m.Customopcode = field.NewString(table, "customopcode")
	m.SyncTime = field.NewInt32(table, "sync_time")
	m.MemberRemark = field.NewString(table, "member_remark")
	m.Employeeid = field.NewString(table, "employeeid")
	m.InviterID = field.NewInt32(table, "inviter_id")
	m.InviteCode = field.NewString(table, "invite_code")
	m.QrInviteCode = field.NewString(table, "qr_invite_code")
	m.IsNatureofbusiness = field.NewInt32(table, "is_natureofbusiness")
	m.SalesmanID = field.NewInt32(table, "salesman_id")
	m.StaffType = field.NewInt32(table, "staff_type")
	m.FristLoginTime = field.NewInt32(table, "frist_login_time")
	m.RegMobile = field.NewString(table, "reg_mobile")
	m.RegPassword = field.NewString(table, "reg_password")
	m.AppointmenSalesmanID = field.NewInt32(table, "appointmen_salesman_id")
	m.CreditBalance = field.NewFloat64(table, "credit_balance")
	m.CreditLimit = field.NewFloat64(table, "credit_limit")
	m.SessionKey = field.NewString(table, "session_key")

	m.fillFieldMap()

	return m
}

func (m *member) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *member) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 91)
	m.fieldMap["member_id"] = m.MemberID
	m.fieldMap["site_id"] = m.SiteID
	m.fieldMap["source_member"] = m.SourceMember
	m.fieldMap["fenxiao_id"] = m.FenxiaoID
	m.fieldMap["is_fenxiao"] = m.IsFenxiao
	m.fieldMap["username"] = m.Username
	m.fieldMap["nickname"] = m.Nickname
	m.fieldMap["mobile"] = m.Mobile
	m.fieldMap["email"] = m.Email
	m.fieldMap["password"] = m.Password
	m.fieldMap["status"] = m.Status
	m.fieldMap["headimg"] = m.Headimg
	m.fieldMap["member_level"] = m.MemberLevel
	m.fieldMap["member_level_name"] = m.MemberLevelName
	m.fieldMap["member_label"] = m.MemberLabel
	m.fieldMap["member_label_name"] = m.MemberLabelName
	m.fieldMap["qq"] = m.Qq
	m.fieldMap["qq_openid"] = m.QqOpenid
	m.fieldMap["wx_openid"] = m.WxOpenid
	m.fieldMap["weapp_openid"] = m.WeappOpenid
	m.fieldMap["wx_unionid"] = m.WxUnionid
	m.fieldMap["ali_openid"] = m.AliOpenid
	m.fieldMap["baidu_openid"] = m.BaiduOpenid
	m.fieldMap["toutiao_openid"] = m.ToutiaoOpenid
	m.fieldMap["douyin_openid"] = m.DouyinOpenid
	m.fieldMap["login_ip"] = m.LoginIP
	m.fieldMap["login_type"] = m.LoginType
	m.fieldMap["login_time"] = m.LoginTime
	m.fieldMap["last_login_ip"] = m.LastLoginIP
	m.fieldMap["last_login_type"] = m.LastLoginType
	m.fieldMap["last_login_time"] = m.LastLoginTime
	m.fieldMap["login_num"] = m.LoginNum
	m.fieldMap["realname"] = m.Realname
	m.fieldMap["sex"] = m.Sex
	m.fieldMap["location"] = m.Location
	m.fieldMap["birthday"] = m.Birthday
	m.fieldMap["reg_time"] = m.RegTime
	m.fieldMap["point"] = m.Point
	m.fieldMap["balance"] = m.Balance
	m.fieldMap["growth"] = m.Growth
	m.fieldMap["balance_money"] = m.BalanceMoney
	m.fieldMap["account5"] = m.Account5
	m.fieldMap["is_auth"] = m.IsAuth
	m.fieldMap["sign_time"] = m.SignTime
	m.fieldMap["sign_days_series"] = m.SignDaysSeries
	m.fieldMap["pay_password"] = m.PayPassword
	m.fieldMap["order_money"] = m.OrderMoney
	m.fieldMap["order_complete_money"] = m.OrderCompleteMoney
	m.fieldMap["order_num"] = m.OrderNum
	m.fieldMap["order_complete_num"] = m.OrderCompleteNum
	m.fieldMap["balance_withdraw_apply"] = m.BalanceWithdrawApply
	m.fieldMap["balance_withdraw"] = m.BalanceWithdraw
	m.fieldMap["is_delete"] = m.IsDelete
	m.fieldMap["member_level_type"] = m.MemberLevelType
	m.fieldMap["level_expire_time"] = m.LevelExpireTime
	m.fieldMap["is_edit_username"] = m.IsEditUsername
	m.fieldMap["login_type_name"] = m.LoginTypeName
	m.fieldMap["can_receive_registergift"] = m.CanReceiveRegistergift
	m.fieldMap["province_id"] = m.ProvinceID
	m.fieldMap["city_id"] = m.CityID
	m.fieldMap["district_id"] = m.DistrictID
	m.fieldMap["full_address"] = m.FullAddress
	m.fieldMap["address"] = m.Address
	m.fieldMap["erp_uid"] = m.ErpUID
	m.fieldMap["erp_code"] = m.ErpCode
	m.fieldMap["erp_name"] = m.ErpName
	m.fieldMap["gid"] = m.Gid
	m.fieldMap["qualifications_images"] = m.QualificationsImages
	m.fieldMap["inputmanid"] = m.Inputmanid
	m.fieldMap["deptid"] = m.Deptid
	m.fieldMap["employeename"] = m.Employeename
	m.fieldMap["erp_agentid"] = m.ErpAgentid
	m.fieldMap["customopcode"] = m.Customopcode
	m.fieldMap["sync_time"] = m.SyncTime
	m.fieldMap["member_remark"] = m.MemberRemark
	m.fieldMap["employeeid"] = m.Employeeid
	m.fieldMap["inviter_id"] = m.InviterID
	m.fieldMap["invite_code"] = m.InviteCode
	m.fieldMap["qr_invite_code"] = m.QrInviteCode
	m.fieldMap["is_natureofbusiness"] = m.IsNatureofbusiness
	m.fieldMap["salesman_id"] = m.SalesmanID
	m.fieldMap["staff_type"] = m.StaffType
	m.fieldMap["frist_login_time"] = m.FristLoginTime
	m.fieldMap["reg_mobile"] = m.RegMobile
	m.fieldMap["reg_password"] = m.RegPassword
	m.fieldMap["appointmen_salesman_id"] = m.AppointmenSalesmanID
	m.fieldMap["credit_balance"] = m.CreditBalance
	m.fieldMap["credit_limit"] = m.CreditLimit
	m.fieldMap["session_key"] = m.SessionKey

}

func (m member) clone(db *gorm.DB) member {
	m.memberDo.ReplaceConnPool(db.Statement.ConnPool)
	m.MemberQualification.db = db.Session(&gorm.Session{Initialized: true})
	m.MemberQualification.db.Statement.ConnPool = db.Statement.ConnPool
	m.MemberAddress.db = db.Session(&gorm.Session{Initialized: true})
	m.MemberAddress.db.Statement.ConnPool = db.Statement.ConnPool
	return m
}

func (m member) replaceDB(db *gorm.DB) member {
	m.memberDo.ReplaceDB(db)
	m.MemberQualification.db = db.Session(&gorm.Session{})
	m.MemberAddress.db = db.Session(&gorm.Session{})
	return m
}

type memberHasManyMemberQualification struct {
	db *gorm.DB

	field.RelationField
}

func (a memberHasManyMemberQualification) Where(conds ...field.Expr) *memberHasManyMemberQualification {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a memberHasManyMemberQualification) WithContext(ctx context.Context) *memberHasManyMemberQualification {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a memberHasManyMemberQualification) Session(session *gorm.Session) *memberHasManyMemberQualification {
	a.db = a.db.Session(session)
	return &a
}

func (a memberHasManyMemberQualification) Model(m *shop_model.Member) *memberHasManyMemberQualificationTx {
	return &memberHasManyMemberQualificationTx{a.db.Model(m).Association(a.Name())}
}

func (a memberHasManyMemberQualification) Unscoped() *memberHasManyMemberQualification {
	a.db = a.db.Unscoped()
	return &a
}

type memberHasManyMemberQualificationTx struct{ tx *gorm.Association }

func (a memberHasManyMemberQualificationTx) Find() (result []*shop_model.MemberQualification, err error) {
	return result, a.tx.Find(&result)
}

func (a memberHasManyMemberQualificationTx) Append(values ...*shop_model.MemberQualification) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a memberHasManyMemberQualificationTx) Replace(values ...*shop_model.MemberQualification) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a memberHasManyMemberQualificationTx) Delete(values ...*shop_model.MemberQualification) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a memberHasManyMemberQualificationTx) Clear() error {
	return a.tx.Clear()
}

func (a memberHasManyMemberQualificationTx) Count() int64 {
	return a.tx.Count()
}

func (a memberHasManyMemberQualificationTx) Unscoped() *memberHasManyMemberQualificationTx {
	a.tx = a.tx.Unscoped()
	return &a
}

type memberHasManyMemberAddress struct {
	db *gorm.DB

	field.RelationField
}

func (a memberHasManyMemberAddress) Where(conds ...field.Expr) *memberHasManyMemberAddress {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a memberHasManyMemberAddress) WithContext(ctx context.Context) *memberHasManyMemberAddress {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a memberHasManyMemberAddress) Session(session *gorm.Session) *memberHasManyMemberAddress {
	a.db = a.db.Session(session)
	return &a
}

func (a memberHasManyMemberAddress) Model(m *shop_model.Member) *memberHasManyMemberAddressTx {
	return &memberHasManyMemberAddressTx{a.db.Model(m).Association(a.Name())}
}

func (a memberHasManyMemberAddress) Unscoped() *memberHasManyMemberAddress {
	a.db = a.db.Unscoped()
	return &a
}

type memberHasManyMemberAddressTx struct{ tx *gorm.Association }

func (a memberHasManyMemberAddressTx) Find() (result []*shop_model.MemberAddress, err error) {
	return result, a.tx.Find(&result)
}

func (a memberHasManyMemberAddressTx) Append(values ...*shop_model.MemberAddress) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a memberHasManyMemberAddressTx) Replace(values ...*shop_model.MemberAddress) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a memberHasManyMemberAddressTx) Delete(values ...*shop_model.MemberAddress) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a memberHasManyMemberAddressTx) Clear() error {
	return a.tx.Clear()
}

func (a memberHasManyMemberAddressTx) Count() int64 {
	return a.tx.Count()
}

func (a memberHasManyMemberAddressTx) Unscoped() *memberHasManyMemberAddressTx {
	a.tx = a.tx.Unscoped()
	return &a
}

type memberDo struct{ gen.DO }

type IMemberDo interface {
	gen.SubQuery
	Debug() IMemberDo
	WithContext(ctx context.Context) IMemberDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMemberDo
	WriteDB() IMemberDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMemberDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMemberDo
	Not(conds ...gen.Condition) IMemberDo
	Or(conds ...gen.Condition) IMemberDo
	Select(conds ...field.Expr) IMemberDo
	Where(conds ...gen.Condition) IMemberDo
	Order(conds ...field.Expr) IMemberDo
	Distinct(cols ...field.Expr) IMemberDo
	Omit(cols ...field.Expr) IMemberDo
	Join(table schema.Tabler, on ...field.Expr) IMemberDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMemberDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMemberDo
	Group(cols ...field.Expr) IMemberDo
	Having(conds ...gen.Condition) IMemberDo
	Limit(limit int) IMemberDo
	Offset(offset int) IMemberDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMemberDo
	Unscoped() IMemberDo
	Create(values ...*shop_model.Member) error
	CreateInBatches(values []*shop_model.Member, batchSize int) error
	Save(values ...*shop_model.Member) error
	First() (*shop_model.Member, error)
	Take() (*shop_model.Member, error)
	Last() (*shop_model.Member, error)
	Find() ([]*shop_model.Member, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*shop_model.Member, err error)
	FindInBatches(result *[]*shop_model.Member, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*shop_model.Member) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMemberDo
	Assign(attrs ...field.AssignExpr) IMemberDo
	Joins(fields ...field.RelationField) IMemberDo
	Preload(fields ...field.RelationField) IMemberDo
	FirstOrInit() (*shop_model.Member, error)
	FirstOrCreate() (*shop_model.Member, error)
	FindByPage(offset int, limit int) (result []*shop_model.Member, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMemberDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m memberDo) Debug() IMemberDo {
	return m.withDO(m.DO.Debug())
}

func (m memberDo) WithContext(ctx context.Context) IMemberDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m memberDo) ReadDB() IMemberDo {
	return m.Clauses(dbresolver.Read)
}

func (m memberDo) WriteDB() IMemberDo {
	return m.Clauses(dbresolver.Write)
}

func (m memberDo) Session(config *gorm.Session) IMemberDo {
	return m.withDO(m.DO.Session(config))
}

func (m memberDo) Clauses(conds ...clause.Expression) IMemberDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m memberDo) Returning(value interface{}, columns ...string) IMemberDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m memberDo) Not(conds ...gen.Condition) IMemberDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m memberDo) Or(conds ...gen.Condition) IMemberDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m memberDo) Select(conds ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m memberDo) Where(conds ...gen.Condition) IMemberDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m memberDo) Order(conds ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m memberDo) Distinct(cols ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m memberDo) Omit(cols ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m memberDo) Join(table schema.Tabler, on ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m memberDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMemberDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m memberDo) RightJoin(table schema.Tabler, on ...field.Expr) IMemberDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m memberDo) Group(cols ...field.Expr) IMemberDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m memberDo) Having(conds ...gen.Condition) IMemberDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m memberDo) Limit(limit int) IMemberDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m memberDo) Offset(offset int) IMemberDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m memberDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMemberDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m memberDo) Unscoped() IMemberDo {
	return m.withDO(m.DO.Unscoped())
}

func (m memberDo) Create(values ...*shop_model.Member) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m memberDo) CreateInBatches(values []*shop_model.Member, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m memberDo) Save(values ...*shop_model.Member) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m memberDo) First() (*shop_model.Member, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.Member), nil
	}
}

func (m memberDo) Take() (*shop_model.Member, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.Member), nil
	}
}

func (m memberDo) Last() (*shop_model.Member, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.Member), nil
	}
}

func (m memberDo) Find() ([]*shop_model.Member, error) {
	result, err := m.DO.Find()
	return result.([]*shop_model.Member), err
}

func (m memberDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*shop_model.Member, err error) {
	buf := make([]*shop_model.Member, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m memberDo) FindInBatches(result *[]*shop_model.Member, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m memberDo) Attrs(attrs ...field.AssignExpr) IMemberDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m memberDo) Assign(attrs ...field.AssignExpr) IMemberDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m memberDo) Joins(fields ...field.RelationField) IMemberDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m memberDo) Preload(fields ...field.RelationField) IMemberDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m memberDo) FirstOrInit() (*shop_model.Member, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.Member), nil
	}
}

func (m memberDo) FirstOrCreate() (*shop_model.Member, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.Member), nil
	}
}

func (m memberDo) FindByPage(offset int, limit int) (result []*shop_model.Member, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m memberDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m memberDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m memberDo) Delete(models ...*shop_model.Member) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *memberDo) withDO(do gen.Dao) *memberDo {
	m.DO = *do.(*gen.DO)
	return m
}
