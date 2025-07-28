package ordersvc

import "time"

type PpmOrdWeiXinLicenceOrder struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	SuitId     string    `db:"suit_id,omitempty" json:"suit_id"`
	CorpId     string    `db:"corp_id,omitempty" json:"corp_id"`
	OrderId    string    `db:"order_id,omitempty" json:"order_id"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
}

func (*PpmOrdWeiXinLicenceOrder) TableName() string {
	return "ppm_ord_weixin_licence_order"
}
