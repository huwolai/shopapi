package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type AccountRecharge struct  {
	Id int64
	No string
	AppId string
	OpenId string
	Amount float64
	Status int
	From int
	Flag string
	Json string
	BaseDModel
}

func NewAccountRecharge() *AccountRecharge {

	return &AccountRecharge{}
}

func (self *AccountRecharge) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("account_recharge").Columns("no","app_id","open_id","amount","status","flag","json","from").Record(self).Exec()

	return err
}

func (self *AccountRecharge) Insert() error  {

	_,err :=db.NewSession().InsertInto("account_recharge").Columns("no","app_id","open_id","amount","status","flag","json").Record(self).Exec()

	return err
}


func (self *AccountRecharge) WithNo(no string,appId string) (*AccountRecharge,error)  {
	var model *AccountRecharge
	_,err :=db.NewSession().Select("*").From("account_recharge").Where("no=?",no).Where("app_id=?",appId).LoadStructs(&model)

	return model,err
}

func (self *AccountRecharge) UpdateStatusWithNo(status int,no string,appId string) error {

	_,err :=db.NewSession().Update("account_recharge").Set("status",status).Where("no=?",no).Where("app_id=?",appId).Exec()

	return err
}
