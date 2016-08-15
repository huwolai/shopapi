package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"github.com/Azure/azure-sdk-for-go/core/http"
)

type MerchantAddParam struct  {
	Name string `json:"name"`
	OpenId string `json:"open_id"`
	Json string `json:"json"`
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
}

func MerchantAdd(c *gin.Context)  {

	appId,err := CheckAuth(c)
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

func merchantAddParamToDll(param MerchantAddParam)  *service.MerchantAddDLL {

	dll :=&service.MerchantAddDLL{}
	dll.OpenId = param.OpenId
	dll.Name = param.Name
	dll.Json = param.Json
	dll.AppId = param.AppId
	return dll
}