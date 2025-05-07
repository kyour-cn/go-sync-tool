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

// MemberAddress 同步ERP商品到商城
type MemberAddress struct {
}

func (g MemberAddress) GetName() string {
    return "MemberAddress"
}

func (g MemberAddress) Run(t *Task) error {
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
    r := erpDb.Db.Raw(t.Config.Sql).Scan(&erpData)
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

    // 添加
    for _, v := range add.Values() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        addOrUpdateMemberAddress(v)
        store.MemberAddressStore.Store.Set(v.ID, v)
    }

    // 更新
    for _, v := range update.Values() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        addOrUpdateMemberAddress(v)
        store.MemberAddressStore.Store.Set(v.ID, v)
    }

    // 删除
    for _, v := range del.Values() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        delMemberAddress(v)
        store.MemberAddressStore.Store.Delete(v.ID)
    }

    return nil
}

func addOrUpdateMemberAddress(item *erp_entity.MemberAddress) {

    //查询是否存在商城表中
    memberInfo, err := shop_query.Member.
        Select().
        Where(shop_query.Member.ErpUID.Eq(item.ErpUID)).First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        slog.Error("memberAddressSync Member First err: " + err.Error())
        return
    }
    if memberInfo == nil { //关联会员未同步直接跳过
        return
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
        if er := updateMemberAddress(item, memberInfo, memberAddress); er != nil {
            slog.Error("memberAddressSync updateMemberAddress err: " + er.Error())
            return
        }
    } else {
        if er := addMemberAddress(item, memberInfo); er != nil {
            slog.Error("memberAddressSync addMemberAddress err: " + er.Error())
            return
        }
    }

}

func addMemberAddress(v *erp_entity.MemberAddress, m *shop_model.Member) error {
    memberAddressData := shop_model.MemberAddress{
        MemberID:    m.MemberID,
        SiteID:      1,
        Name:        v.RealName.String(),
        Mobile:      v.Mobile,
        Telephone:   v.Mobile,
        ProvinceID:  m.ProvinceID,
        CityID:      m.CityID,
        DistrictID:  m.DistrictID,
        Address:     v.Address.String(),
        FullAddress: m.FullAddress,
        IsDefault:   1,
    }
    return shop_query.MemberAddress.Create(&memberAddressData)
}

func updateMemberAddress(v *erp_entity.MemberAddress, m *shop_model.Member, md *shop_model.MemberAddress) error {
    memberAddressData := shop_model.MemberAddress{
        SiteID:      1,
        Name:        v.RealName.String(),
        Mobile:      v.Mobile,
        Telephone:   v.Mobile,
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

func delMemberAddress(goods *erp_entity.MemberAddress) {

}
