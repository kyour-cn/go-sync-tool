package task

import (
	"app/internal/global"
	"app/internal/orm/erp_entity"
	"app/internal/orm/shop_model"
	"app/internal/orm/shop_query"
	"app/internal/store"
	"app/internal/tools/cache"
	"app/internal/tools/safemap"
	"app/internal/tools/sync_tool"
	"app/ui/apptheme"
	"errors"
	"gioui.org/layout"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"log/slog"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func NewMember() *Member {
	return &Member{}
}

// Member 同步ERP商品到商城
type Member struct{}

func (m Member) GetName() string {
	return "Member"
}

func (Member) ClearCache() error {
	return store.MemberStore.Clear()
}

func (m Member) Run(t *Task) error {
	defer func() {
		// 缓存数据到文件
		err := store.MemberStore.Save()
		if err != nil {
			slog.Error("SaveMember err: " + err.Error())
		}
	}()

	// 取出ERP全量数据
	var erpData []erp_entity.Member

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
	newMap := safemap.New[*erp_entity.Member]()
	for _, item := range erpData {
		newMap.Set(item.ErpUID, &item)
	}
	erpData = nil

	// 比对数据差异
	add, update, del := sync_tool.DiffMap[*erp_entity.Member](store.MemberStore.Store, newMap)

	slog.Info("会员同步比对", "old", store.MemberStore.Store.Len(), "new", newMap.Len(), "add", add.Len(), "update", update.Len(), "del", del.Len())
	newMap = nil

	// 统计差异总数
	t.DataCount = add.Len() + update.Len() + del.Len()

	maxConcurrent := 10

	// 新增数据处理
	err := batchProcessor(*add.GetMap(), func(v *erp_entity.Member) error {
		err := m.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberStore.Store.Set(v.ErpUID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 更新数据处理
	err = batchProcessor(*update.GetMap(), func(v *erp_entity.Member) error {
		err := m.addOrUpdate(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberStore.Store.Set(v.ErpUID, v)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)
	if err != nil {
		return err
	}

	// 删除数据处理
	err = batchProcessor(*del.GetMap(), func(v *erp_entity.Member) error {
		err := m.delete(v)
		if err != nil {
			// 这里忽略错误，否则将中断任务
			return nil
		}
		store.MemberStore.Store.Delete(v.ErpUID)
		t.DoneCount++
		return nil
	}, maxConcurrent, t.Ctx)

	return nil
}

func (m Member) addOrUpdate(item *erp_entity.Member) error {

	var memberInfo *shop_model.Member
	var err error
	//查询是否存在商城表中
	if item.MemberID != "" {
		memberId, _ := strconv.Atoi(item.MemberID)
		memberInfo, err = shop_query.Member.
			Where(shop_query.Member.MemberID.Eq(int32(memberId))).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("查询ERP返回的会员ID失败", "err", err)
			return err
		}
	}
	// 根据ERPID匹配
	if memberInfo == nil {
		memberInfo, err = shop_query.Member.
			Where(shop_query.Member.ErpUID.Eq(item.ErpUID)).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("查询ERP返回的ERPID称失败", "err", err)
			return err
		}
	}

	// 根据单位名称匹配
	if memberInfo == nil {
		memberInfo, err = shop_query.Member.
			Where(shop_query.Member.Nickname.Eq(item.Nickname.String())).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("查询ERP返回的单位名称失败", "err", err)
			return err
		}
	}

	if memberInfo != nil {
		if er := m.update(item, memberInfo); er != nil {
			slog.Error("memberSync updateMember", "err", err)
			return er
		}
	} else {
		if er := m.add(item); er != nil {
			slog.Error("memberSync addMember", "err", err)
			return er
		}
	}
	return nil
}

func (m Member) add(v *erp_entity.Member) error {
	areaInfo, _ := getAreaFormCache(v.Province.String(), v.City.String(), v.District.String())

	nowTime := int32(time.Now().Unix())

	// 会员状态
	var status int32 = 1
	if v.Status > -1 {
		status = v.Status
	}

	//查询销售员
	var salesmanId int32 = 0

	if v.SaleerID != "" {
		salesman, err := shop_query.StaffSalesman.
			Where(shop_query.StaffSalesman.ErpSaleerID.Eq(v.SaleerID)).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if salesman != nil {
			salesmanId = salesman.SalesmanID
		}
	}

	label := getMemberLabel(v.MemberType.String())

	var memberData = shop_model.Member{
		SiteID:     1,
		Username:   v.ErpUID,
		Nickname:   v.Nickname.String(),
		ErpUID:     v.ErpUID,
		SalesmanID: salesmanId,
		ErpName:    v.Nickname.String(),
		Mobile:     regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile),
		//Realname:    v.Contacts,
		MemberLabel:            label.LabelID,
		MemberLabelName:        label.LabelName,
		Password:               "e0145437e0c0a26c644efab6f97f2985",
		Status:                 status,
		RegTime:                nowTime,
		SyncTime:               nowTime,
		CanReceiveRegistergift: 1,
	}

	if areaInfo != nil {
		memberData.ProvinceID = areaInfo.ProvinceID
		memberData.CityID = areaInfo.CityID
		memberData.DistrictID = areaInfo.DistrictID
		memberData.FullAddress = areaInfo.ProvinceName + "-" + areaInfo.CityName + "-" + areaInfo.DistrictName
		memberData.Address = v.FullAddress.String()
	}

	// 是否管控
	if v.ScopeControl > -1 {
		memberData.IsNatureofbusiness = v.ScopeControl
	}

	if v.Mobile != "" && memberData.Mobile == "" {
		memberData.Mobile = regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile)
	}
	return shop_query.Member.Create(&memberData)
}

func (m Member) update(v *erp_entity.Member, member *shop_model.Member) error {
	areaInfo, _ := getAreaFormCache(v.Province.String(), v.City.String(), v.District.String())

	nowTime := int32(time.Now().Unix())

	// 要更新的字段
	var updateColumns = []field.Expr{
		shop_query.Member.SiteID,
		//shop_query.Member.Username,
		shop_query.Member.Nickname,
		shop_query.Member.ErpUID,
		shop_query.Member.ErpName,
		shop_query.Member.MemberLabel,
		shop_query.Member.MemberLabelName,
		shop_query.Member.SyncTime,
	}

	//查询销售员
	var salesmanId int32 = 0
	if v.SaleerID != "" {
		salesman, err := shop_query.StaffSalesman.
			Where(shop_query.StaffSalesman.ErpSaleerID.Eq(v.SaleerID)).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if salesman != nil {
			salesmanId = salesman.SalesmanID
			updateColumns = append(updateColumns, shop_query.Member.SalesmanID)
		}
	}

	// 手机号
	if v.Mobile != "" {
		updateColumns = append(updateColumns, shop_query.Member.Mobile)
	}

	label := getMemberLabel(v.MemberType.String())

	memberData := shop_model.Member{
		SiteID:          1,
		Username:        v.ErpUID,
		Nickname:        v.Nickname.String(),
		ErpUID:          v.ErpUID,
		ErpName:         v.Nickname.String(),
		Mobile:          regexp.MustCompile(`\b\d{11}\b`).FindString(v.Mobile),
		SalesmanID:      salesmanId,
		MemberLabel:     label.LabelID,
		MemberLabelName: label.LabelName,
		SyncTime:        nowTime,
	}

	// 状态同步
	if v.Status > -1 {
		memberData.Status = v.Status
		updateColumns = append(updateColumns, shop_query.Member.Status)
	}

	if areaInfo != nil {
		memberData.ProvinceID = areaInfo.ProvinceID
		memberData.CityID = areaInfo.CityID
		memberData.DistrictID = areaInfo.DistrictID
		memberData.FullAddress = areaInfo.ProvinceName + "-" + areaInfo.CityName + "-" + areaInfo.DistrictName

		updateColumns = append(updateColumns,
			shop_query.Member.ProvinceID,
			shop_query.Member.CityID,
			shop_query.Member.DistrictID,
			shop_query.Member.FullAddress,
		)
	}

	// 是否管控
	if v.ScopeControl > -1 {
		memberData.IsNatureofbusiness = v.ScopeControl
		updateColumns = append(updateColumns, shop_query.Member.IsNatureofbusiness)
	}

	if v.Mobile != "" && memberData.Mobile == "" {
		memberData.Mobile = regexp.MustCompile(`\b\d{7}\b`).FindString(v.Mobile)
	}

	_, err := shop_query.Member.
		Where(shop_query.Member.MemberID.Eq(member.MemberID)).
		Select(updateColumns...).
		Updates(&memberData)
	if err != nil {
		return err
	}

	//单独更新状态
	//_, err = shop_query.Member.
	//	Where(shop_query.Member.ErpUID.Eq(v.ErpUID)).
	//	Updates(map[string]interface{}{
	//		"status": status, //这里使用map更新，避免被忽略0值
	//	})
	return err
}

func (m Member) delete(_ *erp_entity.Member) error {
	return nil
}

type AreaInfo struct {
	ProvinceID   int32
	CityID       int32
	DistrictID   int32
	ProvinceName string
	CityName     string
	DistrictName string
}

func getAreaFormCache(province string, city string, district string) (*AreaInfo, error) {
	key := "area:" + province + "_" + city + "_" + district
	// 从缓存中获取
	return cache.Remember(key, 300, func() (*AreaInfo, error) {
		return getArea(province, city, district)
	})
}

// 获取用户区域id
func getArea(province string, city string, district string) (*AreaInfo, error) {
	areaInfo := AreaInfo{
		ProvinceID:   0,
		CityID:       0,
		DistrictID:   0,
		ProvinceName: "暂无",
		CityName:     "暂无",
		DistrictName: "暂无",
	}

	qa := shop_query.Area

	// 省
	provinceInfo, err := qa.Where(
		qa.Where(
			qa.Where(
				qa.Name.Like("%"+province+"%"),
			).Or(
				qa.Shortname.Like("%"+province+"%"),
			),
		),
		qa.Level.Eq(1),
		qa.Pid.Eq(0),
	).First()
	if provinceInfo != nil {
		areaInfo.ProvinceID = provinceInfo.ID
		areaInfo.ProvinceName = provinceInfo.Name
	} else {
		return nil, err
	}

	if city != "" {
		// 市
		cityInfo, _ := qa.Where(
			qa.Where(
				qa.Where(
					qa.Name.Like("%"+city+"%"),
				).Or(
					qa.Shortname.Like("%"+city+"%"),
				),
			),
			qa.Level.Eq(2),
			qa.Pid.Eq(areaInfo.ProvinceID),
		).First()
		if cityInfo != nil {
			areaInfo.CityID = cityInfo.ID
			areaInfo.CityName = cityInfo.Name
		} else {
			//fmt.Println("城市不存在：" + district)

			// 再试试县能否匹配上
			// 获取省下面的所有市
			cityList, _ := qa.Where(
				qa.Level.Eq(2),
				qa.Pid.Eq(areaInfo.ProvinceID),
			).
				Select(qa.ID).
				Find()
			var cityIds []int32
			for _, city := range cityList {
				cityIds = append(cityIds, city.ID)
			}

			// 区/县
			districtInfo, err := qa.Where(
				qa.Where(
					qa.Where(
						qa.Name.Like("%"+district+"%"),
					).Or(
						qa.Shortname.Like("%"+district+"%"),
					),
				),
				qa.Level.Eq(3),
				qa.Pid.In(cityIds...),
			).First()
			if districtInfo != nil {
				areaInfo.DistrictID = districtInfo.ID
				areaInfo.DistrictName = districtInfo.Name
			} else {
				return nil, err
			}

			// 反查市
			cityInfo, err := qa.Where(
				qa.ID.Eq(districtInfo.Pid),
			).First()
			if cityInfo != nil {
				areaInfo.CityID = cityInfo.ID
				areaInfo.CityName = cityInfo.Name
			}

			return nil, err
		}
		// 区/县
		districtInfo, err := qa.Where(
			qa.Where(
				qa.Where(
					qa.Name.Like("%"+district+"%"),
				).Or(
					qa.Shortname.Like("%"+district+"%"),
				),
			),
			qa.Level.Eq(3),
			qa.Pid.Eq(areaInfo.CityID),
		).First()
		if districtInfo != nil {
			areaInfo.DistrictID = districtInfo.ID
			areaInfo.DistrictName = districtInfo.Name
		} else {
			return nil, err
		}
	} else if district != "" { // 没有维护市，只维护区/县的情况

		// 获取省下面的所有市
		cityList, _ := qa.Where(
			qa.Level.Eq(2),
			qa.Pid.Eq(areaInfo.ProvinceID),
		).
			Select(qa.ID).
			Find()
		var cityIds []int32
		for _, city := range cityList {
			cityIds = append(cityIds, city.ID)
		}

		// 区/县
		districtInfo, _ := qa.Where(
			qa.Where(
				qa.Where(
					qa.Name.Like("%"+district+"%"),
				).Or(
					qa.Shortname.Like("%"+district+"%"),
				),
			),
			qa.Level.Eq(3),
			qa.Pid.In(cityIds...),
		).First()
		if districtInfo != nil {
			areaInfo.DistrictID = districtInfo.ID
			areaInfo.DistrictName = districtInfo.Name
		} else {
			return nil, err
		}

		// 反查市
		cityInfo, _ := qa.Where(
			qa.ID.Eq(districtInfo.Pid),
		).First()
		if cityInfo != nil {
			areaInfo.CityID = cityInfo.ID
			areaInfo.CityName = cityInfo.Name
		}
	}

	if areaInfo.CityName == "暂无" && areaInfo.DistrictName == "暂无" {
		return nil, errors.New("地址解析失败")
	}

	return &areaInfo, nil
}

// 获取用户标签id
var mutexMap sync.Map

func getMemberLabel(labelName string) *shop_model.MemberLabel {
	if labelName == "" {
		return &shop_model.MemberLabel{}
	}

	// 获取标签名对应的互斥锁（不存在时自动创建）
	mutex, _ := mutexMap.LoadOrStore(labelName, &sync.Mutex{})
	m := mutex.(*sync.Mutex)

	// 对当前标签名加锁
	m.Lock()
	defer m.Unlock()

	key := "member_label:" + labelName
	// 从缓存中获取
	label, _ := cache.Remember(key, 300, func() (label *shop_model.MemberLabel, err error) {
		label, _ = shop_query.MemberLabel.
			Where(shop_query.MemberLabel.LabelName.Eq(labelName)).
			First()
		if label == nil {
			label, err = shop_query.MemberLabel.
				Where(shop_query.MemberLabel.LabelName.Eq(labelName)).
				First()
		}
		return
	})

	if label == nil {
		lb := &shop_model.MemberLabel{
			LabelName: labelName,
			SiteID:    1,
		}
		_ = shop_query.MemberLabel.Create(lb)
		label = lb
	}

	return label
}

// ConfigLayout 任务配置UI布局
func (m Member) ConfigLayout(_ layout.Context, _ *apptheme.Theme, _ *Task) layout.Dimensions {
	return layout.Dimensions{}
}
