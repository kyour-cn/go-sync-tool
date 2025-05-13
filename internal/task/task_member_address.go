package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_model"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

func NewMemberAddress() *MemberAddress {
	return &MemberAddress{}
}

// MemberAddress 同步ERP商品到商城
type MemberAddress struct{}

func (ma MemberAddress) GetName() string {
	return "MemberAddress"
}

func (MemberAddress) ClearCache() error {
	return store.MemberAddressStore.Clear()
}

func (ma MemberAddress) Run(t *Task) error {
	defer func() {
		// 缓存数据到文件
		err := store.MemberAddressStore.Save()
		if err != nil {
			slog.Error("SaveMemberAddress err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.MemberAddress

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
	newMap := safemap.New[*erp_entity.MemberAddress]()
	for _, item := range erpData {
		newMap.Set(item.ID, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.MemberAddress](store.MemberAddressStore.Store, newMap)
	newMap = nil

	slog.Info("商品价格同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.MemberAddress) error {
		err := ma.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberAddressStore.Store.Set(v.ID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.MemberAddress) error {
		err := ma.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberAddressStore.Store.Set(v.ID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.MemberAddress) error {
		err := ma.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberAddressStore.Store.Delete(v.ID)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (ma MemberAddress) addOrUpdate(item *erp_entity.MemberAddress) error {

	//查询是否存在商城表中
	memberInfo, err := shop_query.Member.
		Select(shop_query.Member.MemberID).
		Where(shop_query.Member.ErpUID.Eq(item.ErpUID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("查询商城会员信息失败", "err", err.Error())
		return err
	}
	if memberInfo == nil { //关联会员未同步直接跳过
		return nil
	}

	//查询用户地址是否存在 (同步ERP只许有一个地址)
	memberAddress, err := shop_query.MemberAddress.
		Where(
			shop_query.MemberAddress.MemberID.Eq(memberInfo.MemberID),
		).
		// 优先查询默认地址
		Order(shop_query.MemberAddress.IsDefault.Desc()).
		First()

	// 如果会员没有地区信息，则重新获取
	// 使用ERP地址信息的省市区信息获取商城的省市区信息
	areaInfo := getAreaFormCache(item.Province.String(), item.City.String(), item.District.String())
	memberInfo.ProvinceID = areaInfo.ProvinceID
	memberInfo.CityID = areaInfo.CityID
	memberInfo.DistrictID = areaInfo.DistrictID
	memberInfo.FullAddress = areaInfo.ProvinceName + "-" + areaInfo.CityName + "-" + areaInfo.DistrictName

	// 更新会员地区信息
	_, _ = shop_query.Member.
		Where(
			shop_query.Member.MemberID.Eq(memberInfo.MemberID),
		).
		Updates(shop_model.Member{
			ProvinceID:  memberInfo.ProvinceID,
			CityID:      memberInfo.CityID,
			DistrictID:  memberInfo.DistrictID,
			FullAddress: memberInfo.FullAddress,
		})

	if memberAddress != nil {
		if er := ma.update(item, memberInfo, memberAddress); er != nil {
			slog.Error("会员地址同步 更新数据失败", "erp_uid", item.ErpUID, "err", er.Error())
			return er
		}
	} else {
		if er := ma.add(item, memberInfo); er != nil {
			slog.Error("会员地址同步 添加数据失败", "erp_uid", item.ErpUID, "err", er.Error())
			return er
		}
	}

	return nil
}

func (ma MemberAddress) add(v *erp_entity.MemberAddress, m *shop_model.Member) error {
	memberAddressData := shop_model.MemberAddress{
		MemberID:    m.MemberID,
		SiteID:      1,
		Name:        v.RealName.String(),
		Mobile:      v.Mobile.String(),
		Telephone:   v.Mobile.String(),
		ProvinceID:  m.ProvinceID,
		CityID:      m.CityID,
		DistrictID:  m.DistrictID,
		Address:     v.Address.String(),
		FullAddress: m.FullAddress,
		IsDefault:   1,
	}
	return shop_query.MemberAddress.Create(&memberAddressData)
}

func (ma MemberAddress) update(v *erp_entity.MemberAddress, m *shop_model.Member, md *shop_model.MemberAddress) error {
	memberAddressData := shop_model.MemberAddress{
		SiteID:      1,
		Name:        v.RealName.String(),
		Mobile:      v.Mobile.String(),
		Telephone:   v.Mobile.String(),
		ProvinceID:  m.ProvinceID,
		CityID:      m.CityID,
		DistrictID:  m.DistrictID,
		Address:     v.Address.String(),
		FullAddress: m.FullAddress,
	}

	// 如果ERP手机号信息为空，则不更新
	if md.Mobile != "" && memberAddressData.Mobile == "" {
		memberAddressData.Mobile = md.Mobile
		memberAddressData.Telephone = md.Mobile
	}

	_, err := shop_query.MemberAddress.
		Where(shop_query.MemberAddress.ID.Eq(md.ID)).
		Select(
			shop_query.MemberAddress.SiteID,
			shop_query.MemberAddress.Name,
			shop_query.MemberAddress.Mobile,
			shop_query.MemberAddress.Telephone,
			shop_query.MemberAddress.ProvinceID,
			shop_query.MemberAddress.CityID,
			shop_query.MemberAddress.DistrictID,
			shop_query.MemberAddress.Address,
			shop_query.MemberAddress.FullAddress,
		).
		Updates(&memberAddressData)
	return err
}

func (ma MemberAddress) delete(_ *erp_entity.MemberAddress) error {

	return nil
}
