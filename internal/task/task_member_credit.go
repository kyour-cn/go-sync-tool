package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_model"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"app/ui/apptheme"
	"errors"
	"gioui.org/layout"
	"gorm.io/gorm"
	"log/slog"
)

func NewMemberCredit() *MemberCredit {
	return &MemberCredit{}
}

// MemberCredit 同步ERP商品到商城
type MemberCredit struct {
	IsRunning bool
}

func (mc MemberCredit) GetName() string {
	return "MemberCredit"
}

func (MemberCredit) ClearCache() error {
	return store.MemberCreditStore.Clear()
}

func (mc MemberCredit) Run(t *Task) error {

	defer func() {
		// 缓存数据到文件
		err := store.MemberCreditStore.Save()
		if err != nil {
			slog.Error("SaveMemberCredit err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.MemberCredit

	erpDb, ok := global.DbPool.Get("erp")
	if !ok {
		return errors.New("获取ERP数据库连接失败")
	}

	// 执行SQL查询
	r := erpDb.Raw(t.Config.Sql).Scan(&erpData)
	if r.Error != nil {
		return r.Error
	}

	// 创建新的Map
	newMap := safemap.New[*erp_entity.MemberCredit]()
	for _, item := range erpData {
		newMap.Set(item.ErpUID, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.MemberCredit](store.MemberCreditStore.Store, newMap)
	newMap = nil

	slog.Info("客户资信同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.MemberCredit) error {
		err := mc.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberCreditStore.Store.Set(v.ErpUID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.MemberCredit) error {
		err := mc.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberCreditStore.Store.Set(v.ErpUID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.MemberCredit) error {
		err := mc.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberCreditStore.Store.Delete(v.ErpUID)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (mc MemberCredit) addOrUpdate(item *erp_entity.MemberCredit) error {

	member, err := shop_query.Member.
		Where(shop_query.Member.ErpUID.Eq(item.ErpUID)).
		Select(
			shop_query.Member.MemberID,
			shop_query.Member.ErpUID,
			shop_query.Member.CreditBalance,
			shop_query.Member.CreditLimit,
		).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("查询客户资信信息 err: " + err.Error())
		return err
	}
	if member == nil {
		return nil
	}

	// 比较是否有更新
	if member.CreditLimit == item.Limit && member.CreditBalance == item.Money {
		return nil
	}

	// 更新Member表
	_, e := shop_query.Member.
		Where(shop_query.Member.MemberID.Eq(member.MemberID)).
		Select(shop_query.Member.CreditLimit, shop_query.Member.CreditBalance).
		Updates(shop_model.Member{
			CreditLimit:   item.Limit,
			CreditBalance: item.Money,
		})
	if e != nil {
		slog.Error("资信更新 err" + e.Error())
		return e
	}
	return nil
}

func (mc MemberCredit) delete(member *erp_entity.MemberCredit) error {
	// 更新为0
	member.Limit = 0
	member.Money = 0
	return mc.addOrUpdate(member)
}

// ConfigLayout 任务配置UI布局
func (mc MemberCredit) ConfigLayout(_ layout.Context, _ *apptheme.Theme, _ *Task) layout.Dimensions {
	return layout.Dimensions{}
}
