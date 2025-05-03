package store

import (
    "app/internal/orm/erp_entity"
    "app/internal/tools/persistence"
    "app/internal/tools/safemap"
    "golang.org/x/exp/slog"
    "os"
)

var (
    tempPath = "./temp"

    GoodsStore   *safemap.Map[*erp_entity.Goods]
    goodsStorage *persistence.Storage[[]*erp_entity.Goods]

    GoodsPriceStore   *safemap.Map[*erp_entity.GoodsPrice]
    goodsPriceStorage *persistence.Storage[[]*erp_entity.GoodsPrice]

    GoodsStockStore   *safemap.Map[*erp_entity.GoodsStock]
    goodsStockStorage *persistence.Storage[[]*erp_entity.GoodsStock]

    MemberStore   *safemap.Map[*erp_entity.Member]
    memberStorage *persistence.Storage[[]*erp_entity.Member]
)

func Init() {

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
    slog.Info("加载缓存商品数据", "num", GoodsStore.Len())

    // 初始化商品价格存储
    GoodsPriceStore = safemap.New[*erp_entity.GoodsPrice]()
    goodsPriceStorage = persistence.NewStorage[[]*erp_entity.GoodsPrice]()
    goodsPriceSlice, err := goodsPriceStorage.Load(tempPath + "/goods_price.dat")
    for _, v := range goodsPriceSlice {
        GoodsPriceStore.Set(v.GoodsErpSpid, v)
    }

    // 初始化商品库存存储
    GoodsStockStore = safemap.New[*erp_entity.GoodsStock]()
    goodsStockStorage = persistence.NewStorage[[]*erp_entity.GoodsStock]()
    goodsStockSlice, err := goodsStockStorage.Load(tempPath + "/goods_stock.dat")
    for _, v := range goodsStockSlice {
        GoodsStockStore.Set(v.GoodsErpSpid, v)
    }

    // 初始化会员存储
    MemberStore = safemap.New[*erp_entity.Member]()
    memberStorage = persistence.NewStorage[[]*erp_entity.Member]()

}

// SaveGoods 持久化商品数据
func SaveGoods() error {
    return goodsStorage.Save(GoodsStore.Values(), tempPath+"/goods.dat")
}

// SaveGoodsPrice 持久化商品价格数据
func SaveGoodsPrice() error {
    return goodsPriceStorage.Save(GoodsPriceStore.Values(), tempPath+"/goods_price.dat")
}

// SaveGoodsStock 持久化商品库存数据
func SaveGoodsStock() error {
    return goodsStockStorage.Save(GoodsStockStore.Values(), tempPath+"/goods_stock.dat")
}

// SaveMember 持久化会员数据
func SaveMember() error {
    return memberStorage.Save(MemberStore.Values(), tempPath+"/member.dat")
}
