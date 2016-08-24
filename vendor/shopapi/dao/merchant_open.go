package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type MerchantOpen struct  {
	Id int64
	AppId string
	MerchantId int64
	IsOpen int
	OpenTimeStart string
	OpenTimeEnd string
	BaseDModel
}

func NewMerchantOpen() *MerchantOpen  {

	return &MerchantOpen{}
}

func (self *MerchantOpen) InsertTx(* dbr.Tx) error  {

	_,err :=db.NewSession().InsertInto("merchant_open").Columns("app_id","merchant_id","is_open","open_time_start","open_time_end").Record(self).Exec()

	return err
}

func (self *MerchantOpen) WithMerchantId(merchantId int64) (*MerchantOpen,error)  {
	var merchantOpen *MerchantOpen
	_,err :=db.NewSession().Select("*").From("merchant_open").Where("merchant_id=?",merchantId).LoadStructs(&merchantOpen)

	return merchantOpen,err
}
