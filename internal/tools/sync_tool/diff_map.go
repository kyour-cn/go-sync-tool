package sync_tool

import (
    "app/internal/tools/safemap"
    "reflect"
)

// DiffMap 对比两个map，返回新增、更新、删除的map
func DiffMap[V any](old *safemap.Map[V], new *safemap.Map[V]) (*safemap.Map[V], *safemap.Map[V], *safemap.Map[V]) {
    add := safemap.New[V]()
    update := safemap.New[V]()
    del := safemap.New[V]()

    for _, k := range new.Keys() {
        if old.Has(k) {
            // 比对数据是否相同
            nv, _ := new.Get(k)
            ov, _ := old.Get(k)

            // 相同则跳过
            if reflect.DeepEqual(nv, ov) {
                continue
            }

            // 更新
            update.Set(k, nv)

        } else {
            nv, _ := new.Get(k)
            add.Set(k, nv)
        }
    }

    // 过滤删除数据
    for _, k := range old.Keys() {
        if !new.Has(k) {
            nv, _ := old.Get(k)
            del.Set(k, nv)
        }
    }

    return add, update, del
}
