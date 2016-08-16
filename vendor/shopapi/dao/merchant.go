package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Merchant struct  {
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	//经度
	Longitude float64
	//维度
	Latitude float64
	BaseDModel
}

type MerchantDetail struct  {
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	//经度
	Longitude float64
	//维度
	Latitude float64
	//权重
	Weight int
	//距离(单位米)
	Distance float64

}

func NewMerchantDetail() *MerchantDetail  {

	return &MerchantDetail{}
}

func NewMerchant() *Merchant  {

	return &Merchant{}
}

func (self *Merchant) InsertTx(tx *dbr.Tx) (int64,error) {

	result,err :=tx.InsertInto("merchant").Columns("name","app_id","open_id","longitude","latitude","status","weight","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err := result.LastInsertId()
	return lastId,err
}

func (self *MerchantDetail) MerchantNear(longitude float64,latitude float64,appId string) ([]*MerchantDetail,error)  {
	var mdetails *MerchantDetail
	_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance  from merchant mt where app_id = ? and getDistance(mt.longitude,latitude,?,?)<= mt.cover_distance ",appId,longitude,latitude).LoadStructs(&mdetails)

	return mdetails,err
}