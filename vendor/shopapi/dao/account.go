package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type Account struct  {
	Id int64
	AppId string
	OpenId string
	Money float64
	Password string
	Status int
}

func NewAccount() *Account  {

	return &Account{}
}

func (self *Account) Insert() error {

	_,err :=db.NewSession().InsertInto("account").Columns("app_id","open_id","money","password","status").Record(self).Exec()

	return err
}

func (self *Account) AccountWithOpenId(openId string,appId string) (*Account,error)  {

	var account *Account
	_,err :=db.NewSession().Select("*").From("account").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&account)

	return account,err
}

func (self *Account) AccountUpdatePwd(pwd string,openId,appId string) error {

	_,err :=db.NewSession().Update("account").Set("password",pwd).Where("open_id=?",openId).Where("app_id=?",appId).Exec()

	return err
}

func (self *Account) AccountUpdateStatus(status int,openId string,appId string) error  {
	_,err :=db.NewSession().Update("account").Set("status",status).Where("open_id=?",openId).Where("app_id=?",appId).Exec()

	return err
}