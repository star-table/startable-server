package ordersvc

import (
	"time"
)

type PpmOrdPricePlanDing struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	GoodsName  string    `db:"goods_name,omitempty" json:"goodsName"`
	GoodsCode  string    `db:"goods_code,omitempty" json:"goodsCode"`
	ItemName   string    `db:"item_name,omitempty" json:"itemName"`
	ItemCode   string    `db:"item_code,omitempty" json:"itemCode"`
	BuyMonth   int       `db:"buy_month,omitempty" json:"buyMonth"`
	Price      int64     `db:"price,omitempty" json:"price"`
	Level      int64     `db:"level,omitempty" json:"level"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrdPricePlanDing) TableName() string {
	return "ppm_ord_price_plan_ding"
}
