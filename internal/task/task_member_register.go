package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_query"
	"app/ui/apptheme"
	"errors"
	"fmt"
	"gioui.org/layout"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func NewMemberRegister() *MemberRegister {
	return &MemberRegister{
		newMemberTable:              "jxkj_new_member",
		newMemberQualificationTable: "jxkj_new_member_qualifications",
	}
}

// MemberRegister 商城的新客户同步到ERP
type MemberRegister struct {
	newMemberTable              string
	newMemberQualificationTable string
}

func (r MemberRegister) GetName() string {
	return "MemberRegister"
}

func (MemberRegister) ClearCache() error {
	return nil
}

func (r MemberRegister) Run(t *Task) error {

	memberList, err := shop_query.Member.
		Preload(shop_query.Member.MemberQualification).
		Preload(shop_query.Member.MemberAddress).
		Where(
			shop_query.Member.Status.Eq(1),
			shop_query.Member.SyncTime.Eq(0),
			shop_query.Member.IsDelete.Is(false),
			shop_query.Member.ErpUID.Eq(""),
		).
		Find()
	if err != nil {
		slog.Error("查询平台新会员异常", "err", err)
		return err
	}
	if len(memberList) == 0 {
		return nil
	}

	t.DataCount = len(memberList)

	//当前时间
	nowTime := time.Now()

	for _, member := range memberList {

		//将full-address 用-分割
		addressArr := strings.Split(member.FullAddress, "-")
		if len(addressArr) < 3 {
			addressArr = []string{"", "", ""}
		}

		//查询商城中的会员销售员
		//salesmanInfo, _ := shop_query.Fenxiao.
		//	Where(shop_query.Fenxiao.FenxiaoID.Eq(member.FenxiaoID)).First()
		//if salesmanInfo == nil {
		//	salesmanInfo = &shop_entity.Fenxiao{}
		//}

		newMember := erp_entity.NewMember{
			MemberID:   int64(member.MemberID),
			Nickname:   erp_entity.UTF8String(member.Nickname),
			Realname:   erp_entity.UTF8String(member.Realname),
			Username:   erp_entity.UTF8String(member.Username),
			Mobile:     erp_entity.UTF8String(member.Mobile),
			MemberType: erp_entity.UTF8String(member.MemberLabelName),
			Province:   erp_entity.UTF8String(addressArr[0]),
			City:       erp_entity.UTF8String(addressArr[1]),
			District:   erp_entity.UTF8String(addressArr[2]),
			Address:    "",
			SyncTime:   erp_entity.UTF8String(nowTime.Format("2006-01-02 15:04:05")),
		}
		if len(member.MemberAddress) > 0 {
			newMember.Address = erp_entity.UTF8String(member.MemberAddress[0].Address)
		}

		// 查询客户经营范围
		businessScopeList, err := shop_query.MemberBusinessScope.
			Where(shop_query.MemberBusinessScope.MemberID.Eq(member.MemberID)).
			Find()
		if err == nil && len(businessScopeList) > 0 {
			businessScopeNames := make([]string, 0)
			for _, v := range businessScopeList {
				businessScopeNames = append(businessScopeNames, v.BusinessScope)
			}
			newMember.BusinessScope = erp_entity.UTF8String(strings.Join(businessScopeNames, ","))
		}

		// 获取ERP数据库连接
		erpDb, ok := global.DbPool.Get("erp")
		if !ok {
			return errors.New("获取ERP数据库连接失败")
		}

		if len(member.MemberQualification) > 0 {

			for _, v := range member.MemberQualification {
				if v.Status != 1 { //不同步未完成审核资质
					continue
				}
				mq := erp_entity.NewMemberQualification{
					MemberID: int(v.MemberID),
					Name:     erp_entity.UTF8String(v.Name),
					Identify: erp_entity.UTF8String(v.Identify),
					LongTerm: int(v.LongTerm),
					Image:    erp_entity.UTF8String(v.Image),
					CardNo:   erp_entity.UTF8String(v.CardNo),
				}

				// 自定义表单
				if v.Custom != "" {
					mq.CustomForm = erp_entity.UTF8String(v.Custom)
				}

				if v.ExpirationStartDate != nil {
					mq.ExpirationStartDate = v.ExpirationStartDate.Format("2006-01-02")
				}
				if v.ExpirationEndDate != nil {
					mq.ExpirationEndDate = v.ExpirationEndDate.Format("2006-01-02")
				}

				// 查询资质经营范围
				if v.BusinessScope != "" {
					var busIds []int32
					for _, v := range strings.Split(v.BusinessScope, ",") {
						id, _ := strconv.Atoi(v)
						busIds = append(busIds, int32(id))
					}
					businessScopeList, err := shop_query.MemberBusinessScopeRow.
						Where(shop_query.MemberBusinessScopeRow.ID.In(busIds...)).
						Find()
					if err == nil && len(businessScopeList) > 0 {
						businessScopeNames := make([]string, 0)
						for _, v := range businessScopeList {
							businessScopeNames = append(businessScopeNames, v.Name)
						}
						mq.BusinessScope = erp_entity.UTF8String(strings.Join(businessScopeNames, ","))
					}
				}

				// 写入ERP
				result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
					return tx.Table(r.newMemberQualificationTable).Create(&mq)
				}))
				if result.Error != nil {
					slog.Error(fmt.Sprintf("ERP CreateOrder err:%s,args:%+v", result.Error, mq))
					return nil
				}
			}
		}

		// 写入ERP
		result := erpDb.Exec(erpDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Table(r.newMemberTable).Create(&newMember)
		}))
		if result.Error != nil {
			slog.Error(fmt.Sprintf("ERP CreateMewMember err:%s,args:%+v", result.Error, newMember))
			return nil
		}

		member.SyncTime = int32(nowTime.Unix())
		//更新商城中的sync_time
		_, err = shop_query.Member.
			Where(shop_query.Member.MemberID.Eq(member.MemberID)).
			Update(shop_query.Member.SyncTime, int32(nowTime.Unix()))
		if err != nil {
			slog.Error("newMemberSync Member Save err:%s", err)
			continue
		}

		t.DoneCount++
	}

	return nil
}

// ConfigLayout 任务配置UI布局
func (r MemberRegister) ConfigLayout(_ layout.Context, _ *apptheme.Theme, _ *Task) layout.Dimensions {
	return layout.Dimensions{}
}
