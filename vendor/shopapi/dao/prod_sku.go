package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type ProdSku struct  {
	Id int64
	SkuNo string
	ProdId int64
	AppId string
	Price float64
	DisPrice float64
	AttrSymbolPath string
	Stock int
	Json string
}

func NewProdSku() *ProdSku  {

	return &ProdSku{}
}

func (self *ProdSku) InsertTx(tx *dbr.Tx) (error)  {
	_,err :=tx.InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()
	return err
}

func (self *ProdSku) Insert() (error)  {

	_,err :=db.NewSession().InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()

	return err
}

func (self *ProdSku) WithProdIdAndSymbolPath(attrSymbolPath string,prodId int64) (*ProdSku,error)  {
	var prodSku *ProdSku
	_,err :=db.NewSession().Select("*").From("prod_sku").Where("attr_symbol_path=?",attrSymbolPath).Where("prod_id=?",prodId).LoadStructs(&prodSku)

	return prodSku,err

}