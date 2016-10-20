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

type ProdSkuDetail struct  {
	Id int64
	SkuNo string
	//商品标题
	Title string
	ProdId int64
	AppId string
	SoldNum int
	Price float64
	DisPrice float64
	AttrSymbolPath string
	Stock int
	Json string
}

func NewProdSkuDetail() *ProdSkuDetail  {

	return &ProdSkuDetail{}
}

func NewProdSku() *ProdSku  {

	return &ProdSku{}
}

func (self *ProdSku) InsertTx(tx *dbr.Tx) (error)  {
	_,err :=tx.InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","sold_num","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()
	return err
}

func (self *ProdSku) Insert() (int64,error)  {

	result,err :=db.NewSession().InsertInto("prod_sku").Columns("sku_no","prod_id","app_id","sold_num","price","dis_price","attr_symbol_path","stock","json").Record(self).Exec()
	if err!=nil{
		return int64(0),err
	}
	id,err :=result.LastInsertId()
	return id,err
}

func (self *ProdSku) WithSkuNo(skuNo string) (*ProdSku,error)   {
	var prodSku *ProdSku
	_,err :=db.NewSession().Select("*").From("prod_sku").Where("sku_no=?",skuNo).LoadStructs(&prodSku)

	return prodSku,err
}

func (self *ProdSkuDetail) WithSkuNo(skuNo string) (*ProdSkuDetail,error)  {
	var prodSkuDetail *ProdSkuDetail
	_,err :=db.NewSession().Select("prod_sku.*,product.title").From("prod_sku").Join("product","prod_sku.prod_id=product.id").Where("sku_no=?",skuNo).LoadStructs(&prodSkuDetail)

	return prodSkuDetail,err
}

func (self *ProdSku) SoldNumInc(num int,skuNo string,appId string,tx *dbr.Tx) error  {

	_,err :=tx.UpdateBySql("update prod_sku set sold_num=sold_num+? where sku_no=? and app_id=?",num,skuNo,appId).Exec()

	return err
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

func (self *ProdSku) UpdatePriceWithSkuNo(price float64,disPrice float64,skuNo string) error  {
	_,err :=db.NewSession().Update("prod_sku").Set("price",price).Set("dis_price",disPrice).Where("sku_no=?",skuNo).Exec()

	return err
}

func (self *ProdSku) UpdatePriceWithProdIdAndSymbolPath(price float64,disPrice float64,attrSymbolPath string,prodId int64) error  {
	_,err :=db.NewSession().Update("prod_sku").Set("price",price).Set("dis_price",disPrice).Where("attr_symbol_path=?",attrSymbolPath).Where("prod_id=?",prodId).Exec()

	return err
}

//修改SKU 库存
func (self *ProdSku) UpdateStockWithSkuNoTx(stock int ,skuNo string,tx *dbr.Tx) error {
	_,err :=tx.Update("prod_sku").Set("stock",stock).Where("sku_no=?",skuNo).Exec()

	return err
}