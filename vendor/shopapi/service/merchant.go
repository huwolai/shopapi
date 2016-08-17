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
	//经度
	Longitude float64
	//维度
	Latitude float64
	Address string
	//覆盖距离
	CoverDistance float64
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
	merchant.Longitude = dll.Longitude
	merchant.Latitude = dll.Latitude
	merchant.Address = dll.Address
	merchant.CoverDistance = dll.CoverDistance
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

func  MerchantNear(longitude float64,latitude float64,appId string) ([]*dao.MerchantDetail,error)   {
	mDetail :=dao.NewMerchantDetail()
	mDetailList,err := mDetail.MerchantNear(longitude,latitude,appId)

	return mDetailList,err
}

