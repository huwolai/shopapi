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
	attrValues []ProdAttrVal
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

func (self *ProdAttrKey) DetailWithProdId(prodId int64) ([]*ProdAttrKey,error)  {
	var prodAttrKeys []*ProdAttrKey
	_,err :=db.NewSession().Select("*").From("prod_attr_key").Where("prod_id=?",prodId).LoadStructs(&prodAttrKeys)

	if prodAttrKeys!=nil&&len(prodAttrKeys)>0 {
		err =fillAttrValues(prodAttrKeys,prodId)
		if err!=nil{
			return nil,err
		}
	}

	return prodAttrKeys,err
}

func fillAttrValues(prodAttrKeys []*ProdAttrKey,prodId int64) error  {
	prodAttrVal := NewProdAttrVal()
	prodAttrVals,err :=prodAttrVal.WithProdId(prodId)
	if err!=nil{
		return err
	}
	if prodAttrVals==nil{
		return nil
	}
	for _,prodAttrK :=range prodAttrKeys {
		for _,prodAttrV :=range prodAttrVals {
			if prodAttrK.AttrKey==prodAttrV.AttrKey {
				attrValues :=prodAttrK.attrValues
				if attrValues==nil {
					attrValues = make([]ProdAttrVal,0)
				}
				attrValues = append(attrValues,*prodAttrV)
				prodAttrK.attrValues = attrValues
			}
		}

	}
	return nil
}


