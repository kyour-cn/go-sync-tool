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

// MemberBusinessScope 同步ERP商品到商城
type MemberBusinessScope struct {
    memberMap *safemap.Map[*shop_model.Member] // 用于临时缓存member数据
}

func (bs MemberBusinessScope) GetName() string {
    return "MemberBusinessScope"
}

func (bs MemberBusinessScope) Run(t *Task) error {
    defer func() {
        bs.memberMap = nil
        // 缓存数据到文件
        err := store.MemberBusinessScopeStore.Save()
        if err != nil {
            slog.Error("SaveMemberBusinessScope err: " + err.Error())
        }
    }()

    bs.memberMap = safemap.New[*shop_model.Member]()

    // 取出ERP全量数据
    var erpData []erp_entity.MemberBusinessScope

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
    newMap := safemap.New[*erp_entity.MemberBusinessScope]()
    for _, item := range erpData {
        newMap.Set(item.ID.String(), &item)
    }
    erpData = nil

    // 比对数据差异
    add, update, del := sync_tool.DiffMap[*erp_entity.MemberBusinessScope](store.MemberBusinessScopeStore.Store, newMap)
    newMap = nil

    slog.Info("商品价格同步比对", "add", add.Len(), "update", update.Len(), "del", del.Len())

    // 统计差异总数
    t.DataCount = add.Len() + update.Len() + del.Len()

    // 添加
    for _, v := range *add.GetMap() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        bs.addOrUpdate(v)
        store.MemberBusinessScopeStore.Store.Set(v.ID.String(), v)
        t.DoneCount++
    }

    // 更新
    for _, v := range *update.GetMap() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        bs.addOrUpdate(v)
        store.MemberBusinessScopeStore.Store.Set(v.ID.String(), v)
        t.DoneCount++
    }

    // 删除
    for _, v := range *del.GetMap() {
        // 优先检查退出信号
        if t.Ctx.Err() != nil {
            return nil
        }
        bs.delete(v)
        store.MemberBusinessScopeStore.Store.Delete(v.ID.String())
        t.DoneCount++
    }

    return nil
}

// 查找会员
func (bs MemberBusinessScope) findMember(erpUid string) (*shop_model.Member, error) {
    // 优先查询缓存
    if v, ok := bs.memberMap.Get(erpUid); ok {
        return v, nil
    }
    m, err := shop_query.Member.
        Where(shop_query.Member.ErpUID.Eq(erpUid)).
        Select(
            shop_query.Member.MemberID,
        ).
        First()
    if err != nil {
        return nil, err
    }
    bs.memberMap.Set(erpUid, m)
    return m, nil
}

func (bs MemberBusinessScope) addOrUpdate(item *erp_entity.MemberBusinessScope) {
    m, err := bs.findMember(item.ErpUID.String())
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        slog.Error("memberBusinessScopeSync Member First err: " + err.Error())
        return
    }
    if m == nil {
        return
    }
    //查询数据是否存在
    memberScope, err := shop_query.MemberBusinessScope.
        Where(
            shop_query.MemberBusinessScope.MemberID.Eq(m.MemberID),
            shop_query.MemberBusinessScope.Medicinetype.Eq(item.UserBusinessID.String()),
        ).First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        slog.Error("memberBusinessScopeSync MemberBusinessScope First err: " + err.Error())
        return
    }
    if memberScope != nil {
        if er := bs.updateMemberScope(item, memberScope); er != nil {
            slog.Error("memberBusinessScopeSync updateMemberScope err: " + er.Error())
            return
        }
    } else {
        if er := bs.addMemberScope(item, m); er != nil {
            slog.Error("memberBusinessScopeSync addMemberScope err: " + er.Error())
            return
        }
    }
}

func (bs MemberBusinessScope) addMemberScope(v *erp_entity.MemberBusinessScope, m *shop_model.Member) error {
    memberScopeData := shop_model.MemberBusinessScope{
        MemberID:      m.MemberID,
        ErpUID:        v.ErpUID.String(),
        BusinessScope: v.UserBusiness.String(),
        Medicinetype:  v.UserBusinessID.String(),
    }
    return shop_query.MemberBusinessScope.Create(&memberScopeData)
}

func (bs MemberBusinessScope) updateMemberScope(v *erp_entity.MemberBusinessScope, md *shop_model.MemberBusinessScope) error {
    memberScopeData := shop_model.MemberBusinessScope{
        BusinessScope: v.UserBusiness.String(),
        Medicinetype:  v.UserBusinessID.String(),
    }
    _, err := shop_query.MemberBusinessScope.
        Where(shop_query.MemberBusinessScope.ID.Eq(md.ID)).
        Select(
            shop_query.MemberBusinessScope.BusinessScope,
            shop_query.MemberBusinessScope.Medicinetype,
        ).
        Updates(&memberScopeData)
    return err
}

func (bs MemberBusinessScope) delete(item *erp_entity.MemberBusinessScope) {
    m, err := bs.findMember(item.ErpUID.String())
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        slog.Error("memberBusinessScopeSync Member First err: " + err.Error())
        return
    }
    if m == nil {
        return
    }

    _, err = shop_query.MemberBusinessScope.
        Where(
            shop_query.MemberBusinessScope.MemberID.Eq(m.MemberID),
            shop_query.MemberBusinessScope.Medicinetype.Eq(item.UserBusinessID.String()),
        ).
        Delete()
    if err != nil {
        slog.Error("memberBusinessScopeSync delete err: " + err.Error())
    }
}
