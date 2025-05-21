package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_model"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"strings"
)

func NewGoods() *Goods {
	return &Goods{}
}

// Goods 同步ERP商品到商城
type Goods struct {
}

func (g Goods) GetName() string {
	return "Goods"
}

func (Goods) ClearCache() error {
	return store.GoodsStore.Clear()
}

func (g Goods) Run(t *Task) error {
	defer func() {
		// 缓存数据到文件
		err := store.GoodsStore.Save()
		if err != nil {
			slog.Error("SaveGoods err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.Goods

	erpDb, ok := global.DbPool.Get("erp")
	if !ok {
		return errors.New("获取ERP数据库连接失败")
	}

	// 执行SQL查询
	r := erpDb.Raw(t.Config.Sql).Scan(&erpData)
	if r.Error != nil {
		return r.Error
	}

	// 创建新的Map
	newMap := safemap.New[*erp_entity.Goods]()
	for _, item := range erpData {
		newMap.Set(item.GoodsErpSpid, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.Goods](store.GoodsStore.Store, newMap)

	slog.Info("商品同步比对", "old", store.GoodsStore.Store.Len(), "new", newMap.Len(), "add", add.Len(), "update", update.Len(), "del", del.Len())
	newMap = nil

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.Goods) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.Goods) error {
		err := g.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStore.Store.Set(v.GoodsErpSpid, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.Goods) error {
		err := g.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.GoodsStore.Store.Delete(v.GoodsErpSpid)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (g Goods) addOrUpdate(item *erp_entity.Goods) error {
	// 查询商城里面是否存在该商品
	shopGoodsInfo, err := shop_query.Goods.
		Where(shop_query.Goods.GoodsErpSpid.Eq(item.GoodsErpSpid)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("查询商品信息失败: " + err.Error())
		return err
	}
	if shopGoodsInfo != nil {
		if err := g.update(item, *shopGoodsInfo); err != nil {
			slog.Error("更新商品资料异常", "err", err.Error(), "ERP ID", item.GoodsErpSpid)
			return err
		}
	} else {
		if err = g.add(item); err != nil {
			slog.Error("添加商品资料异常", "err", err.Error(), "ERP ID", item.GoodsErpSpid)
			return err
		}
	}
	return nil
}

func (g Goods) delete(goods *erp_entity.Goods) error {
	// 查询商城里面是否存在该商品
	_, _ = shop_query.Goods.
		Where(shop_query.Goods.GoodsErpSpid.Eq(goods.GoodsErpSpid)).
		Select(
			shop_query.Goods.GoodsState,
			shop_query.Goods.IsDelete,
		).
		Updates(&shop_model.Goods{
			GoodsState: 0,
			IsDelete:   1,
		})
	return nil
}

func (g Goods) update(syncGoods *erp_entity.Goods, shopGoodsInfo shop_model.Goods) error {

	attrValue := attrGoods(syncGoods)
	yddGoodsData := shop_model.Goods{
		// 更新时不修改名称，因为电商可能会自定义
		//GoodsName:         syncGoods.GoodsName.String(),
		Currencyname:      syncGoods.GoodsNickname.String(), // 通用名
		GoodsAttrFormat:   attrValue,
		Unit:              syncGoods.Unit.String(),
		GoodsErpSpid:      syncGoods.GoodsErpSpid,
		BusinessScope:     syncGoods.BusinessTypeID.String(),
		BusinessScopeName: syncGoods.BusinessTypeName.String(),
		Keywords:          tools.GenFirstPinyin(syncGoods.GoodsName.String()),
		GoodsArea:         syncGoods.GoodsArea.String(),
		BarCode:           syncGoods.BarCode.String(),
		Packsize:          int32(syncGoods.MediumPackageNum), // 件装量
		Manufactor:        syncGoods.AttrFactory.String(),
		MinBuy:            int32(syncGoods.BuyMinNum),
		MaxBuy:            int32(syncGoods.BuyMaxNum),
		GoodsNum:          syncGoods.GoodsNo.String(),
		GoodsLabel:        syncGoods.IsPrescription.String(),
		YibaoType:         syncGoods.YiBaoType.String(),
		YibaoNo:           syncGoods.YiBaoNo.String(),
		IsMedicinal:       int32(syncGoods.IsMedicinal),
		IsDelete:          0,
	}
	if yddGoodsData.MinBuy == 0 {
		yddGoodsData.MinBuy = 1
	}

	extensionData := FormatGoodsAttr(attrValue)
	extensionDataStr, err := json.Marshal(extensionData)
	if err == nil {
		yddGoodsData.ExtensionData = string(extensionDataStr)
	}

	//if syncGoods.CategoryID > 0 {
	//	checkCategory(syncGoods.CategoryID, syncGoods.Category)
	//	yddGoodsData.CategoryID = "," + strconv.FormatInt(syncGoods.CategoryID, 10) + ","
	//}

	//更新goods表数据
	if _, er := shop_query.Goods.
		Where(shop_query.Goods.GoodsErpSpid.Eq(yddGoodsData.GoodsErpSpid)).
		Select(
			//shop_query.Goods.GoodsName,
			shop_query.Goods.BusinessScope,
			shop_query.Goods.BusinessScopeName,
			shop_query.Goods.Currencyname,
			shop_query.Goods.GoodsAttrFormat,
			shop_query.Goods.Unit,
			shop_query.Goods.GoodsErpSpid,
			shop_query.Goods.Keywords,
			shop_query.Goods.GoodsArea,
			shop_query.Goods.BarCode,
			shop_query.Goods.Packsize,
			shop_query.Goods.Manufactor,
			shop_query.Goods.MinBuy,
			shop_query.Goods.MaxBuy,
			shop_query.Goods.GoodsNum,
			shop_query.Goods.GoodsLabel,
			shop_query.Goods.YibaoType,
			shop_query.Goods.YibaoNo,
			shop_query.Goods.IsMedicinal,
			shop_query.Goods.ExtensionData,
			shop_query.Goods.IsDelete,
		).
		Updates(&yddGoodsData); er != nil {
		slog.Error("updateShopGoods Updates err: " + er.Error())
		return er
	}

	yddGoodsSkuData := shop_model.GoodsSku{
		Keywords: yddGoodsData.Keywords,
		//SkuName:         yddGoodsData.GoodsName,
		//GoodsName:       yddGoodsData.GoodsName,
		GoodsClassName:  yddGoodsData.GoodsClassName,
		GoodsAttrFormat: yddGoodsData.GoodsAttrFormat,
		Unit:            yddGoodsData.Unit,
		GoodsArea:       yddGoodsData.GoodsArea,
		MinBuy:          yddGoodsData.MinBuy,
		MaxBuy:          yddGoodsData.MaxBuy,
		IsDelete:        0,
		//GoodsState:      yddGoodsData.GoodsState,
	}

	//更新GoodsSku表数据
	if _, er := shop_query.GoodsSku.
		Where(shop_query.GoodsSku.GoodsID.Eq(shopGoodsInfo.GoodsID)).
		Select(
			shop_query.GoodsSku.Keywords,
			//shop_query.GoodsSku.SkuName,
			//shop_query.GoodsSku.GoodsName,
			shop_query.GoodsSku.GoodsClassName,
			shop_query.GoodsSku.GoodsAttrFormat,
			shop_query.GoodsSku.Unit,
			shop_query.GoodsSku.GoodsArea,
			shop_query.GoodsSku.MinBuy,
			shop_query.GoodsSku.MaxBuy,
			shop_query.GoodsSku.IsDelete,
		).Updates(&yddGoodsSkuData); er != nil {
		slog.Error("updateShopGoods GoodsSku update err: " + er.Error())
		return er
	}
	// 单独更新上下架状态字段 （自己实现自动上下架）
	//if yddGoodsData.GoodsArea == "HZZ00000002" {
	//	_, _ = shop_query.GoodsSku.Where(shop_query.GoodsSku.GoodsID.Eq(shopGoods.GoodsID)).Update(shop_query.GoodsSku.GoodsState, 0)
	//	_, _ = shop_query.Goods.Where(shop_query.Goods.GoodsErpSpid.Eq(yddGoodsData.GoodsErpSpid)).Update(shop_query.Goods.GoodsState, 0)
	//}
	return nil
}

func (g Goods) add(syncGoods *erp_entity.Goods) error {

	attrValue := attrGoods(syncGoods)
	yddGoodsData := shop_model.Goods{
		GoodsName:         syncGoods.GoodsName.String(),
		GoodsAttrFormat:   attrValue,
		Unit:              syncGoods.Unit.String(),
		GoodsErpSpid:      syncGoods.GoodsErpSpid,
		BusinessScope:     syncGoods.BusinessTypeID.String(),
		BusinessScopeName: syncGoods.BusinessTypeName.String(),
		Keywords:          tools.GenFirstPinyin(syncGoods.GoodsName.String()),
		GoodsArea:         syncGoods.GoodsArea.String(),
		BarCode:           syncGoods.BarCode.String(),
		GoodsClassName:    "实物商品",
		SiteID:            1,
		CategoryID:        ",1,",
		CategoryJSON:      "[\"1\"]",
		IsFreeShipping:    0,                                 // 是否免邮
		ShippingTemplate:  1,                                 //默认运费模板
		Packsize:          int32(syncGoods.MediumPackageNum), // 中包装
		Manufactor:        syncGoods.AttrFactory.String(),
		GoodsState:        0,
		IsFenxiao:         false, //是否参与分销
		MinBuy:            int32(syncGoods.BuyMinNum),
		MaxBuy:            int32(syncGoods.BuyMaxNum),
		GoodsNum:          syncGoods.GoodsNo.String(),
		GoodsLabel:        syncGoods.IsPrescription.String(),
		YibaoType:         syncGoods.YiBaoType.String(),
		YibaoNo:           syncGoods.YiBaoNo.String(),
		IsMedicinal:       int32(syncGoods.IsMedicinal),
		IsJc:              syncGoods.IsJc,
		DzjgCode:          syncGoods.DzjgCode,
		TraceabilityCode:  syncGoods.TraceabilityCode,
	}
	if yddGoodsData.MinBuy == 0 {
		yddGoodsData.MinBuy = 1
	}

	extensionData := FormatGoodsAttr(attrValue)
	extensionDataStr, err := json.Marshal(extensionData)
	if err == nil {
		yddGoodsData.ExtensionData = string(extensionDataStr)
	}

	//默认上架
	//if autoSale {
	//    yddGoodsData.GoodsState = 1
	//}

	if er := shop_query.Goods.Create(&yddGoodsData); er != nil {
		slog.Error("addShopGoods Goods Create err: " + er.Error())
		return er
	}
	yddGoodsSkuData := shop_model.GoodsSku{
		Keywords:        yddGoodsData.Keywords,
		GoodsID:         yddGoodsData.GoodsID,
		SkuName:         yddGoodsData.GoodsName,
		GoodsName:       yddGoodsData.GoodsName,
		GoodsClassName:  yddGoodsData.GoodsClassName,
		GoodsAttrFormat: yddGoodsData.GoodsAttrFormat,
		Unit:            yddGoodsData.Unit,
		SiteID:          1,
		GoodsArea:       yddGoodsData.GoodsArea,
		GoodsState:      yddGoodsData.GoodsState,
		MinBuy:          yddGoodsData.MinBuy,
		MaxBuy:          yddGoodsData.MaxBuy,
	}

	if er := shop_query.GoodsSku.Create(&yddGoodsSkuData); er != nil {
		slog.Error("addShopGoods GoodsSku Create err: " + er.Error())
		return er
	}

	//更新sku_id到主表
	if _, er := shop_query.Goods.
		Where(shop_query.Goods.GoodsID.Eq(yddGoodsData.GoodsID)).
		Update(shop_query.Goods.SkuID, yddGoodsSkuData.SkuID); er != nil {
		slog.Error("addShopGoods shop  SkuID err: " + er.Error())
		return er
	}
	return nil
}

// 返回商品json格式属性
func attrGoods(goods *erp_entity.Goods) string {
	type Attr struct {
		AttrName      string `json:"attr_name"`
		AttrValueName string `json:"attr_value_name"`
		AttrClassId   int32  `json:"attr_class_id"`
		AttrId        int32  `json:"attr_id"`
		AttrValueId   int32  `json:"attr_value_id"`
		Sort          int32  `json:"sort"`
	}
	m := []Attr{
		{
			AttrName:      "效期",
			AttrValueName: goods.AttrValidity.String(),
			AttrClassId:   -3444,
			AttrId:        -3444,
			AttrValueId:   -3444,
			Sort:          0,
		},
		{
			AttrName:      "保质期",
			AttrValueName: goods.AttrShelfLife.String(),
			AttrClassId:   -3452,
			AttrId:        -3452,
			AttrValueId:   -3452,
		},
		{
			AttrName:      "生产日期",
			AttrValueName: goods.FactoryDate.String(),
			AttrClassId:   -3445,
			AttrId:        -3445,
			AttrValueId:   -3445,
			Sort:          0,
		},
		{
			AttrName:      "规格",
			AttrValueName: goods.AttrSpecs.String(),
			AttrClassId:   -3446,
			AttrId:        -3446,
			AttrValueId:   -3446,
		},
		{
			AttrName:      "商品规格",
			AttrValueName: goods.GoodsSpecs.String(),
			AttrClassId:   -3451,
			AttrId:        -3451,
			AttrValueId:   -3451,
		},
		{
			AttrName:      "批准文号",
			AttrValueName: goods.AttrApprovalNumber.String(),
			AttrClassId:   -3447,
			AttrId:        -3447,
			AttrValueId:   -3447,
		},
		{
			AttrName:      "剂型",
			AttrValueName: goods.AttrDosageForm.String(),
			AttrClassId:   -3448,
			AttrId:        -3448,
			AttrValueId:   -3448,
		},
		{
			AttrName:      "生产厂家",
			AttrValueName: goods.AttrFactory.String(),
			AttrClassId:   -3449,
			AttrId:        -3449,
			AttrValueId:   -3449,
		},
		{
			AttrName:      "国家码",
			AttrValueName: goods.AttrCountryCode.String(),
			AttrClassId:   -3450,
			AttrId:        -3450,
			AttrValueId:   -3450,
		},
		{
			AttrName:      "产地",
			AttrValueName: goods.Place.String(),
			AttrClassId:   -3451,
			AttrId:        -3451,
			AttrValueId:   -3451,
		},
		{
			AttrName:      "产品批号",
			AttrValueName: goods.GoodsBatch.String(),
			AttrClassId:   -3453,
			AttrId:        -3453,
			AttrValueId:   -3453,
		},
	}

	bytes, _ := json.Marshal(m)
	return string(bytes)
}

type GoodsAttr struct {
	AttrName      string `json:"attr_name"`
	AttrValueName string `json:"attr_value_name"`
}

// FormatGoodsAttr 格式化商品属性
func FormatGoodsAttr(goodsAttrFormat string) map[string]string {
	returnData := map[string]string{
		"attr_specs":           "",
		"attr_factory":         "无",
		"attr_approval_number": "",
		"attr_dosage_form":     "",
		"attr_validity":        "暂无",
		"attr_production_date": "暂无",
		"attr_country_code":    "",
		"attr_place":           "",
	}

	var goodsAttrs []GoodsAttr
	err := json.Unmarshal([]byte(goodsAttrFormat), &goodsAttrs)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return returnData
	}

	for _, goodsAttr := range goodsAttrs {
		switch goodsAttr.AttrName {
		case "规格":
			returnData["attr_specs"] = goodsAttr.AttrValueName
		case "生产厂家":
			returnData["attr_factory"] = goodsAttr.AttrValueName
		case "批准文号":
			returnData["attr_approval_number"] = goodsAttr.AttrValueName
		case "剂型":
			returnData["attr_dosage_form"] = goodsAttr.AttrValueName
		case "效期":
			atv := strings.TrimSpace(goodsAttr.AttrValueName)
			if len(atv) > 10 {
				returnData["attr_validity"] = atv[:10]
			} else {
				returnData["attr_validity"] = atv
			}
		case "生产日期":
			atv := strings.TrimSpace(goodsAttr.AttrValueName)
			if len(atv) > 10 {
				returnData["attr_production_date"] = atv[:10]
			} else {
				returnData["attr_production_date"] = atv
			}
		case "国家码":
			returnData["attr_country_code"] = goodsAttr.AttrValueName
		case "产地":
			returnData["attr_place"] = goodsAttr.AttrValueName

		case "商品规格":
			returnData["attr_goods_attr"] = goodsAttr.AttrValueName
		case "产品批号":
			returnData["attr_batch_number"] = goodsAttr.AttrValueName
		case "保质期":
			returnData["attr_shelf_life"] = goodsAttr.AttrValueName
		}
	}

	return returnData
}
