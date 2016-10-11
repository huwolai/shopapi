package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
)
//商品初始化售出数量
func ProductInitNum(c *gin.Context)  {
	err :=service.ProductInitNum()
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}	
	util.ResponseSuccess(c.Writer)
}