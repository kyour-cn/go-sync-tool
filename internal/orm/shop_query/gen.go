// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                      = new(Query)
	Area                   *area
	Config                 *config
	ErpInvoice             *erpInvoice
	ErpOrderOutbound       *erpOrderOutbound
	Goods                  *goods
	GoodsCategory          *goodsCategory
	GoodsSku               *goodsSku
	Member                 *member
	MemberAddress          *memberAddress
	MemberBusinessScope    *memberBusinessScope
	MemberBusinessScopeRow *memberBusinessScopeRow
	MemberLabel            *memberLabel
	MemberQualification    *memberQualification
	Order                  *order
	OrderGoods             *orderGoods
	OrderSettlementType    *orderSettlementType
	StaffSalesman          *staffSalesman
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Area = &Q.Area
	Config = &Q.Config
	ErpInvoice = &Q.ErpInvoice
	ErpOrderOutbound = &Q.ErpOrderOutbound
	Goods = &Q.Goods
	GoodsCategory = &Q.GoodsCategory
	GoodsSku = &Q.GoodsSku
	Member = &Q.Member
	MemberAddress = &Q.MemberAddress
	MemberBusinessScope = &Q.MemberBusinessScope
	MemberBusinessScopeRow = &Q.MemberBusinessScopeRow
	MemberLabel = &Q.MemberLabel
	MemberQualification = &Q.MemberQualification
	Order = &Q.Order
	OrderGoods = &Q.OrderGoods
	OrderSettlementType = &Q.OrderSettlementType
	StaffSalesman = &Q.StaffSalesman
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                     db,
		Area:                   newArea(db, opts...),
		Config:                 newConfig(db, opts...),
		ErpInvoice:             newErpInvoice(db, opts...),
		ErpOrderOutbound:       newErpOrderOutbound(db, opts...),
		Goods:                  newGoods(db, opts...),
		GoodsCategory:          newGoodsCategory(db, opts...),
		GoodsSku:               newGoodsSku(db, opts...),
		Member:                 newMember(db, opts...),
		MemberAddress:          newMemberAddress(db, opts...),
		MemberBusinessScope:    newMemberBusinessScope(db, opts...),
		MemberBusinessScopeRow: newMemberBusinessScopeRow(db, opts...),
		MemberLabel:            newMemberLabel(db, opts...),
		MemberQualification:    newMemberQualification(db, opts...),
		Order:                  newOrder(db, opts...),
		OrderGoods:             newOrderGoods(db, opts...),
		OrderSettlementType:    newOrderSettlementType(db, opts...),
		StaffSalesman:          newStaffSalesman(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Area                   area
	Config                 config
	ErpInvoice             erpInvoice
	ErpOrderOutbound       erpOrderOutbound
	Goods                  goods
	GoodsCategory          goodsCategory
	GoodsSku               goodsSku
	Member                 member
	MemberAddress          memberAddress
	MemberBusinessScope    memberBusinessScope
	MemberBusinessScopeRow memberBusinessScopeRow
	MemberLabel            memberLabel
	MemberQualification    memberQualification
	Order                  order
	OrderGoods             orderGoods
	OrderSettlementType    orderSettlementType
	StaffSalesman          staffSalesman
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                     db,
		Area:                   q.Area.clone(db),
		Config:                 q.Config.clone(db),
		ErpInvoice:             q.ErpInvoice.clone(db),
		ErpOrderOutbound:       q.ErpOrderOutbound.clone(db),
		Goods:                  q.Goods.clone(db),
		GoodsCategory:          q.GoodsCategory.clone(db),
		GoodsSku:               q.GoodsSku.clone(db),
		Member:                 q.Member.clone(db),
		MemberAddress:          q.MemberAddress.clone(db),
		MemberBusinessScope:    q.MemberBusinessScope.clone(db),
		MemberBusinessScopeRow: q.MemberBusinessScopeRow.clone(db),
		MemberLabel:            q.MemberLabel.clone(db),
		MemberQualification:    q.MemberQualification.clone(db),
		Order:                  q.Order.clone(db),
		OrderGoods:             q.OrderGoods.clone(db),
		OrderSettlementType:    q.OrderSettlementType.clone(db),
		StaffSalesman:          q.StaffSalesman.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:                     db,
		Area:                   q.Area.replaceDB(db),
		Config:                 q.Config.replaceDB(db),
		ErpInvoice:             q.ErpInvoice.replaceDB(db),
		ErpOrderOutbound:       q.ErpOrderOutbound.replaceDB(db),
		Goods:                  q.Goods.replaceDB(db),
		GoodsCategory:          q.GoodsCategory.replaceDB(db),
		GoodsSku:               q.GoodsSku.replaceDB(db),
		Member:                 q.Member.replaceDB(db),
		MemberAddress:          q.MemberAddress.replaceDB(db),
		MemberBusinessScope:    q.MemberBusinessScope.replaceDB(db),
		MemberBusinessScopeRow: q.MemberBusinessScopeRow.replaceDB(db),
		MemberLabel:            q.MemberLabel.replaceDB(db),
		MemberQualification:    q.MemberQualification.replaceDB(db),
		Order:                  q.Order.replaceDB(db),
		OrderGoods:             q.OrderGoods.replaceDB(db),
		OrderSettlementType:    q.OrderSettlementType.replaceDB(db),
		StaffSalesman:          q.StaffSalesman.replaceDB(db),
	}
}

type queryCtx struct {
	Area                   IAreaDo
	Config                 IConfigDo
	ErpInvoice             IErpInvoiceDo
	ErpOrderOutbound       IErpOrderOutboundDo
	Goods                  IGoodsDo
	GoodsCategory          IGoodsCategoryDo
	GoodsSku               IGoodsSkuDo
	Member                 IMemberDo
	MemberAddress          IMemberAddressDo
	MemberBusinessScope    IMemberBusinessScopeDo
	MemberBusinessScopeRow IMemberBusinessScopeRowDo
	MemberLabel            IMemberLabelDo
	MemberQualification    IMemberQualificationDo
	Order                  IOrderDo
	OrderGoods             IOrderGoodsDo
	OrderSettlementType    IOrderSettlementTypeDo
	StaffSalesman          IStaffSalesmanDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Area:                   q.Area.WithContext(ctx),
		Config:                 q.Config.WithContext(ctx),
		ErpInvoice:             q.ErpInvoice.WithContext(ctx),
		ErpOrderOutbound:       q.ErpOrderOutbound.WithContext(ctx),
		Goods:                  q.Goods.WithContext(ctx),
		GoodsCategory:          q.GoodsCategory.WithContext(ctx),
		GoodsSku:               q.GoodsSku.WithContext(ctx),
		Member:                 q.Member.WithContext(ctx),
		MemberAddress:          q.MemberAddress.WithContext(ctx),
		MemberBusinessScope:    q.MemberBusinessScope.WithContext(ctx),
		MemberBusinessScopeRow: q.MemberBusinessScopeRow.WithContext(ctx),
		MemberLabel:            q.MemberLabel.WithContext(ctx),
		MemberQualification:    q.MemberQualification.WithContext(ctx),
		Order:                  q.Order.WithContext(ctx),
		OrderGoods:             q.OrderGoods.WithContext(ctx),
		OrderSettlementType:    q.OrderSettlementType.WithContext(ctx),
		StaffSalesman:          q.StaffSalesman.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
