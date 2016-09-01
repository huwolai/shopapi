package service

import "shopapi/dao"

type UserBank struct  {
	Id int64
	AppId string
	OpenId string
	AccountName string
	BankName string
	BankNo string
}

func UserBankGet(openId,appId string) ([]*dao.UserBank,error) {

	return dao.NewUserBank().WithOpenId(openId,appId)
}

func UserBankAdd(userBank *UserBank) (*dao.UserBank,error) {
	uBank :=dao.NewUserBank()
	uBank.AccountName = userBank.AccountName
	uBank.AppId = userBank.AppId
	uBank.OpenId = userBank.OpenId
	uBank.BankName = userBank.BankName
	uBank.BankNo = userBank.BankNo

	id,err :=uBank.Insert()
	if err!=nil{
		return nil,err
	}
	uBank.Id = id
	return uBank,nil
}

func UserBankDel(id int64) error  {
	uBank :=dao.NewUserBank()

	return uBank.DeleteWithId(id)
}

func UserBankUpdate(userBank *UserBank) (*dao.UserBank,error)  {
	uBank :=dao.NewUserBank()
	uBank.AccountName = userBank.AccountName
	uBank.OpenId = userBank.OpenId
	uBank.BankName = userBank.BankName
	uBank.BankNo = userBank.BankNo
	uBank.AppId = userBank.AppId
	uBank.Id = userBank.Id

	err := uBank.UpdateWithId(userBank.Id)

	return uBank,err
}
