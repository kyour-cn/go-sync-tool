package task

import (
    "app/internal/global"
    "app/internal/orm/erp_entity"
    "app/internal/store"
    "app/internal/tools/safemap"
    "app/internal/tools/sync_tool"
    "errors"
    "log/slog"
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

func addOrUpdateMember(member *erp_entity.Member) {
    // TODO 执行业务操作

}

func delMember(member *erp_entity.Member) {
    // TODO 执行业务操作
}
