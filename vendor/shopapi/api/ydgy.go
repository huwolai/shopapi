package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
	"strconv"
	"shopapi/service"
)
//一点公益ID号绑定
func YdgySetId(c *gin.Context)  {
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	ydgyId,_ := strconv.ParseInt(c.Param("id"),10,64)	
	if ydgyId<1 {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	
	err :=service.YdgySetId(openId,ydgyId)
	if err!=nil{
		util.ResponseError400(c.Writer,"设置失败!")
		return
	}
	
	util.ResponseSuccess(c.Writer)
	
}
//一点公益ID号获取
func YdgyGetId(c *gin.Context)  {
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}

	ydgyId,err :=service.YdgyGetId(openId)	
	if err!=nil{
		util.ResponseError400(c.Writer,"获取失败!")
		return
	}
	
	c.JSON(http.StatusOK,gin.H{
		"id":ydgyId,
	})	
}