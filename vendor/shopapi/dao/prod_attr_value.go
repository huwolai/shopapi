package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
)

type ProdAttrVal struct  {
	Id int64
	ProdId int64
	AttrKey string
	AttrValue string
	Flag string
	Json string
}

func NewProdAttrVal() *ProdAttrVal  {

	return &ProdAttrVal{}
}

func (self *ProdAttrVal) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=db.NewSession().InsertInto("prod_attr_val").Columns("attr_key","attr_value","prod_id","flag","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err := result.LastInsertId()
	return lastId,err
}

//特殊业务接口
func (self *ProdAttrVal) WithAttrKeyStock(vsearch string,attrKey string,prodId int64)  ([]*ProdAttrVal,error) {
	var prodAttrvals []*ProdAttrVal
	_,err :=db.NewSession().SelectBySql("select pv.* from prod_attr_val pv,prod_sku ps where pv.prod_id=ps.prod_id and pv.id=ps.attr_symbol_path and ps.stock=0 and pv.attr_key=? and pv.prod_id=? and pv.attr_value like '?%'",attrKey,prodId,vsearch).LoadStructs(&prodAttrvals)

	return prodAttrvals,err
}

func (self *ProdAttrVal) WithProdValue(value string,prodId int64) (*ProdAttrVal,error)  {
	var model *ProdAttrVal
	_,err :=db.NewSession().Select("*").From("prod_attr_val").Where("attr_value=?",value).Where("prod_id=?",prodId).LoadStructs(&model)

	return model,err
}