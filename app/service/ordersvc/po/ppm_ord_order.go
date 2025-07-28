package ordersvc

import (
	"time"
)

type PpmOrdOrder struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	OrgId           int64     `db:"org_id,omitempty" json:"orgId"`
	OutOrderNo      int64     `db:"out_order_no,omitempty" json:"outOrderNo"`
	Status          int64     `db:"status,omitempty" json:"status"`
	OrderCreateTime time.Time `db:"order_create_time,omitempty" json:"orderCreateTime"`
	PaidTime        time.Time `db:"paid_time,omitempty" json:"paidTime"`
	BuyCount        int       `db:"buy_count,omitempty" json:"buyCount"`
	EffectiveTime   int       `db:"effective_time,omitempty" json:"effectiveTime"`
	Seats           int       `db:"seats,omitempty" json:"seats"`
	TotalPrice      int64     `db:"total_price,omitempty" json:"totalPrice"`
	OrderPayPrice   int64     `db:"order_pay_price,omitempty" json:"orderPayPrice"`
	SourceChannel   string    `db:"source_channel,omitempty" json:"sourceChannel"`
	BuyType         string    `db:"buy_type,omitempty" json:"buyType"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrdOrder) TableName() string {
	return "ppm_ord_order"
}
