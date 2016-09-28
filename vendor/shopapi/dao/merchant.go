package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
)

type Merchant struct  {
	Id int64
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	Address string
	Mobile string
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
	//覆盖范围
	CoverDistance float64
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

	result,err :=tx.InsertInto("merchant").Columns("name","mobile","app_id","open_id","address","longitude","latitude","status","weight","cover_distance","json","flag").Record(self).Exec()
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
	_,err :=db.NewSession().Update("merchant").Set("status",status).Where("id=?",merchantId).Exec()
	return err
}

func (self *Merchant) UpdateStatusTx(status int,merchantId int64,tx *dbr.Tx) error  {

	_,err :=tx.Update("merchant").Set("status=?",status).Where("id=?",merchantId).Exec()

	return err

}

//权重递增 num 递增子数
func (self *Merchant) IncrWeightWithIdTx(num int,id int64,tx *dbr.Tx) error {
	_,err :=tx.UpdateBySql("update merchant set weight=weight + ? where id=?",num,id).Exec()

	return err
}

func (self *MerchantDetail) MerchantNear(longitude float64,latitude float64,openId string,appId string, pageIndex uint64, pageSize uint64) ([]*MerchantDetail,error)  {
	var mdetails []*MerchantDetail
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = 1 and mt.open_id<>? and mt.flag<>'default' and getDistance(mt.longitude,latitude,?,?)<cover_distance order by distance limit ?,?",longitude,latitude,appId,openId,longitude,latitude,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)	
	_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = 1 and mt.open_id<>? and mt.flag<>'default' order by distance limit ?,?",longitude,latitude,appId,openId,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)	
	
	/* builder :=db.NewSession().Select(fmt.Sprintf("mt.*,getDistance(mt.longitude,latitude,%f,%f) distance,mt.cover_distance",longitude,latitude)).From("merchant mt")
	
	builder = builder.Where("app_id = ?",appId)
	builder = builder.Where("mt.status = ?",1)
	builder = builder.Where("mt.open_id <> ?",openId)
	builder = builder.Where("mt.flag <> ?","default")
	builder = builder.Where("getDistance(mt.longitude,latitude,?,?)<cover_distance",longitude,latitude)
	
	_,err :=builder.OrderBy("distance").Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&mdetails) */
	
	return mdetails,err
}
//附近商户搜索 可提供服务的厨师
func (self *MerchantDetail) MerchantNearSearch(longitude float64,latitude float64,openId string,appId string, pageIndex uint64, pageSize uint64, serviceTime string, serviceHour uint64) ([]*MerchantDetail,error)  {
	var mdetails []*MerchantDetail	
	
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = 1 and mt.open_id<>? and mt.flag<>'default' and getDistance(mt.longitude,latitude,?,?)<cover_distance and id not in (SELECT merchant_prod.merchant_id from prod_attr_val,merchant_prod where merchant_prod.prod_id=prod_attr_val.prod_id and prod_attr_val.attr_key='time' and prod_attr_val.attr_value=?) order by distance limit ?,?",longitude,latitude,appId,openId,longitude,latitude,serviceTime,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
	_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = 1 and mt.open_id<>? and mt.flag<>'default' and id not in (SELECT merchant_prod.merchant_id from prod_attr_val,merchant_prod where merchant_prod.prod_id=prod_attr_val.prod_id and prod_attr_val.attr_key='time' and prod_attr_val.attr_value=?) and id not in (SELECT merchant_id from merchant_service_time where stime=? ) order by distance limit ?,?",longitude,latitude,appId,openId,serviceTime,fmt.Sprintf("%d:00",serviceHour),(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
	
	return mdetails,err
}

func (self *Merchant)  With(flags []string,noflags []string,status string,orderBy string,pageIndex uint64,pageSize uint64,appId string) ([]*Merchant,error){

	builder :=db.NewSession().Select("*").From("merchant").Where("app_id=?",appId)
	if flags!=nil&&len(flags)>0{
		builder = builder.Where("flag in ?",flags)
	}

	if noflags!=nil&&len(noflags) >0 {
		builder = builder.Where("flag not in ?",noflags)
	}
	if status!="" {
		builder = builder.Where("status=?",status)
	}
	if orderBy!="" {
		builder = builder.OrderDir(orderBy,false)
	}
	var list []*Merchant
	_,err :=builder.Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&list)

	return list,err
}

func (self *Merchant) CountWith(flags []string,noflags []string,status string,appId string) (int64,error)  {
	builder :=db.NewSession().Select("count(*)").From("merchant").Where("app_id=?",appId)
	if flags!=nil&&len(flags)>0{
		builder = builder.Where("flag in ?",flags)
	}

	if noflags!=nil&&len(noflags) >0 {
		builder = builder.Where("flag not in ?",noflags)
	}
	if status!="" {
		builder = builder.Where("status=?",status)
	}
	var count int64
	err :=builder.LoadValue(&count)

	return count,err
}



