package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"net/http"
)
//应用更新日志
func AppUpdateLog(c *gin.Context)  {
	appLog,err :=service.AppUpdateLog()
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}	
	c.JSON(http.StatusOK,appLog)
}