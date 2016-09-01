package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type UserBank struct  {
	Id int64
	AppId string
	OpenId string
	AccountName string
	BankName string
	BankNo string
	BaseDModel
}

func NewUserBank() *UserBank  {

	return &UserBank{}
}

func (self *UserBank) Insert() (int64,error)  {
	result,err :=db.NewSession().InsertInto("user_bank").Columns("app_id","open_id","account_name","bank_name","bank_no").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err := result.LastInsertId()
	return lastId,err
}

func (self *UserBank) WithOpenId(openId string,appId string) ([]*UserBank,error)  {
	var list []*UserBank
	_,err :=db.NewSession().Select("*").From("user_bank").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&list)

	return list,err
}

func (self *UserBank) DeleteWithId(id int64) error  {

	_,err :=db.NewSession().DeleteFrom("user_bank").Where("id=?",id).Exec()
	return err
}

func (self *UserBank) UpdateWithId(id int64) error  {
	_,err :=db.NewSession().Update("user_bank").Set("account_name",self.AccountName).Set("bank_name",self.BankName).Set("bank_no",self.BankNo).Exec()

	return err
}

