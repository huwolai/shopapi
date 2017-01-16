package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"strconv"
	"math/rand"  
)

type Merchant struct  {
	Id int64
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	Address string
	AddressId int64
	Mobile string
	Flag string
	CoverDistance float64
	//权重
	Weight int
	//经度
	Longitude float64
	//维度
	Latitude float64
	
	FailRes string	
	
	Money float64
	
	BaseDModel
	
	ServiceArea string
}

type MerchantDetail struct  {
	Id string
	Name string
	AppId string
	OpenId string
	Status int
	Json string
	Address string
	AddressId int64
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
	
	ServiceArea string
}

type MerchantOnline struct  {
	Online int
	Status int
}

func NewMerchantDetail() *MerchantDetail  {

	return &MerchantDetail{}
}

func NewMerchant() *Merchant  {

	return &Merchant{}
}

func NewMerchantOnline() *MerchantOnline  {

	return &MerchantOnline{}
}

func (self *Merchant) InsertTx(tx *dbr.Tx) (int64,error) {

	result,err :=tx.InsertInto("merchant").Columns("name","mobile","app_id","open_id","address","address_id","longitude","latitude","status","weight","cover_distance","json","flag","service_area").Record(self).Exec()
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

func (self *MerchantDetail) MerchantWithIdDistance(id int64,longitude float64,latitude float64) (*MerchantDetail,error)  {

	var model *MerchantDetail
	_,err :=db.NewSession().Select(fmt.Sprintf("*,getDistance(longitude,latitude,%f,%f) distance",longitude,latitude)).From("merchant").Where("id=?",id).LoadStructs(&model)

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
	//_,err :=tx.Update("merchant").Set("name",merchant.Name).Set("address",merchant.Address).Set("address_id",merchant.AddressId).Set("longitude",merchant.Longitude).Set("latitude",merchant.Latitude).Set("cover_distance",merchant.CoverDistance).Set("json",merchant.Json).Where("id=?",merchant.Id).Exec()
	//return err
	
	builder :=db.NewSession().Update("merchant")
	builder = builder.Set("name"		,merchant.Name)
	builder = builder.Set("address"		,merchant.Address)
	builder = builder.Set("address_id"	,merchant.AddressId)
	builder = builder.Set("longitude"	,merchant.Longitude)
	builder = builder.Set("latitude"	,merchant.Latitude)	
	builder = builder.Set("json"		,merchant.Json)
	builder = builder.Set("service_area",merchant.ServiceArea)
	
	if merchant.CoverDistance>0 {
		builder = builder.Set("cover_distance",merchant.CoverDistance)
	}
	
	_,err :=builder.Where("id=?",merchant.Id).Exec()
	return err
}

func (self *Merchant) UpdateStatus(status int,merchantId int64,res string) error  {
	builder :=db.NewSession().Update("merchant").Set("status",status)
	if len(res)>0 {
		builder = builder.Set("fail_res",res)
	}
	_,err :=builder.Where("id=?",merchantId).Exec()
	return err
}

func (self *Merchant) UpdateStatusTx(status int,merchantId int64,tx *dbr.Tx) error  {

	_,err :=tx.Update("merchant").Set("status=?",status).Where("id=?",merchantId).Exec()

	return err

}

//权重递增 num 递增子数
func (self *Merchant) IncrWeightWithIdTx(num int,id int64,tx *dbr.Tx) error {
	tx.UpdateBySql("update merchant set online=4 where id=? and online in (1,3)",id).Exec()
	
	_,err :=tx.UpdateBySql("update merchant set weight=weight + ? where id=?",num,id).Exec()
	return err
}

func (self *MerchantDetail) MerchantNear(longitude float64,latitude float64,openId string,appId string,serviceArea string, pageIndex uint64, pageSize uint64) ([]*MerchantDetail,error)  {
	
	status := make([]uint64, 0)
	status = append(status,1)
	status = append(status,5)
	
	var mdetails []*MerchantDetail
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = ? and mt.open_id<>? and mt.flag<>'default' and getDistance(mt.longitude,latitude,?,?)<cover_distance order by distance limit ?,?",longitude,latitude,appId,status,openId,longitude,latitude,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)	
	_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status in ? and mt.open_id<>? and id>3 and find_in_set(?,service_area) order by status asc,weight desc limit ?,?",longitude,latitude,appId,status,openId,serviceArea,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = ? and mt.open_id<>? and mt.flag<>'default' order by distance limit ?,?",longitude,latitude,appId,status,openId,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
/* 	b:=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = ? and mt.open_id<>? and mt.flag<>'default' and find_in_set(?,service_area) order by distance limit ?,?",longitude,latitude,appId,status,openId,serviceArea,(pageIndex-1)*pageSize,pageSize)
	log.Error( b.ToSql() ) */
	
	
	//log.Error( len(mdetails) )
	
	//首页固定20个
	if uint64(len(mdetails))>pageSize {
		return mdetails,err
	}
	
	//首页补充到20个
	var mdetails20 []*MerchantDetail
	l:=pageSize-uint64(len(mdetails))	
	log.Info(l)
	
	//排除已存在的厨师
	existsId := make([]uint64, l)
	var id uint64
	for _,mDetail :=range mdetails {
		id,_ = strconv.ParseUint(mDetail.Id,10,64)
		existsId = append(existsId,id)
	}
	//log.Info(existsId)
	
	
	builder :=db.NewSession().Select(fmt.Sprintf("mt.*,getDistance(mt.longitude,latitude,%f,%f) distance,mt.cover_distance",longitude,latitude)).From("merchant mt")	
	builder = builder.Where("app_id = ?",appId)
	//builder = builder.Where("mt.status = ?",1)
	builder = builder.Where("mt.status in ?",status)
	builder = builder.Where("mt.open_id <> ?",openId)
	//builder = builder.Where("mt.flag <> ?","default")
	builder = builder.Where("mt.id > ?",3)
	builder = builder.Where("id not in ?",existsId)	
	_,err =builder.OrderBy("is_recom desc,id desc").Limit(uint64(l)).Offset(0).LoadStructs(&mdetails20)
	
	for _,mDetail :=range mdetails20 {
		mdetails = append(mdetails,mDetail)
	}
	
	return mdetails,err
}

func (self *MerchantDetail) MerchantNearCount(longitude float64,latitude float64,openId string,appId string,serviceArea string) (int64,error)  {
	
	status := make([]uint64, 0)
	status = append(status,1)
	status = append(status,5)

	total,err :=db.NewSession().SelectBySql("select count(app_id) from merchant mt where app_id = ? and mt.status in ? and mt.open_id<>? and id>3 and find_in_set(?,service_area) limit 1",appId,status,openId,serviceArea).ReturnInt64()

	return total,err
}

//附近商户搜索 可提供服务的厨师
func (self *MerchantDetail) MerchantNearSearch(longitude float64,latitude float64,openId string,appId string, pageIndex uint64, pageSize uint64, serviceTime string, serviceHour uint64,serviceArea string) ([]*MerchantDetail,error)  {

	status := make([]uint64, 0)
	status = append(status,1)
	status = append(status,5)

	var mdetails []*MerchantDetail	
	
	
	
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where app_id = ? and mt.status = 1 and mt.open_id<>? and mt.flag<>'default' and getDistance(mt.longitude,latitude,?,?)<cover_distance and id not in (SELECT merchant_prod.merchant_id from prod_attr_val,merchant_prod where merchant_prod.prod_id=prod_attr_val.prod_id and prod_attr_val.attr_key='time' and prod_attr_val.attr_value=?) order by distance limit ?,?",longitude,latitude,appId,openId,longitude,latitude,serviceTime,(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
	//_,err :=db.NewSession().SelectBySql("select mt.*,getDistance(mt.longitude,latitude,?,?) distance,mt.cover_distance from merchant mt where ?>=11 and ?<=22 and app_id = ? and mt.status = ? and mt.open_id<>? and mt.flag<>'default' and id not in (SELECT merchant_prod.merchant_id from prod_attr_val,merchant_prod where merchant_prod.prod_id=prod_attr_val.prod_id and prod_attr_val.attr_key='time' and prod_attr_val.attr_value=?) and id not in (SELECT merchant_id from merchant_service_time where stime=? ) order by distance limit ?,?",longitude,latitude,serviceHour,serviceHour,appId,status,openId,serviceTime,fmt.Sprintf("%d:00",serviceHour),(pageIndex-1)*pageSize,pageSize).LoadStructs(&mdetails)
	
	
	builder :=db.NewSession().Select(fmt.Sprintf("mt.*,getDistance(mt.longitude,latitude,%f,%f) distance,mt.cover_distance",longitude,latitude)).From("merchant mt")	
	builder = builder.Where("app_id = ?",appId)
	//builder = builder.Where("mt.status = ?",1)
	builder = builder.Where("mt.status in ?",status)
	builder = builder.Where("mt.open_id <> ?",openId)
	//builder = builder.Where("mt.flag <> ?","default")
	builder = builder.Where("mt.id > ?",3)
	builder = builder.Where("find_in_set(?,service_area)",serviceArea)
	builder = builder.Where("id not in (SELECT merchant_prod.merchant_id from prod_attr_val,merchant_prod where merchant_prod.prod_id=prod_attr_val.prod_id and prod_attr_val.attr_key='time' and prod_attr_val.attr_value=?)",serviceTime)	
	builder = builder.Where("id not in (SELECT merchant_id from merchant_service_time where stime=? )",fmt.Sprintf("%d:00",serviceHour))


	
	
	_,err :=builder.OrderBy("is_recom desc,id desc").Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&mdetails)
		
	//log.Error( builder.ToSql() )
	
	
	
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

func (self *Merchant)  WithOrderby(flags []string,noflags []string,status string,orderBy string,pageIndex uint64,pageSize uint64,appId string) ([]*Merchant,error){

	builder :=db.NewSession().Select("merchant.*").From("merchant").LeftJoin("account","account.open_id=merchant.open_id").Where("merchant.app_id=?",appId)
	
	if flags!=nil&&len(flags)>0{
		builder = builder.Where("merchant.flag in ?",flags)
	}

	if noflags!=nil&&len(noflags) >0 {
		builder = builder.Where("merchant.flag not in ?",noflags)
	}
	if status!="" {
		builder = builder.Where("merchant.status=?",status)
	}
	/* if orderBy!="" {
		builder = builder.OrderDir(orderBy,false)
	} */
	
	builder = builder.OrderDir("merchant.weight",false)
	builder = builder.OrderDir("account.money",false)
	//fmt.Println( builder.ToSql() )
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

//用户在线状态
func (self *MerchantOnline) MerchantOnline(openId string,appId string) (*MerchantOnline,error)  {

	sess := db.NewSession()
	var online *MerchantOnline

	builder :=sess.Select("online,status").From("`merchant`")
	builder = builder.Where("open_id=?",openId)
	builder = builder.Where("app_id=?",appId)
	
	//log.Error(builder.ToSql())
	
	_,err :=builder.LoadStructs(&online)
	if err!=nil{
		return nil,err
	}
	if online==nil{
		return nil,nil
	}

	return online,err
}
//用户在线状态更改
func (self *Merchant) MerchantOnlineAndChange(openId string,appId string,status int) error {
	_,err :=db.NewSession().Update("merchant").Set("online",status).Where("open_id=?",openId).Where("app_id=?",appId).Exec()
	return err
}

//厨师随机增加服务数量 0到2
func MerchantServiceAdd() error  {
	var MerchantService []*Merchant		
	builder :=db.NewSession().Select("*").From("merchant")
	builder = builder.Where("id>=?",4)	
	builder = builder.Where("id<=?",8)	
	builder.LoadStructs(&MerchantService)	
	
	x := 0
	for _,Service :=range MerchantService  {
		x=rand.Intn(2)
		db.NewSession().UpdateBySql(fmt.Sprintf("update merchant set weight=weight+%d where id=%d limit 1",x,Service.Id)).Exec()
	}
	return nil
}
//删除商户
func (self *Merchant) MerchantDelete(id string,appId string) error {
	_,err :=db.NewSession().UpdateBySql(fmt.Sprintf("insert into merchant_backup select * from merchant where app_id='%s' and id='%s' limit 1",appId,id)).Exec()
	
	if err==nil {
		_,err :=db.NewSession().UpdateBySql(fmt.Sprintf("delete from merchant where app_id='%s' and id='%s' limit 1",appId,id)).Exec()
		return err
	}
	
	return err
}










