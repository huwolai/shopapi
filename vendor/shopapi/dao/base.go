package dao

import "time"

type BaseDModel  struct{

	Id int64

	CreateTime time.Time	`json:"create_time"`
	UpdateTime time.Time	`json:"update_time"`


}