// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package shop_query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"app/internal/orm/shop_model"
)

func newMemberLabel(db *gorm.DB, opts ...gen.DOOption) memberLabel {
	_memberLabel := memberLabel{}

	_memberLabel.memberLabelDo.UseDB(db, opts...)
	_memberLabel.memberLabelDo.UseModel(&shop_model.MemberLabel{})

	tableName := _memberLabel.memberLabelDo.TableName()
	_memberLabel.ALL = field.NewAsterisk(tableName)
	_memberLabel.LabelID = field.NewInt32(tableName, "label_id")
	_memberLabel.SiteID = field.NewInt32(tableName, "site_id")
	_memberLabel.LabelName = field.NewString(tableName, "label_name")
	_memberLabel.CreateTime = field.NewUint(tableName, "create_time")
	_memberLabel.ModifyTime = field.NewInt32(tableName, "modify_time")
	_memberLabel.Remark = field.NewString(tableName, "remark")
	_memberLabel.Sort = field.NewInt32(tableName, "sort")
	_memberLabel.QualificationsConfig = field.NewString(tableName, "qualifications_config")

	_memberLabel.fillFieldMap()

	return _memberLabel
}

// memberLabel 会员标签
type memberLabel struct {
	memberLabelDo

	ALL                  field.Asterisk
	LabelID              field.Int32  // 标签id
	SiteID               field.Int32  // site_id
	LabelName            field.String // 标签名称
	CreateTime           field.Uint   // 创建时间
	ModifyTime           field.Int32  // 修改时间
	Remark               field.String // 备注
	Sort                 field.Int32  // 排序
	QualificationsConfig field.String // 资质配置 [{'id':1,'is_must':1}]  is_must是否必须 1是 0否

	fieldMap map[string]field.Expr
}

func (m memberLabel) Table(newTableName string) *memberLabel {
	m.memberLabelDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m memberLabel) As(alias string) *memberLabel {
	m.memberLabelDo.DO = *(m.memberLabelDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *memberLabel) updateTableName(table string) *memberLabel {
	m.ALL = field.NewAsterisk(table)
	m.LabelID = field.NewInt32(table, "label_id")
	m.SiteID = field.NewInt32(table, "site_id")
	m.LabelName = field.NewString(table, "label_name")
	m.CreateTime = field.NewUint(table, "create_time")
	m.ModifyTime = field.NewInt32(table, "modify_time")
	m.Remark = field.NewString(table, "remark")
	m.Sort = field.NewInt32(table, "sort")
	m.QualificationsConfig = field.NewString(table, "qualifications_config")

	m.fillFieldMap()

	return m
}

func (m *memberLabel) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *memberLabel) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 8)
	m.fieldMap["label_id"] = m.LabelID
	m.fieldMap["site_id"] = m.SiteID
	m.fieldMap["label_name"] = m.LabelName
	m.fieldMap["create_time"] = m.CreateTime
	m.fieldMap["modify_time"] = m.ModifyTime
	m.fieldMap["remark"] = m.Remark
	m.fieldMap["sort"] = m.Sort
	m.fieldMap["qualifications_config"] = m.QualificationsConfig
}

func (m memberLabel) clone(db *gorm.DB) memberLabel {
	m.memberLabelDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m memberLabel) replaceDB(db *gorm.DB) memberLabel {
	m.memberLabelDo.ReplaceDB(db)
	return m
}

type memberLabelDo struct{ gen.DO }

