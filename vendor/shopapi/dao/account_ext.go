package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type AccountExt struct  {
	OpenId string
	AppId string
	Api string
	Status int
}
func NewAccountExt() *AccountExt  {

	return &AccountExt{}
}

func (self *AccountExt) Insert() error  {

	_,err :=db.NewSession().InsertInto("account_ext").Columns("open_id","app_id","api","status").Record(self).Exec()

	return err
}

func (self *AccountExt) AccountExtWithApi(api string,openId,appId string) (*AccountExt,error)  {

	var accountExt *AccountExt
	_,err :=db.NewSession().SelectBySql("select * from account_ext where open_id=? and app_id=? and api=?",openId,appId,api).LoadStructs(&accountExt)

	return accountExt,err
}

func (self *AccountExt) AccountExtUpdateStatus(status int,api string,openId,appId string)  error {
	_,err :=db.NewSession().Update("account_ext").Set("status=?",status).Where("open_id=?",openId).Where("app_id=?",appId).Where("api=?",api).Exec()

	return err
}