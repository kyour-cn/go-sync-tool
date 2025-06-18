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
	"time"
)

func NewSalesman() *Salesman {
	return &Salesman{}
}

// Salesman 同步ERP商品到商城
type Salesman struct{}

func (m Salesman) GetName() string {
	return "Salesman"
}

func (Salesman) ClearCache() error {
	return store.SalesmanStore.Clear()
}

func (m Salesman) Run(t *Task) error {
	defer func() {
		// 缓存数据到文件
		err := store.SalesmanStore.Save()
		if err != nil {
			slog.Error("SaveSalesman err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.Salesman

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
	newMap := safemap.New[*erp_entity.Salesman]()
	for _, item := range erpData {
		newMap.Set(item.SaleID, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.Salesman](store.SalesmanStore.Store, newMap)

	slog.Info("销售员比对", "old", store.SalesmanStore.Store.Len(), "new", newMap.Len(), "add", add.Len(), "update", update.Len(), "del", del.Len())
	newMap = nil

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.Salesman) error {
		err := m.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.SalesmanStore.Store.Set(v.SaleID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.Salesman) error {
		err := m.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.SalesmanStore.Store.Set(v.SaleID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.Salesman) error {
		err := m.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.SalesmanStore.Store.Delete(v.SaleID)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (m Salesman) addOrUpdate(item *erp_entity.Salesman) error {

	salesmanInfo, err := shop_query.StaffSalesman.Where(shop_query.StaffSalesman.ErpSaleerID.Eq(item.SaleID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("业务员同步 addOrUpdate", "err", err)
		return err
	}

	if salesmanInfo != nil {
		if er := m.update(item, salesmanInfo); er != nil {
			slog.Error("salesmanSync updateMember", "err", err)
			return er
		}
	} else {
		if er := m.add(item); er != nil {
			slog.Error("memberSync addMember", "err", err)
			return er
		}
	}
	return nil
}

func (m Salesman) add(v *erp_entity.Salesman) error {
	memberData, err := shop_query.Member.Where(shop_query.Member.Username.Eq(v.SaleID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if memberData == nil {
		memberData = &shop_model.Member{
			SiteID:    1,
			Username:  "ywy" + v.SaleID,
			Nickname:  v.Realname.String(),
			ErpName:   v.Realname.String(),
			Realname:  v.Realname.String(),
			Mobile:    v.Mobile.String(),
			StaffType: int32(1),
			Password:  "e0145437e0c0a26c644efab6f97f2985",
			RegTime:   int32(time.Now().Unix()),
			SyncTime:  int32(time.Now().Unix()),
			Status:    1,
		}
		if er := shop_query.Member.Create(memberData); er != nil {
			return er
		}
	}

	sm := shop_model.StaffSalesman{
		ErpSaleerID:     v.SaleID,
		MemberID:        memberData.MemberID,
		FirstID:         0,
		SecondID:        0,
		RegionID:        0,
		SalesmanName:    v.Realname.String(),
		SalesmanMobile:  v.Mobile.String(),
		SalesmanAccount: "ywy" + v.SaleID,
		Level:           3,
		Status:          0,
	}
	er := shop_query.StaffSalesman.Create(&sm)
	if er != nil {
		slog.Error("addSalesman StaffSalesman Create err:", "err", er)
		return er
	}
	// 更新用户表
	_, ers := shop_query.Member.Where(shop_query.Member.MemberID.Eq(sm.MemberID)).Update(shop_query.Member.SalesmanID, sm.SalesmanID)
	if ers != nil {
		slog.Error("addSalesman Member Update err:", "err", ers)
		return ers
	}
	return nil
}

func (m Salesman) update(v *erp_entity.Salesman, salesman *shop_model.StaffSalesman) error {

	nowTime := int32(time.Now().Unix())

	member := shop_model.Member{
		Nickname: v.Realname.String(),
		ErpName:  v.Realname.String(),
		Realname: v.Realname.String(),
		Mobile:   v.Mobile.String(),
		SyncTime: nowTime,
	}
	_, err := shop_query.Member.
		Where(shop_query.Member.MemberID.Eq(salesman.MemberID)).
		Where(shop_query.Member.StaffType.Eq(1)).
		Select(
			shop_query.Member.Nickname,
			shop_query.Member.ErpName,
			shop_query.Member.Realname,
			shop_query.Member.Mobile,
			shop_query.Member.SyncTime,
		).
		Updates(&member)
	if err != nil {
		return err
	}

	smData := shop_model.StaffSalesman{
		//SiteID:       1,
		SalesmanName:   v.Realname.String(),
		SalesmanMobile: v.Mobile.String(),
		//Introduction: v.Realname.String(),
	}
	_, er := shop_query.StaffSalesman.
		Where(shop_query.StaffSalesman.MemberID.Eq(salesman.MemberID)).
		Select(
			//shop_query.StaffSalesman.SiteID,
			shop_query.StaffSalesman.SalesmanName,
			shop_query.StaffSalesman.SalesmanMobile,
			//shop_query.StaffSalesman.Introduction,
		).
		Updates(&smData)
	return er
}

func (m Salesman) delete(_ *erp_entity.Salesman) error {
	return nil
}

// ConfigLayout 任务配置UI布局
func (m Salesman) ConfigLayout(_ layout.Context, _ *apptheme.Theme, _ *Task) layout.Dimensions {
	return layout.Dimensions{}
}
