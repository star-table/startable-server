package ordersvc

import (
	"time"
)

type PpmOrdFunction struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	Name       string    `db:"name,omitempty" json:"name"`
	Code       string    `db:"code,omitempty" json:"code"`
	Type       int       `db:"type,omitempty" json:"type"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrdFunction) TableName() string {
	return "ppm_ord_function"
}
