package service

import (
	"shopapi/dao"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"strings"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"fmt"
)

type MerchantDetailDLL struct  {
	Id int64
	Name string
	AppId string
	OpenId string
	//经度
	Longitude float64
	//维度
	Latitude float64
	Address string
	AddressId int64
	Mobile string
	//覆盖距离
	CoverDistance float64
	Json string
	Imgs []MerchantImgDLL
	
	ServiceArea string
	ServiceCity string
}

type MerchantImgDLL struct  {
	Id int64
	//商户ID
	MerchantId int64
	OpenId string
	AppId string
	Url string
	Flag string
	Json string
}

func MerchantProds(merchantId int64,appId string,flags []string,noflags []string)([]*dao.ProductDetail,error)  {

	productDetail :=dao.NewProductDetail()
	productDetailList,err :=productDetail.ProductListWithMerchant(merchantId,appId,flags,noflags)

	return productDetailList,err
}

func MerchantUpdate(dll *MerchantDetailDLL) error  {
	sesson := db.NewSession()
	tx,_  :=sesson.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	merchant := dao.NewMerchant()
	merchant,err := merchant.MerchantWithId(dll.Id)
	if err!=nil {
		return err
	}
	if merchant==nil {
		return errors.New("商户没找到!")
	}

	//如果选择了地址 将地址信息填充到商户信息里
	if merchant.AddressId!=0 {
		address,err :=dao.NewAddress().WithId(merchant.AddressId)
		if err!=nil{
			log.Error("查询地址错误！")
			return err
		}
		if address!=nil {
			merchant.Address = address.Address
			merchant.Latitude = address.Latitude
			merchant.Longitude = address.Longitude
		}		
	}
	fillMerchant(merchant,dll)
	err =merchant.MerchantUpdateTx(merchant,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	err = dao.NewMerchantImgs().DeleteWithMerchantIdTx(merchant.Id,merchant.AppId,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	imgs := dll.Imgs
	if imgs!=nil&&len(imgs)>0{

		for _,img:=range imgs {

			merchantImgs := dao.NewMerchantImgs()
			merchantImgs.AppId = dll.AppId
			merchantImgs.MerchantId = merchant.Id
			merchantImgs.OpenId = dll.OpenId
			merchantImgs.Json = img.Json
			merchantImgs.Flag = img.Flag
			merchantImgs.Url = img.Url
			err :=merchantImgs.InsertTx(tx)
			if err!=nil {
				tx.Rollback()
				return err
			}

		}
	}
	err =tx.Commit()
	return err

}

func MerchantWith(flags []string,noflags []string,status string,orderBy string,pageIndex uint64,pageSize uint64,appId string) ([]*dao.Merchant,error)  {

	return dao.NewMerchant().WithOrderby(flags,noflags,status,orderBy,pageIndex,pageSize,appId)
}

func MerchantCountWith(flags []string,noflags []string,status string,orderBy string,appId string) (int64,error) {

	return dao.NewMerchant().CountWith(flags,noflags,status,appId)
}

func fillMerchantImg(merchantImg *dao.MerchantImgs,dll *MerchantImgDLL)  {

	if dll.Json!="" {
		merchantImg.Json = dll.Json
	}

	if dll.Flag!="" {
		merchantImg.Flag = dll.Flag
	}

	if dll.Url!="" {
		merchantImg.Url = dll.Url
	}
}

func fillMerchant(merchant *dao.Merchant,dll *MerchantDetailDLL)  {
	if dll.ServiceArea!="" {
		merchant.ServiceArea = dll.ServiceArea
	}
	
	if dll.ServiceCity!="" {
		merchant.ServiceCity = dll.ServiceCity
	}
	
	if dll.Json!="" {
		merchant.Json = dll.Json
	}
	if dll.Name!="" {
		merchant.Name = dll.Name
	}

	if dll.Address!="" {
		merchant.Address = dll.Address
	}
	if dll.AddressId!=0 {
		merchant.AddressId = dll.AddressId
	}
	if dll.CoverDistance!=0 {
		merchant.CoverDistance = dll.CoverDistance
	}

	if dll.Latitude!=0 {
		merchant.Latitude = dll.Latitude
	}

	if dll.Longitude != 0{
		merchant.Longitude = dll.Longitude
	}

}
func MerchantAdd(dll *MerchantDetailDLL) (*MerchantDetailDLL,error)  {

	merchant := dao.NewMerchant()
	isMerchant,err :=merchant.MerchantExistWithOpenId(dll.OpenId,dll.AppId)
	if err!=nil{
		log.Error()
		return nil,errors.New("商户查询错误!")
	}
	if isMerchant {
		return nil,errors.New("已经是商户了!")
	}

	account,err :=dao.NewAccount().AccountWithOpenId(dll.OpenId,dll.AppId)
	if err!=nil{
		return nil,errors.New("用户信息未找到!")
	}

	sesson := db.NewSession()
	tx,_  :=sesson.Begin()
	defer func() {

		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()

	//如果有地址ID 将以地址ID对应的信息为准
	if dll.AddressId!=0 {
		address,err :=dao.NewAddress().WithId(dll.AddressId)
		if err!=nil {
			log.Error("地址查询失败！")
			return nil,err
		}
		if address==nil{
			log.Error("地址查询没有找到！")
			return nil,errors.New("地址没有找到！")
		}
		dll.Longitude = address.Longitude
		dll.Latitude = address.Latitude
		dll.Address = address.Address
	}

	merchant 					= dao.NewMerchant()
	merchant.Json				= dll.Json
	merchant.Name 				= dll.Name
	merchant.OpenId 			= dll.OpenId
	merchant.Mobile  			= account.Mobile
	merchant.Status 			= comm.MERCHANT_STATUS_WAIT_AUIT //申请中
	merchant.AppId 				= dll.AppId
	merchant.Longitude 			= dll.Longitude
	merchant.Latitude 			= dll.Latitude
	merchant.Address 			= dll.Address
	merchant.AddressId 			= dll.AddressId
	merchant.CoverDistance 		= dll.CoverDistance
	merchant.ServiceArea 		= dll.ServiceArea
	merchant.ServiceCity 		= dll.ServiceCity
	mid,err := merchant.InsertTx(tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	if dll.Imgs!=nil {
		for _,merchantImg :=range dll.Imgs {
			merchantImgs := dao.NewMerchantImgs()
			merchantImgs.AppId = merchantImg.AppId
			merchantImgs.MerchantId = mid
			merchantImgs.OpenId = dll.OpenId
			merchantImgs.Json = merchantImg.Json
			merchantImgs.Flag = merchantImg.Flag
			merchantImgs.Url = merchantImg.Url
			err :=merchantImgs.InsertTx(tx)
			if err!=nil {
				tx.Rollback()
				return nil,err
			}
		}
	}

	if err :=tx.Commit();err!=nil{
		tx.Rollback()

		return nil,err
	}
	dll.Id = mid

	return dll,nil

}

//审核商户
func MerchantAudit(merchantId int64,appId string,state int64,failRes string) error  {

	merchant :=dao.NewMerchant()
	merchant,err :=merchant.MerchantWithId(merchantId)
	if err!=nil{
		log.Error(err)
		return errors.New("查询商户信息错误!")
	}
	if merchant==nil{
		return errors.New("商户不存在!")
	}

	if merchant.Status!=comm.MERCHANT_STATUS_WAIT_AUIT {
		return errors.New("商户不是等待商户状态!")
	}	
	
	//审核未通过
	if  state == comm.MERCHANT_STATUS_FAIL  {
		//推送
		PushSingle(merchant.OpenId,appId,"审核未通过",failRes,"chefApplication")
		err =merchant.UpdateStatus(comm.MERCHANT_STATUS_FAIL,merchant.Id,failRes)
		if err!=nil{
			log.Error("更新商户状态失败",err)
			return err
		}
		return nil
	}	
	
	if  state != comm.MERCHANT_STATUS_NORMAL  {
		log.Info(state)
		return errors.New("商户状态错误!")
	}
	
	err =merchant.UpdateStatus(comm.MERCHANT_STATUS_NORMAL,merchant.Id,"")
	if err!=nil{
		log.Error("更新商户状态失败",err)
		return err
	}
	//推送
	PushSingle(merchant.OpenId,appId,"审核通过","审核通过","chefApplication")
	//------------- 特殊要求--------------
	//为商户添加默认产品
	err =ProdDefaultAddForMerchant(merchant)
	if err!=nil{
		log.Error("为商户添加默认商品失败",err)
		return err
	}
	
	return nil
}

//为商户添加默认商品
func ProdDefaultAddForMerchant(merchant *dao.Merchant) error {


	prodbll :=&ProdBLL{}
	prodbll.Title = "四菜一汤"
	prodbll.Description="四菜一汤"
	prodbll.AppId = merchant.AppId
	prodbll.DisPrice = 180.0
	prodbll.Price = 180.0
	prodbll.CategoryId=1
	prodbll.Flag=comm.MERCHANT_DEFAULT_PRODUCT_FLAG
	prodbll.MerchantId = merchant.Id
	_,err :=ProdAdd(prodbll)
	
	if err!=nil{
		return err
	}	
	
	//prodbll =&ProdBLL{}
	prodbll.Title = "六菜一汤"
	prodbll.Description="六菜一汤"
	prodbll.DisPrice = 240.0
	prodbll.Price = 240.0
	_,err =ProdAdd(prodbll)
	
	if err!=nil{
		return err
	}	

	return nil
}

func MerchantServiceTimeGet(merchantId int64) ([]*dao.MerchantServiceTime,error) {

	merchantServiceTime:=dao.NewMerchantServiceTime()
	return merchantServiceTime.WithMerchantId(merchantId)
}

//merchantId 商户ID  stimes服务时间(例如 0901)
func MerchantServiceTimeAdd(merchantId int64,stimes []string) error  {

	sesison :=db.NewSession()
	tx,_ :=sesison.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
		}
	}()
	if stimes!=nil&&len(stimes)>0 {
		merchantServiceTime := dao.NewMerchantServiceTime()
		err :=merchantServiceTime.DeleteWithMerchantIdTx(merchantId,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
		for _,stime :=range stimes {
			merchantServiceTime := dao.NewMerchantServiceTime()
			merchantServiceTime.MerchantId = merchantId
			merchantServiceTime.Stime = stime
			err :=merchantServiceTime.InsertTx(tx)
			if err!=nil{
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil

}

func  MerchantNear(longitude float64,latitude float64,openId string,appId string,serviceArea string,serviceCity string, pageIndex uint64, pageSize uint64) ([]*dao.MerchantDetail,error)   {
	mDetail :=dao.NewMerchantDetail()
	mDetailList,err := mDetail.MerchantNear(longitude,latitude,openId,appId,serviceArea,serviceCity,pageIndex,pageSize)
	return mDetailList,err
}
func  Merchants(longitude float64,latitude float64,openId string,appId string,serviceCity string, pageIndex uint64, pageSize uint64) ([]*dao.MerchantDetail,int64,error)   {
	mDetail 		:=dao.NewMerchantDetail()
	mDetailList,err := mDetail.Merchants(longitude,latitude,openId,appId,serviceCity,pageIndex,pageSize)
	count,err 		:= mDetail.MerchantsCount(openId,appId,serviceCity)
	return mDetailList,count,err
}
//附近商户搜索 可提供服务的厨师
func  MerchantNearSearch(longitude float64,latitude float64,openId string,appId string, pageIndex uint64, pageSize uint64, serviceTime string, serviceHour uint64,serviceArea string) ([]*dao.MerchantDetail,error)   {
	mDetail :=dao.NewMerchantDetail()
	mDetailList,err := mDetail.MerchantNearSearch(longitude,latitude,openId,appId,pageIndex,pageSize,serviceTime,serviceHour,serviceArea)

	return mDetailList,err
}

func MerchantImgWithMerchantId(merchantId int64,flags []string,appId string) ([]*dao.MerchantImgs,error)  {
	merchantimgs :=dao.NewMerchantImgs()

	return merchantimgs.MerchantImgsWithMerchantId(merchantId,flags)
}

func MerchantImgWithFlag(flags string,mopenId string,appId string)  ([]*dao.MerchantImgs,error) {

	merchantimg := dao.NewMerchantImgs()
	if flags=="" {
		return merchantimg.MerchantImgs(mopenId,appId)
	}else {
		flagsArray :=strings.Split(flags,",")
		return merchantimg.MerchantImgsWithFlag(flagsArray,mopenId,appId)
	}
}

func MerchantWithOpenId(openId,appId string) (*dao.Merchant,error)   {

	merchant := dao.NewMerchant()
	merchant,err := merchant.MerchantWithOpenId(openId,appId)

	return merchant,err
}

func MerchantWithId(id int64,appId string) (*dao.Merchant,error)  {
	merchant := dao.NewMerchant()
	merchant,err := merchant.MerchantWithId(id)

	return merchant,err
}

func MerchantWithIdDistance(id int64,appId string,longitude float64,latitude float64) (*dao.MerchantDetail,error)  {
	merchant := dao.NewMerchantDetail()
	merchant,err := merchant.MerchantWithIdDistance(id,longitude,latitude)

	return merchant,err
}

func MerchantOpenWithMerchantId(merchantId int64)  (*dao.MerchantOpen,error)  {
	merchantOpen :=  dao.NewMerchantOpen()
	merchantOpen,err :=merchantOpen.WithMerchantId(merchantId)

	return merchantOpen,err
}

//用户在线状态
func MerchantOnline(openId string,appId string)  (*dao.MerchantOnline,error)  {

	MerchantOnline :=dao.NewMerchantOnline()
	online,err := MerchantOnline.MerchantOnline(openId,appId)

	return online,err
}
//用户在线状态更改
func MerchantOnlineAndChange(openId string,appId string,status int) error  {
	merchant :=dao.NewMerchant()
	return merchant.MerchantOnlineAndChange(openId,appId,status)
}

//商户菜品图片批量命名.
func MerchantImgsWithNamed(appId string,names string) error {
	
	type NameStruct struct{
		Id 		int64	`json:"id"`
		Name 	string	`json:"name"`
	}

	var nameMap []NameStruct
	err:=util.ReadJsonByByte([]byte(names),&nameMap)
	
	if err!=nil {		
		return err
	}
	
	for _,v :=range nameMap {
		err = dao.NewMerchantImgs().MerchantImgsWithJson(v.Id,fmt.Sprintf("{\"name\":\"%s\"}",v.Name))
		if err!=nil {		
			return err
		}
	}

	return nil
}


//删除商户
func MerchantDelete(merchantId string,appId string) error  {	
	return dao.NewMerchant().MerchantDelete(merchantId,appId)
}


























