package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/orm/shop_model"
    "app/internal/orm/shop_query"
    "app/internal/store"
    tools "app/internal/tools/cache"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
    "gorm.io/gen/field"
    "gorm.io/gorm"
    "log/slog"
    "regexp"
    "strconv"
    "time"
)

// MemberSync 同步ERP商品到商城
type MemberSync struct {
}

func (g MemberSync) GetName() string {
    return "MemberSync"
}

func (g MemberSync) Run(t *Task) error {

    // 取出ERP全量数据
    var erpData []erp_entity.Member

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
    newMap := safemap.New[*erp_entity.Member]()
    for _, item := range erpData {
        newMap.Set(item.ErpUID, &item)
    }
    erpData = nil

    // 比对数据差异
    add, update, del := sync_tool.DiffMap[*erp_entity.Member](store.MemberStore, newMap)
    newMap = nil

    slog.Info("会员同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

    // 添加
    for _, v := range add.Values() {
        addOrUpdateMember(v)
        store.MemberStore.Set(v.ErpUID, v)
    }

    // 更新
    for _, v := range update.Values() {
        addOrUpdateMember(v)
        store.MemberStore.Set(v.ErpUID, v)
    }

    // 删除
    for _, v := range del.Values() {
        delMember(v)
        store.MemberStore.Delete(v.ErpUID)
    }

    // 缓存数据到文件
    err := store.SaveMember()
    if err != nil {
        return err
    }

    return nil
}

func addOrUpdateMember(item *erp_entity.Member) {
    // TODO 执行业务操作
    //v.MemberID = strings.TrimSpace(v.MemberID)

    var memberInfo *shop_model.Member
    var err error
    //查询是否存在商城表中
    if item.MemberID != "" {
        memberId, _ := strconv.Atoi(item.MemberID)
        memberInfo, err = shop_query.Member.
            Where(shop_query.Member.MemberID.Eq(int32(memberId))).
            First()
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            slog.Error("memberSync Member First", "err", err)
            return
        }
    }
    // 根据ERPID匹配
    if memberInfo == nil {
        memberInfo, err = shop_query.Member.
            Where(shop_query.Member.ErpUID.Eq(item.ErpUID)).
            First()
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            slog.Error("memberSync Member First", "err", err)
            return
        }
    }

    // 根据单位名称匹配
    if memberInfo == nil {
        memberInfo, err = shop_query.Member.
            Where(shop_query.Member.Nickname.Eq(item.Nickname.String())).
            First()
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            slog.Error("memberSync Member First", "err", err)
            return
        }
    }

    if memberInfo != nil {
        if er := updateMember(item, memberInfo); er != nil {
            slog.Error("memberSync updateMember", "err", err)
            return
        }
    } else {
        if er := addMember(item); er != nil {
            slog.Error("memberSync addMember", "err", err)
            return
        }
    }
}

func delMember(member *erp_entity.Member) {
    // TODO 执行业务操作
}

func addMember(v *erp_entity.Member) error {
    areaInfo := getAreaFormCache(v.Province.String(), v.City.String(), v.District.String())

    nowTime := int32(time.Now().Unix())

    // TODO: （可配置项）默认会员状态
    var status int32 = -0

    //查询销售员
    var salesmanId int32 = 0

    if v.SaleerID != "" {
        salesman, err := shop_query.StaffSalesman.
            Where(shop_query.StaffSalesman.ErpSaleerID.Eq(v.SaleerID)).
            First()
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        if salesman != nil {
            salesmanId = salesman.SalesmanID
        }
    }

    label := getMemberLabel(v.MemberType.String())

    var memberData = shop_model.Member{
        SiteID:     1,
        Username:   v.ErpUID,
        Nickname:   v.Nickname.String(),
        ErpUID:     v.ErpUID,
        SalesmanID: salesmanId,
        ErpName:    v.Nickname.String(),
        Mobile:     regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile),
        //Realname:    v.Contacts,
        MemberLabel:            label.LabelID,
        MemberLabelName:        label.LabelName,
        ProvinceID:             areaInfo.ProvinceID,
        CityID:                 areaInfo.CityID,
        DistrictID:             areaInfo.DistrictID,
        FullAddress:            areaInfo.ProvinceName + "-" + areaInfo.CityName + "-" + areaInfo.DistrictName,
        Password:               "e0145437e0c0a26c644efab6f97f2985",
        Status:                 status,
        RegTime:                nowTime,
        SyncTime:               nowTime,
        CanReceiveRegistergift: 1,
    }

    // 是否管控
    if v.ScopeControl > -1 {
        memberData.IsNatureofbusiness = v.ScopeControl
    }

    if v.Mobile != "" && memberData.Mobile == "" {
        memberData.Mobile = regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile)
    }
    return shop_query.Member.Create(&memberData)
}

