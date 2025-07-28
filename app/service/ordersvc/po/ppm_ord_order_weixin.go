package ordersvc

import (
	"time"
)

type PpmOrdOrderWeiXin struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	OutOrgId      string    `db:"out_org_id,omitempty" json:"outOrgId"`
	OrderId       string    `db:"order_id,omitempty" json:"orderId"`
	EditionId     string    `db:"edition_id,omitempty" json:"editionId"`
	EditionName   string    `db:"edition_name,omitempty" json:"editionName"`
	OrderType     int       `db:"order_type,omitempty" json:"orderType"`
	OrderStatus   int       `db:"order_status,omitempty" json:"orderStatus"`
	UserCount     int       `db:"user_count,omitempty" json:"userCount"`
	OrderPeriod   int       `db:"order_period,omitempty" json:"orderPeriod"`
	PaidTime      time.Time `db:"paid_time,omitempty" json:"paidTime"`
	BeginTime     time.Time `db:"begin_time,omitempty" json:"beginTime"`
	EndTime       time.Time `db:"end_time,omitempty" json:"endTime"`
	OrderPayPrice int64     `db:"order_pay_price,omitempty" json:"orderPayPrice"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrdOrderWeiXin) TableName() string {
	return "ppm_ord_order_weixin"
}
