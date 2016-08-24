package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
)

type MerchantServiceTime struct {
	Id int64
	//商户ID
	MerchantId int64
	//服务时间
	Stime string
}

func NewMerchantServiceTime() *MerchantServiceTime  {

	return &MerchantServiceTime{}
}

func (self *MerchantServiceTime) DeleteWithMerchantIdTx(merchantId int64,tx *dbr.Tx) error {

	_,err :=tx.DeleteFrom("merchant_service_time").Where("merchant_id=?",merchantId).Exec()

	return err
}

func (self *MerchantServiceTime) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("merchant_service_time").Columns("merchant_id","stime").Record(self).Exec()

	return err
}

func (self *MerchantServiceTime) WithMerchantId(merchantId int64) ([]*MerchantServiceTime,error)  {
	var list []*MerchantServiceTime
	_,err :=db.NewSession().Select("*").From("merchant_service_time").Where("merchant_id=?",merchantId).LoadStructs(&list)

	return list,err
}