func updateMember(v *erp_entity.Member, m *shop_model.Member) error {
    areaInfo := getAreaFormCache(v.Province.String(), v.City.String(), v.District.String())

    nowTime := int32(time.Now().Unix())

    var status int32 = 1 // 默认会员是激活状态

    // 要更新的字段
    var updateCloumns = []field.Expr{
        shop_query.Member.SiteID,
        //shop_query.Member.Username,
        shop_query.Member.Nickname,
        shop_query.Member.ErpUID,
        shop_query.Member.ErpName,
        shop_query.Member.MemberLabel,
        shop_query.Member.MemberLabelName,
        shop_query.Member.ProvinceID,
        shop_query.Member.CityID,
        shop_query.Member.DistrictID,
        shop_query.Member.FullAddress,
        shop_query.Member.SyncTime,
    }

    //查询销售员
    var salesmanId int32 = 0
    if v.SaleerID != "" {
        salesman, err := shop_query.StaffSalesman.
            Where(shop_query.StaffSalesman.ErpSaleerID.Eq(v.SaleerID)).
            First()
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        if salesman != nil {
            salesmanId = salesman.SalesmanID
            updateCloumns = append(updateCloumns, shop_query.Member.SalesmanID)
        }
    }

    // 手机号
    if v.Mobile != "" {
        updateCloumns = append(updateCloumns, shop_query.Member.Mobile)
    }

    label := getMemberLabel(v.MemberType.String())

    memberData := shop_model.Member{
        SiteID:          1,
        Username:        v.ErpUID,
        Nickname:        v.Nickname.String(),
        ErpUID:          v.ErpUID,
        ErpName:         v.Nickname.String(),
        Mobile:          regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile),
        SalesmanID:      salesmanId,
        MemberLabel:     label.LabelID,
        MemberLabelName: label.LabelName,
        ProvinceID:      areaInfo.ProvinceID,
        CityID:          areaInfo.CityID,
        DistrictID:      areaInfo.DistrictID,
        FullAddress:     areaInfo.ProvinceName + "-" + areaInfo.CityName + "-" + areaInfo.DistrictName,
        SyncTime:        nowTime,
    }

    // 是否管控
    if v.ScopeControl > -1 {
        memberData.IsNatureofbusiness = v.ScopeControl
        updateCloumns = append(updateCloumns, shop_query.Member.IsNatureofbusiness)
    }

    if v.Mobile != "" && memberData.Mobile == "" {
        memberData.Mobile = regexp.MustCompile(`\b\d{7}\b`).FindString(v.Mobile)
    }

    _, err := shop_query.Member.
        Where(shop_query.Member.MemberID.Eq(m.MemberID)).
        Select(updateCloumns...).
        Updates(&memberData)
    if err != nil {
        return err
    }

    //单独更新状态
    _, err = shop_query.Member.
        Where(shop_query.Member.ErpUID.Eq(v.ErpUID)).
        Updates(map[string]interface{}{
            "status": status, //这里使用map更新，避免被忽略0值
        })
    return err
}

type AreaInfo struct {
    ProvinceID   int32
    CityID       int32
    DistrictID   int32
    ProvinceName string
    CityName     string
    DistrictName string
}

func getAreaFormCache(province string, city string, district string) AreaInfo {
    key := "area:" + province + "_" + city + "_" + district
    // 从缓存中获取
    area, _ := tools.Remember(key, 300, func() (*AreaInfo, error) {
        area := getAreaFormCache(province, city, district)
        return &area, nil
    })
    return *area
}

// 获取用户标签id
func getMemberLabel(labelName string) *shop_model.MemberLabel {
    if labelName == "" {
        return &shop_model.MemberLabel{}
    }
    key := "member_label:" + labelName
    // 从缓存中获取
    label, _ := tools.Remember(key, 300, func() (label *shop_model.MemberLabel, err error) {
        label, _ = shop_query.MemberLabel.
            Where(shop_query.MemberLabel.LabelName.Eq(labelName)).
            First()
        if label == nil {
            label, err = shop_query.MemberLabel.
                Where(shop_query.MemberLabel.LabelName.Eq(labelName)).
                First()
        }
        return
    })

    if label == nil {
        lb := &shop_model.MemberLabel{
            LabelName: labelName,
            SiteID:    1,
        }
        _ = shop_query.MemberLabel.Create(lb)
        label = lb
    }

    return label
}
