package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
)

type  ProdAttrKey struct {
	Id int64
	ProdId int64
	AttrKey string
	AttrName string
	status int
	Flag string
	Json string
}

func NewProdAttrKey() *ProdAttrKey {

	return &ProdAttrKey{}
}

func (self*ProdAttrKey) Insert() (int64,error)  {

	result,err :=db.NewSession().InsertInto("prod_attr_key").Columns("attr_key","status","attr_name","flag","prod_id","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err :=result.LastInsertId()
	return lastId,err
}

func (self*ProdAttrKey) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("prod_attr_key").Columns("attr_key","status","attr_name","flag","prod_id","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err :=result.LastInsertId()
	return lastId,err
}

