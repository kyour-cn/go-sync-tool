package store

import (
    "app/internal/orm/erp_entity"
    "app/internal/tools/safemap"
)

var GoodsStore *safemap.Map[*erp_entity.Goods]

func init() {
    GoodsStore = safemap.New[*erp_entity.Goods]()
}
