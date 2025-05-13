package main

import (
	"app/cmd/gorm/gen_tool"
	"app/cmd/gorm/methods"
	"app/cmd/gorm/tags"
	"app/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// main 模型代码生成
func main() {

	// 初始化数据库
	dbConfig, err := config.GetDBConfig("shop")
	if err != nil {
		panic(err)
	}

	mysqlDb, err := gorm.Open(mysql.Open(dbConfig.GenerateDsn()))
	if err != nil {
		panic("mysql connect failed: " + err.Error())
	}

	// 公共属性
	comOpts := []gen.ModelOpt{
		// 自动时间戳字段属性
		gen.FieldGORMTag("create_time", tags.CreateField),
		gen.FieldGORMTag("update_time", tags.UpdateField),
		gen.FieldType("create_time", "uint"),
		gen.FieldType("update_time", "uint"),

		// 软删除字段属性
		gen.FieldType("delete_time", "soft_delete.DeletedAt"),

		// Json序列化
		gen.WithMethod(methods.JsonMethod{}),
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/orm/shop_query",
		ModelPkgPath: "shop_model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(mysqlDb)

	// 使用工具生成模型
	db := gen_tool.Database{
		Generator:   g,
		ComOpts:     &comOpts,
		TablePrefix: "ydd_",
		Tables: []gen_tool.Table{
			{Name: "goods_sku"},
			{Name: "goods"},
			{Name: "config"},
			{Name: "member_address"},
			{Name: "member_label"},
			{
				Name: "member",
				Relate: &[]gen_tool.TableRelate{
					{
						TableName:  "member_qualification",
						FieldName:  "MemberQualification",
						Type:       field.HasMany,
						ForeignKey: "member_id",
						LocalKey:   "member_id",
					},
					{
						TableName:  "member_address",
						FieldName:  "MemberAddress",
						Type:       field.HasMany,
						ForeignKey: "member_id",
						LocalKey:   "member_id",
					},
				},
			},
			{Name: "area"},
			{Name: "member_business_scope"},
			{
				Name: "erp_order_outbound",
				Opts: []gen.ModelOpt{
					gen.FieldIgnoreReg("execution_time"),
				},
			},
			{Name: "member_business_scope_row"},
			{Name: "order_settlement_type"},
			{Name: "staff_salesman"},
			{Name: "erp_invoice"},
			{
				Name: "order_goods",
				Relate: &[]gen_tool.TableRelate{
					// 关联表
					{TableName: "goods", FieldName: "Goods", Type: field.HasOne, ForeignKey: "goods_id", LocalKey: "goods_id"},
				},
			},
			{
				Name: "order",
				Relate: &[]gen_tool.TableRelate{
					// 关联表
					{
						TableName:  "order_goods",
						FieldName:  "OrderGoods",
						Type:       field.HasMany,
						ForeignKey: "order_id",
						LocalKey:   "order_id",
						Relate: &[]gen_tool.TableRelate{
							{TableName: "goods", FieldName: "Goods", Type: field.HasOne, ForeignKey: "goods_id", LocalKey: "goods_id"},
						},
					},
					{TableName: "staff_salesman", FieldName: "StaffSalesman", Type: field.HasOne, ForeignKey: "member_id", LocalKey: "salesman_member_id"},
					{TableName: "order_settlement_type", FieldName: "SettlementType", Type: field.HasOne, ForeignKey: "ID", LocalKey: "settle_type_id"},
					{TableName: "member", FieldName: "Member", Type: field.HasOne, ForeignKey: "member_id", LocalKey: "member_id"},
				},
			},
		},
	}

	db.GenTable()
}
