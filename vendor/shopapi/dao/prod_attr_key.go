package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type  ProdAttrKey struct {
	Id int64
	ProdId int64
	AttrKey string
	AttrName string
	status int
	Flag string
	Json string
}


func (self*ProdAttrKey) Insert() error  {

	_,err :=db.NewSession().InsertInto("prod_attr_key").Columns("attr_key","status","attr_name","flag","prod_id","json").Record(self).Exec()

	return err
}

