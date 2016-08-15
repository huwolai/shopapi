package dao

import "time"

type BaseDModel  struct{

	Id int64

	CreateTime time.Time
	UpdateTime time.Time


}