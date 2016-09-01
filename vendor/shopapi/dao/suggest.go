package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type Suggest struct  {
	Id int64
	OpenId string
	Contact string
	Content string
	BaseDModel
}

func NewSuggest() *Suggest  {

	return &Suggest{}
}

func (self *Suggest) Insert() error  {

	_,err :=db.NewSession().InsertInto("suggest").Columns("open_id","contact","content").Record(self).Exec()

	return err
}
