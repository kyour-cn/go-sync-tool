package store

import (
    "app/internal/orm/erp_entity"
    "app/internal/tools/persistence"
    "app/internal/tools/safemap"
    "os"
)

var (
    tempPath = "./temp"

    GoodsStore   *safemap.Map[*erp_entity.Goods]
    goodsStorage *persistence.Storage[[]*erp_entity.Goods]
)

func init() {

    // 创建temp目录，检查目录是否存在
    _, err := os.Stat(tempPath)
    if os.IsNotExist(err) {
        // 如果目录不存在，则创建
        _ = os.MkdirAll(tempPath, os.ModePerm)
    }

    // 初始化商品存储
    GoodsStore = safemap.New[*erp_entity.Goods]()
    goodsStorage = persistence.NewStorage[[]*erp_entity.Goods]()
    goodsSlice, err := goodsStorage.Load(tempPath + "/goods.dat")
    for _, v := range goodsSlice {
        GoodsStore.Set(v.GoodsErpSpid, v)
    }

}

// SaveGoods 持久化商品数据
func SaveGoods() error {
    return goodsStorage.Save(GoodsStore.Values(), tempPath+"/goods.dat")
}