type IMemberLabelDo interface {
	gen.SubQuery
	Debug() IMemberLabelDo
	WithContext(ctx context.Context) IMemberLabelDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMemberLabelDo
	WriteDB() IMemberLabelDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMemberLabelDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMemberLabelDo
	Not(conds ...gen.Condition) IMemberLabelDo
	Or(conds ...gen.Condition) IMemberLabelDo
	Select(conds ...field.Expr) IMemberLabelDo
	Where(conds ...gen.Condition) IMemberLabelDo
	Order(conds ...field.Expr) IMemberLabelDo
	Distinct(cols ...field.Expr) IMemberLabelDo
	Omit(cols ...field.Expr) IMemberLabelDo
	Join(table schema.Tabler, on ...field.Expr) IMemberLabelDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMemberLabelDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMemberLabelDo
	Group(cols ...field.Expr) IMemberLabelDo
	Having(conds ...gen.Condition) IMemberLabelDo
	Limit(limit int) IMemberLabelDo
	Offset(offset int) IMemberLabelDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMemberLabelDo
	Unscoped() IMemberLabelDo
	Create(values ...*shop_model.MemberLabel) error
	CreateInBatches(values []*shop_model.MemberLabel, batchSize int) error
	Save(values ...*shop_model.MemberLabel) error
	First() (*shop_model.MemberLabel, error)
	Take() (*shop_model.MemberLabel, error)
	Last() (*shop_model.MemberLabel, error)
	Find() ([]*shop_model.MemberLabel, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*shop_model.MemberLabel, err error)
	FindInBatches(result *[]*shop_model.MemberLabel, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*shop_model.MemberLabel) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMemberLabelDo
	Assign(attrs ...field.AssignExpr) IMemberLabelDo
	Joins(fields ...field.RelationField) IMemberLabelDo
	Preload(fields ...field.RelationField) IMemberLabelDo
	FirstOrInit() (*shop_model.MemberLabel, error)
	FirstOrCreate() (*shop_model.MemberLabel, error)
	FindByPage(offset int, limit int) (result []*shop_model.MemberLabel, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMemberLabelDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m memberLabelDo) Debug() IMemberLabelDo {
	return m.withDO(m.DO.Debug())
}

func (m memberLabelDo) WithContext(ctx context.Context) IMemberLabelDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m memberLabelDo) ReadDB() IMemberLabelDo {
	return m.Clauses(dbresolver.Read)
}

func (m memberLabelDo) WriteDB() IMemberLabelDo {
	return m.Clauses(dbresolver.Write)
}

func (m memberLabelDo) Session(config *gorm.Session) IMemberLabelDo {
	return m.withDO(m.DO.Session(config))
}

func (m memberLabelDo) Clauses(conds ...clause.Expression) IMemberLabelDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m memberLabelDo) Returning(value interface{}, columns ...string) IMemberLabelDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m memberLabelDo) Not(conds ...gen.Condition) IMemberLabelDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m memberLabelDo) Or(conds ...gen.Condition) IMemberLabelDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m memberLabelDo) Select(conds ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m memberLabelDo) Where(conds ...gen.Condition) IMemberLabelDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m memberLabelDo) Order(conds ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m memberLabelDo) Distinct(cols ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m memberLabelDo) Omit(cols ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m memberLabelDo) Join(table schema.Tabler, on ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m memberLabelDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m memberLabelDo) RightJoin(table schema.Tabler, on ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m memberLabelDo) Group(cols ...field.Expr) IMemberLabelDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m memberLabelDo) Having(conds ...gen.Condition) IMemberLabelDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m memberLabelDo) Limit(limit int) IMemberLabelDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m memberLabelDo) Offset(offset int) IMemberLabelDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m memberLabelDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMemberLabelDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m memberLabelDo) Unscoped() IMemberLabelDo {
	return m.withDO(m.DO.Unscoped())
}

func (m memberLabelDo) Create(values ...*shop_model.MemberLabel) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m memberLabelDo) CreateInBatches(values []*shop_model.MemberLabel, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m memberLabelDo) Save(values ...*shop_model.MemberLabel) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m memberLabelDo) First() (*shop_model.MemberLabel, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.MemberLabel), nil
	}
}

func (m memberLabelDo) Take() (*shop_model.MemberLabel, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.MemberLabel), nil
	}
}

func (m memberLabelDo) Last() (*shop_model.MemberLabel, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.MemberLabel), nil
	}
}

func (m memberLabelDo) Find() ([]*shop_model.MemberLabel, error) {
	result, err := m.DO.Find()
	return result.([]*shop_model.MemberLabel), err
}

func (m memberLabelDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*shop_model.MemberLabel, err error) {
	buf := make([]*shop_model.MemberLabel, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m memberLabelDo) FindInBatches(result *[]*shop_model.MemberLabel, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m memberLabelDo) Attrs(attrs ...field.AssignExpr) IMemberLabelDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m memberLabelDo) Assign(attrs ...field.AssignExpr) IMemberLabelDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m memberLabelDo) Joins(fields ...field.RelationField) IMemberLabelDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m memberLabelDo) Preload(fields ...field.RelationField) IMemberLabelDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m memberLabelDo) FirstOrInit() (*shop_model.MemberLabel, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.MemberLabel), nil
	}
}

func (m memberLabelDo) FirstOrCreate() (*shop_model.MemberLabel, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*shop_model.MemberLabel), nil
	}
}

func (m memberLabelDo) FindByPage(offset int, limit int) (result []*shop_model.MemberLabel, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m memberLabelDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m memberLabelDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m memberLabelDo) Delete(models ...*shop_model.MemberLabel) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *memberLabelDo) withDO(do gen.Dao) *memberLabelDo {
	m.DO = *do.(*gen.DO)
	return m
}
