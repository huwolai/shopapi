package service

import (
	"shopapi/dao"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"strings"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
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
	//覆盖距离
	CoverDistance float64
	Json string
	Imgs []MerchantImgDLL
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
	fillMerchant(merchant,dll)
	err =merchant.MerchantUpdateTx(merchant,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	imgs := dll.Imgs
	if imgs!=nil{
		for _,img:=range imgs {

			if img.Id!=0 {
				merchantImg := dao.NewMerchantImgs()
				merchantImg,err := merchantImg.MerchantImgsWithId(img.Id)
				if err!=nil{
					tx.Rollback()
					return err
				}

				if merchantImg!=nil {
					fillMerchantImg(merchantImg,&img)
					err :=merchantImg.MerchantImgsUpdateTx(merchantImg,tx)
					if err!=nil{
						tx.Rollback()
						return err
					}
				}
			}else{
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
	}

	err =tx.Commit()
	return err

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

	if dll.Json!="" {
		merchant.Json = dll.Json
	}
	if dll.Name!="" {
		merchant.Name = dll.Name
	}

	if dll.Address!="" {
		merchant.Address = dll.Address
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

	sesson := db.NewSession()
	tx,_  :=sesson.Begin()
	defer func() {

		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()

	merchant = dao.NewMerchant()
	merchant.Json=dll.Json
	merchant.Name = dll.Name
	merchant.OpenId = dll.OpenId
	merchant.Status = comm.MERCHANT_STATUS_WAIT_AUIT //等待审核
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
func MerchantAudit(merchantId int64,appId string) error  {

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

	err =merchant.UpdateStatus(comm.MERCHANT_STATUS_NORMAL,merchant.Id)
	if err!=nil{
		log.Error("更新商户状态失败",err)
		return err
	}

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
	prodbll.Title = "默认商品"
	prodbll.Description="默认商品描述"
	prodbll.AppId = merchant.AppId
	prodbll.DisPrice = 158.0
	prodbll.Price = 158.0
	prodbll.CategoryId=1
	prodbll.Flag=comm.MERCHANT_DEFAULT_PRODUCT_FLAG
	prodbll.MerchantId = merchant.Id
	_,err :=ProdAdd(prodbll)

	return err
}

func  MerchantNear(longitude float64,latitude float64,appId string) ([]*dao.MerchantDetail,error)   {
	mDetail :=dao.NewMerchantDetail()
	mDetailList,err := mDetail.MerchantNear(longitude,latitude,appId)

	return mDetailList,err
}

func MerchantImgWithMerchantId(merchantId int64,flags []string,appId string) ([]*dao.MerchantImgs,error) {
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

