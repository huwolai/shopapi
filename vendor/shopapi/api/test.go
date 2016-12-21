package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
)

func Test(c *gin.Context) {
	appId	:="shopapi"
	openId	:="2bd209e36084479cbbb7258f12fce02f"
	title	:="@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事"
	content	:="@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事@订单已付款事"
	types	:="chefOrder"

	err:=service.PushSingle(openId,appId,title,content,types)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
	
	
	
	
	
}

















