package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"net/http"
	"shopapi/dao"
	"strconv"
)

type MerchantAddParam struct  {
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

func MerchantAdd(c *gin.Context)  {

	appId,err := CheckAppAuth(c)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	var param MerchantAddParam
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
	param.AppId = appId

	mdll,err :=service.MerchantAdd(merchantAddParamToDll(param))
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	param.Id =mdll.Id
	c.JSON(http.StatusOK,param)

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

func merchantAddParamToDll(param MerchantAddParam)  *service.MerchantAddDLL {

	dll :=&service.MerchantAddDLL{}
	dll.OpenId = param.OpenId
	dll.Name = param.Name
	dll.Json = param.Json
	dll.AppId = param.AppId
	return dll
}