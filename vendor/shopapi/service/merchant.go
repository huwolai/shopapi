package service

import (
	"shopapi/dao"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type MerchantAddDLL struct  {
	Id int64
	Name string
	AppId string
	OpenId string
	Json string
}
func MerchantAdd(dll *MerchantAddDLL) (*MerchantAddDLL,error)  {
	sesson := db.NewSession()
	tx,_  :=sesson.Begin()
	defer func() {

		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()


	merchant := dao.NewMerchant()
	merchant.Json=dll.Json
	merchant.Name = dll.Name
	merchant.OpenId = dll.OpenId
	merchant.Status = comm.MERCHANT_STATUS_NORMAL
	merchant.AppId = dll.AppId
	mid,err := merchant.InsertTx(tx)
	if err!=nil{
		tx.Rollback()

		return nil,err

	}
	if err :=tx.Commit();err!=nil{
		tx.Rollback()

		return nil,err
	}
	dll.Id = mid

	return dll,nil

}

