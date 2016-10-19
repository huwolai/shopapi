package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type MerchantProd struct  {
	MerchantId int64
	ProdId int64
	AppId string `json:"app_id"`
	//附加数据
	Json string
	BaseDModel
}

func NewMerchantProd() *MerchantProd  {

	return &MerchantProd{}
}

func (self *MerchantProd) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("merchant_prod").Columns("app_id","merchant_id","prod_id","json").Record(self).Exec()

	return err
}

func (self *MerchantProd) UpdateTx(prodId int64,merchantId int64,tx *dbr.Tx) error  {

	_,err :=tx.Update("merchant_prod").Set("merchant_id=?",merchantId).Where("prod_id=?",prodId).Exec()

	return err
}

//根据商品ID查询
func (self *MerchantProd) WithProdId(prodId int64,appId string) (*MerchantProd,error)  {
	var model *MerchantProd
	_,err :=db.NewSession().Select("*").From("merchant_prod").Where("prod_id=?",prodId).Where("app_id=?",appId).LoadStructs(&model)

	return model,err
}