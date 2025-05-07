package store

import (
    "app/internal/orm/erp_entity"
    "app/internal/tools/persistence"
    "app/internal/tools/safemap"
    "log/slog"
    "os"
)

var (
    tempPath = "./temp"

    GoodsStore                 = safemap.New[*erp_entity.Goods]()
    goodsStorage               = persistence.NewStorage[[]*erp_entity.Goods]()
    GoodsPriceStore            = safemap.New[*erp_entity.GoodsPrice]()
    goodsPriceStorage          = persistence.NewStorage[[]*erp_entity.GoodsPrice]()
    GoodsStockStore            = safemap.New[*erp_entity.GoodsStock]()
    goodsStockStorage          = persistence.NewStorage[[]*erp_entity.GoodsStock]()
    MemberStore                = safemap.New[*erp_entity.Member]()
    memberStorage              = persistence.NewStorage[[]*erp_entity.Member]()
    MemberAddressStore         = safemap.New[*erp_entity.MemberAddress]()
    memberAddressStorage       = persistence.NewStorage[[]*erp_entity.MemberAddress]()
    MemberBusinessScopeStore   = safemap.New[*erp_entity.MemberBusinessScope]()
    memberBusinessScopeStorage = persistence.NewStorage[[]*erp_entity.MemberBusinessScope]()
)

func Init() {

    // 创建temp目录，检查目录是否存在
    _, err := os.Stat(tempPath)
    if os.IsNotExist(err) {
        // 如果目录不存在，则创建
        _ = os.MkdirAll(tempPath, os.ModePerm)
    }

    // 初始化商品存储
    goodsSlice, err := goodsStorage.Load(tempPath + "/goods.dat")
    for _, v := range goodsSlice {
        GoodsStore.Set(v.GoodsErpSpid, v)
    }
    goodsSlice = nil
    slog.Debug("加载缓存商品数据", "num", GoodsStore.Len())

    // 初始化商品价格存储
    goodsPriceSlice, err := goodsPriceStorage.Load(tempPath + "/goods_price.dat")
    for _, v := range goodsPriceSlice {
        GoodsPriceStore.Set(v.GoodsErpSpid, v)
    }
    goodsPriceSlice = nil
    slog.Debug("加载缓存商品价格数据", "num", GoodsPriceStore.Len())

    // 初始化商品库存存储
    goodsStockSlice, err := goodsStockStorage.Load(tempPath + "/goods_stock.dat")
    for _, v := range goodsStockSlice {
        GoodsStockStore.Set(v.GoodsErpSpid, v)
    }
    goodsStockSlice = nil
    slog.Debug("加载缓存商品库存数据", "num", GoodsStockStore.Len())

    // 初始化会员存储
    memberSlice, err := memberStorage.Load(tempPath + "/member.dat")
    for _, v := range memberSlice {
        MemberStore.Set(v.ErpUID, v)
    }
    memberSlice = nil
    slog.Debug("加载缓存会员数据", "num", MemberStore.Len())

    // 初始化会员地址存储
    memberAddressSlice, err := memberAddressStorage.Load(tempPath + "/member_address.dat")
    for _, v := range memberAddressSlice {
        MemberAddressStore.Set(v.ID, v)
    }
    memberAddressSlice = nil
    slog.Debug("加载缓存会员地址数据", "num", MemberAddressStore.Len())

    memberBusinessScopeSlice, err := memberBusinessScopeStorage.Load(tempPath + "/member_business_scope.dat")
    for _, v := range memberBusinessScopeSlice {
        MemberBusinessScopeStore.Set(v.ID.String(), v)
    }
    memberBusinessScopeSlice = nil
    slog.Debug("加载缓存客户经营范围数据", "num", MemberBusinessScopeStore.Len())

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

// SaveMemberAddress 持久化会员地址数据
func SaveMemberAddress() error {
    return memberAddressStorage.Save(MemberAddressStore.Values(), tempPath+"/member_address.dat")
}

// SaveMemberBusinessScope 持久化客户经营范围数据
func SaveMemberBusinessScope() error {
    return memberBusinessScopeStorage.Save(MemberBusinessScopeStore.Values(), tempPath+"/member_business_scope.dat")
}
