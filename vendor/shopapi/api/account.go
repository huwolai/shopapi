package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
	"shopapi/service"
)

type AccountPreRechargeDto struct  {
	OpenId string `json:"open_id"`
	Money float64 `json:"money"`
	PayType int `json:"pay_type"`
}

//账户充值
func AccountPreRecharge(c *gin.Context)  {
	_,err :=CheckAppAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	openId,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	var param AccountPreRechargeDto
	err =c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if openId!=c.Param("open_id") {
		util.ResponseError400(c.Writer,"不能跟别人充值!")
		return
	}
	param.OpenId = c.Param("open_id")


	model :=&service.AccountRechargeModel{}
	model.Money = param.Money
	model.OpenId = param.OpenId
	model.PayType = param.PayType
	resultMap,err := service.AccountPreRecharge(model)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,resultMap)
}
