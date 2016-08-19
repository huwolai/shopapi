package service

import (
	"shopapi/dao"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"strings"
	"errors"
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

func  MerchantNear(longitude float64,latitude float64,appId string) ([]*dao.MerchantDetail,error)   {
	mDetail :=dao.NewMerchantDetail()
	mDetailList,err := mDetail.MerchantNear(longitude,latitude,appId)

	return mDetailList,err
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

