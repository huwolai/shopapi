package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"net/http"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)



func AddressWithRecom(c *gin.Context) {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return 
	}
	appId := security.GetAppId2(c.Request)
	address,err := service.AddressWithRecom(openId,appId)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}
	if address==nil {
		util.ResponseError400(c.Writer,"请先完善用餐地址!")
		return
	}

	c.JSON(http.StatusOK,service.AddressToDto(address))
}

func AddressWithId(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return
	}

	sid :=c.Param("id")
	id,err := strconv.ParseInt(sid,10,64)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"参数有误!")
		return
	}
	address,err := service.AddressWithId(id)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}

	if address==nil{
		util.ResponseError400(c.Writer,"地址未找到!")
		return
	}

	c.JSON(http.StatusOK,service.AddressToDto(address))
}

func AddressWithOpenId(c *gin.Context)  {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return
	}
	appId := security.GetAppId2(c.Request)

	addresses,err := service.AddressWithOpenId(openId,appId)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}
	addressdtos :=make([]*service.AddressDto,0)
	if addresses!=nil {


		for _,address :=range addresses {
			addressdtos = append(addressdtos,service.AddressToDto(address))
		}
	}
	c.JSON(http.StatusOK,addressdtos)
}

func AddressDelete(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return
	}

	 sid :=c.Param("id")
	id,err := strconv.ParseInt(sid,10,64)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"参数有误!")
		return
	}
	 err =service.AddressDelete(id)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"删除失败!")
		return
	}

	util.ResponseSuccess(c.Writer)

}

func AddressUpdate(c *gin.Context)  {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return
	}

	appId := security.GetAppId2(c.Request)

	var param *service.AddressDto
	err =c.BindJSON(&param)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,"参数有误!")
		return
	}
	param.OpenId = openId
	param.AppId = appId

	dto,err := service.AddressUpdate(param)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}

	c.JSON(http.StatusOK,dto)
}

func AddressAdd(c *gin.Context)  {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证有误!")
		return
	}
	appId := security.GetAppId2(c.Request)

	var param *service.AddressDto
	err =c.BindJSON(&param)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,"参数有误!")
		return
	}

	param.OpenId = openId
	param.AppId = appId

	dto,err :=service.AddressAdd(param)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}

	c.JSON(http.StatusOK,dto)
}

