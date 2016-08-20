package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Merchant struct  {
	Id int64
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	Address string
	Flag string
	CoverDistance float64
	//权重
	Weight int
	//经度
	Longitude float64
	//维度
	Latitude float64
	BaseDModel
}

type MerchantDetail struct  {
	Id string
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	Address string
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

	result,err :=tx.InsertInto("merchant").Columns("name","app_id","open_id","address","longitude","latitude","status","weight","cover_distance","json","flag").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err := result.LastInsertId()
	return lastId,err
}

func (self *Merchant) MerchantWithId(id int64) (*Merchant,error)  {

	var model *Merchant
	_,err :=db.NewSession().Select("*").From("merchant").Where("id=?",id).LoadStructs(&model)

	return model,err
}

func (self *Merchant) MerchantWithOpenId(openId string,appId string) (*Merchant,error)  {

	var model *Merchant
	_,err :=db.NewSession().Select("*").From("merchant").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&model)

	return model,err
}

func (self*Merchant) MerchantExistWithOpenId(openId string,appId string) (bool,error)  {

	var count int64
	err :=db.NewSession().Select("count(*)").From("merchant").Where("open_id=?",openId).Where("app_id=?",appId).LoadValue(&count)

	if err!=nil {
		return false,err
	}
	if count>0 {
		return true,nil
	}
	return false,nil
}

func (self *Merchant) MerchantUpdateTx(merchant *Merchant,tx *dbr.Tx) error  {
	_,err :=tx.Update("merchant").Set("name",merchant.Name).Set("address",merchant.Address).Set("longitude",merchant.Longitude).Set("latitude",merchant.Latitude).Set("json",merchant.Json).Where("id=?",merchant.Id).Exec()
	return err
}

func (self *Merchant) UpdateStatus(status int,merchantId int64) error  {

	_,err :=db.NewSession().Update("merchant").Set("status",status).Where("merchant_id=?",merchantId).Exec()

	return err

}

func (self *Merchant) UpdateStatusTx(status int,merchantId int64,tx *dbr.Tx) error  {

	_,err :=tx.Update("merchant").Set("status=?",status).Where("merchant_id=?",merchantId).Exec()

	return err

}

func (self *MerchantDetail) MerchantNear(longitude float64,latitude float64,appId string) ([]*MerchantDetail,error)  {
	var mdetails []*MerchantDetail
	_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance  from merchant mt where app_id = ? and mt.status = 1 order by distance",longitude,latitude,appId).LoadStructs(&mdetails)

	return mdetails,err
}



