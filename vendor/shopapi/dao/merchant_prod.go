package dao

import "github.com/gocraft/dbr"

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