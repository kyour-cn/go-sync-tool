package store

import (
	"app/internal/orm/erp_entity"
	"app/internal/tools/persistence"
	"app/internal/tools/safemap"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	tempPath = "./temp"
)

// 需要初始化的存储列表
var stores = []Loader{
	GoodsStore,
	GoodsPriceStore,
	GoodsStockStore,
	MemberStore,
	MemberAddressStore,
	MemberBusinessScopeStore,
	OrderOutboundStore,
	MemberCreditStore,
}

// 初始化各个实体存储（注意类型参数声明方式）
var (
	GoodsStore = &EntityStore[*erp_entity.Goods]{
		Store:    safemap.New[*erp_entity.Goods](), // 使用正确构造函数
		Storage:  persistence.NewStorage[[]*erp_entity.Goods](),
		KeyFunc:  func(g *erp_entity.Goods) string { return g.GoodsErpSpid }, // 直接使用指针类型
		FileName: "goods.dat",
	}

	GoodsPriceStore = &EntityStore[*erp_entity.GoodsPrice]{
		Store:    safemap.New[*erp_entity.GoodsPrice](),
		Storage:  persistence.NewStorage[[]*erp_entity.GoodsPrice](),
		KeyFunc:  func(gp *erp_entity.GoodsPrice) string { return gp.GoodsErpSpid },
		FileName: "goods_price.dat",
	}

	GoodsStockStore = &EntityStore[*erp_entity.GoodsStock]{
		Store:    safemap.New[*erp_entity.GoodsStock](),
		Storage:  persistence.NewStorage[[]*erp_entity.GoodsStock](),
		KeyFunc:  func(gs *erp_entity.GoodsStock) string { return gs.GoodsErpSpid },
		FileName: "goods_stock.dat",
	}

	MemberStore = &EntityStore[*erp_entity.Member]{
		Store:    safemap.New[*erp_entity.Member](),
		Storage:  persistence.NewStorage[[]*erp_entity.Member](),
		KeyFunc:  func(m *erp_entity.Member) string { return m.ErpUID },
		FileName: "member.dat",
	}

	MemberAddressStore = &EntityStore[*erp_entity.MemberAddress]{
		Store:    safemap.New[*erp_entity.MemberAddress](),
		Storage:  persistence.NewStorage[[]*erp_entity.MemberAddress](),
		KeyFunc:  func(ma *erp_entity.MemberAddress) string { return ma.ID },
		FileName: "member_address.dat",
	}

	MemberBusinessScopeStore = &EntityStore[*erp_entity.MemberBusinessScope]{
		Store:    safemap.New[*erp_entity.MemberBusinessScope](),
		Storage:  persistence.NewStorage[[]*erp_entity.MemberBusinessScope](),
		KeyFunc:  func(mbs *erp_entity.MemberBusinessScope) string { return mbs.ID.String() },
		FileName: "member_business_scope.dat",
	}

	OrderOutboundStore = &EntityStore[*erp_entity.OrderOutBound]{
		Store:    safemap.New[*erp_entity.OrderOutBound](),
		Storage:  persistence.NewStorage[[]*erp_entity.OrderOutBound](),
		KeyFunc:  func(mbs *erp_entity.OrderOutBound) string { return mbs.OutboundNo.String() },
		FileName: "member_business_scope.dat",
	}

	MemberCreditStore = &EntityStore[*erp_entity.MemberCredit]{
		Store:    safemap.New[*erp_entity.MemberCredit](),
		Storage:  persistence.NewStorage[[]*erp_entity.MemberCredit](),
		KeyFunc:  func(mbs *erp_entity.MemberCredit) string { return mbs.ErpUID },
		FileName: "member_credit.dat",
	}
)

// Loader 接口用于统一初始化
type Loader interface {
	Load()
}

// EntityStore 泛型数据存储结构体
type EntityStore[T any] struct {
	Store    *safemap.Map[T] // 修正类型名称
	Storage  *persistence.Storage[[]T]
	KeyFunc  func(T) string // 修正参数类型（去掉指针）
	FileName string
}

// Load 加载数据到存储
func (es *EntityStore[T]) Load() {
	filePath := filepath.Join(tempPath, es.FileName)

	// 判断文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	items, err := es.Storage.Load(filePath)
	if err != nil {
		slog.Error("加载数据失败", "entity", es.FileName, "error", err)
		return
	}

	for _, item := range items {
		es.Store.Set(es.KeyFunc(item), item)
	}

	slog.Debug("加载缓存数据", "entity", es.FileName, "num", es.Store.Len())
}

// Save 持久化存储数据
func (es *EntityStore[T]) Save() error {
	filePath := filepath.Join(tempPath, es.FileName)
	return es.Storage.Save(es.Store.Values(), filePath)
}

func (es *EntityStore[T]) Clear() error {
	es.Store.Clear()
	return es.Save()
}

func Init() {
	// 创建临时目录
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		if err := os.MkdirAll(tempPath, os.ModePerm); err != nil {
			slog.Error("创建临时目录失败", "error", err)
		}
	}

	// 统一加载所有存储
	for _, store := range stores {
		store.Load()
	}
}
