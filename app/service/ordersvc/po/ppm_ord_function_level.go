package ordersvc

import (
	"time"
)

type PpmOrdFunctionLevel struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	FunctionId int64     `db:"function_id,omitempty" json:"functionId"`
	Level      int       `db:"level,omitempty" json:"level"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrdFunctionLevel) TableName() string {
	return "ppm_ord_function_level"
}

type FunctionLevelPo struct {
	FunctionId int64  `db:"function_id,omitempty" json:"functionId"`
	Code       string `db:"code,omitempty" json:"code"`
	Name       string `db:"name,omitempty" json:"name"`
	LimitInfo  string `db:"limit_info,omitempty" json:"limitInfo"` // 功能校验信息
}
