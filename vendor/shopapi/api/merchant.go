package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"net/http"
	"shopapi/dao"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/security"
)


type MerchantDto struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Status int `json:"status"`
	Json string `json:"json"`
	Address string `json:"address"`
	CoverDistance float64 `json:"cover_distance"`
	//经度
	Longitude float64 `json:"longitude"`
	//维度
	Latitude float64 `json:"latitude"`
}
type MerchantDetailParam struct  {
	Name string `json:"name"`
	OpenId string `json:"open_id"`
	Json string `json:"json"`
	Id int64 `json:"id"`
	//经度
	Longitude float64 `json:"longitude"`
	//纬度
	Latitude float64 `json:"latitude"`
	//覆盖距离 (单位米)
	CoverDistance float64 `json:"cover_distance"`
	AppId string `json:"app_id"`
	//商户图片
	Imgs []MerchantImgDto
}

//商户主要图片DTO
type MerchantImgDto struct  {
	Id int64 `json:"id"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Url string `json:"url"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}


type MerchantDetailDto struct  {
	Name string `json:"name"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Status int `json:"status"`
	Json string `json:"json"`
	//商户地址
	Address string `json:"address"`
	//权重
	Weight int `json:"weight"`
	//距离(单位米)
	Distance float64 `json:"distance"`

}

func MerchantUpdate(c *gin.Context)  {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	appId,err :=security.GetAppId(c.Request)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}


	var param MerchantDetailParam
	err =c.BindJSON(&param)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	if param.Id==0 {
		util.ResponseError400(c.Writer,"商户ID不能为空")
		return
	}
	param.OpenId = openId
	param.AppId =appId

	err =service.MerchantUpdate(merchantDetailParamToDll(param))
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

func MerchantAdd(c *gin.Context)  {

	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	appId,err :=security.GetAppId(c.Request)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	var param MerchantDetailParam
	err =c.BindJSON(&param)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	param.OpenId = openId
	if param.Name=="" {
		util.ResponseError400(c.Writer,"名称不能为空!")
		return
	}

	param.AppId =appId

	mdll,err :=service.MerchantAdd(merchantDetailParamToDll(param))
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	param.Id =mdll.Id
	c.JSON(http.StatusOK,param)

}

func MerchantWithOpenId(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	appId :=security.GetAppId2(c.Request)
	openId := c.Param("open_id")

	merchant,err := service.MerchantWithOpenId(openId,appId)

	if merchant!=nil{
		c.JSON(http.StatusOK,merchantToDto(merchant))
		return
	}

	util.ResponseError400(c.Writer,"没有找到信息!")


}

//附近商户
func MerchatNear(c *gin.Context)  {
	appId,err := CheckAppAuth(c)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	longitude :=c.Query("longitude")
	latitude :=c.Query("latitude")

	flongitude,_ :=strconv.ParseFloat(longitude,20)
	flatitude,_ :=strconv.ParseFloat(latitude,20)
	mDetailList,err := service.MerchantNear(flongitude,flatitude,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	mDetailDtos :=make([]*MerchantDetailDto,0)
	if mDetailList!=nil {

		for _,mDetail :=range mDetailList {
			mDetailDtos = append(mDetailDtos,merchantDetailToDto(mDetail))
		}
	}

	c.JSON(http.StatusOK,mDetailDtos)
}

//根据图片标记查询商户图片
func MerchantImgWithFlag(c *gin.Context)  {

	appId,err := CheckAppAuth(c)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}

	flags :=c.Query("flags")
	mopenId := c.Param("open_id")

	merchantImgs,err := service.MerchantImgWithFlag(flags,mopenId,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	imgddtos := make([]*MerchantImgDto,0)
	if merchantImgs!=nil {
		for _,mImgsModel :=range merchantImgs{
			imgddtos = append(imgddtos,merchantImgToDto(mImgsModel))
		}
	}
	c.JSON(http.StatusOK,imgddtos)

}

func merchantImgToDto(model *dao.MerchantImgs) *MerchantImgDto   {
	dto := &MerchantImgDto{}
	dto.MerchantId = model.MerchantId
	dto.Json = model.Json
	dto.AppId = model.AppId
	dto.Flag = model.Flag
	dto.Url = model.Url
	dto.Id = model.Id

	return dto
}

func merchantDetailToDto(model *dao.MerchantDetail) *MerchantDetailDto  {
	dto :=&MerchantDetailDto{}
	dto.AppId=model.AppId
	dto.Distance = model.Distance
	dto.Json = model.Json
	dto.Name = model.Name
	dto.OpenId = model.OpenId
	dto.Status = model.Status
	dto.Address = model.Address
	dto.Weight = model.Weight


	return dto

}

func merchantDetailParamToDll(param MerchantDetailParam)  *service.MerchantDetailDLL {

	dll :=&service.MerchantDetailDLL{}
	dll.OpenId = param.OpenId
	dll.Name = param.Name
	dll.Json = param.Json
	dll.AppId = param.AppId
	dll.CoverDistance = param.CoverDistance
	dll.Latitude = param.Latitude
	dll.Longitude = param.Longitude
	dll.Id = param.Id

	if param.Imgs!=nil {
		imgdlls := make([]service.MerchantImgDLL,0)
		for _,imgDto :=range param.Imgs  {
			imgdlls = append(imgdlls,merchantImgToDLL(imgDto))
		}

		dll.Imgs=imgdlls
	}

	return dll
}

func merchantImgToDLL(model MerchantImgDto) service.MerchantImgDLL  {
	dll :=service.MerchantImgDLL{}
	dll.Url = model.Url
	dll.Flag = model.Flag
	dll.AppId = model.AppId
	dll.Json = model.Json
	dll.MerchantId = model.MerchantId

	return dll
}

func merchantToDto(model *dao.Merchant)  *MerchantDto {
	dto:=&MerchantDto{}
	dto.Json = model.Json
	dto.Address = model.Address
	dto.AppId = model.AppId
	dto.CoverDistance = model.CoverDistance
	dto.Latitude=model.Latitude
	dto.Longitude = model.Longitude
	dto.Name = model.Name
	dto.Id = model.Id
	dto.OpenId = model.OpenId

	return dto
}