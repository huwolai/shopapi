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
	SoldNum int
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
	_,err :=tx.InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","sold_num","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()
	return err
}

func (self *ProdSku) Insert() (error)  {

	_,err :=db.NewSession().InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","sold_num","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()

	return err
}

func (self *ProdSku) WithSkuNo(skuNo string) (*ProdSku,error)   {
	var prodSku *ProdSku
	_,err :=db.NewSession().Select("*").From("prod_sku").Where("sku_no=?",skuNo).LoadStructs(&prodSku)

	return prodSku,err
}

func (self *ProdSku) WithProdIdAndSymbolPath(attrSymbolPath string,prodId int64) (*ProdSku,error)  {
	var prodSku *ProdSku
	_,err :=db.NewSession().Select("*").From("prod_sku").Where("attr_symbol_path=?",attrSymbolPath).Where("prod_id=?",prodId).LoadStructs(&prodSku)

	return prodSku,err

}

//修改SKU 库存
func (self *ProdSku) UpdateStockWithSkuNo(stock int ,skuNo string) error {
	_,err :=db.NewSession().Update("prod_sku").Set("stock",stock).Where("sku_no=?",skuNo).Exec()

	return err
}
//修改SKU 库存
func (self *ProdSku) UpdateStockWithSkuNoTx(stock int ,skuNo string,tx *dbr.Tx) error {
	_,err :=tx.Update("prod_sku").Set("stock",stock).Where("sku_no=?",skuNo).Exec()

	return err
}