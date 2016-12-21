package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
	"strconv"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/security"
)
type Ydgy struct {
	Id		string   `json:"id"`
	Name	string  `json:"name"`
	Mine	int64   `json:"mine"`
}
//一点公益ID号绑定
func YdgySetId(c *gin.Context)  {
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	
	var ydgy Ydgy
	c.BindJSON(&ydgy)
	/* ydgy.Id		= ydgyId
	ydgy.Name	= c.PostForm("name")
	ydgy.Mine	= 1 */
	
	err :=service.YdgySetId(openId,ydgy.Id,ydgy.Name,ydgy.Mine)
	if err!=nil{
		util.ResponseError400(c.Writer,"设置失败!")
		log.Info("设置一点公益失败!")
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

	ydgy,err :=service.YdgyGetId(openId)	
	if err!=nil{
		util.ResponseError400(c.Writer,"获取失败!")
		return
	}
	
	c.JSON(http.StatusOK,ydgy)
}
//一点公益ID号状态审核
func YdgySetIdWithStatus(c *gin.Context)  {
	appId 	:= security.GetAppId2(c.Request)
	openId  := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	
	YdgyStatus,_ := strconv.ParseInt(c.Param("status"),10,64)	
	if YdgyStatus<1 {
		util.ResponseError400(c.Writer,"状态不能为空!")
		return
	}	
	
	err :=service.YdgySetIdWithStatus(openId,YdgyStatus,c.Query("failres"))	
	if err!=nil{
		util.ResponseError400(c.Writer,"操作失败!")
		return
	}
	
	//推送
	if len(appId)>0 {
		if len(c.Query("failres"))>0 {
			service.PushSingle(openId,appId,"审核失败",c.Query("failres"),"userYYGY")
		}else{
			service.PushSingle(openId,appId,"审核成功","审核成功","userYYGY")
		}
	}
	
	
	util.ResponseSuccess(c.Writer)
}
//一点公益ID号删除
func YdgySetIdWithDelete(c *gin.Context)  {
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}

	err :=service.YdgySetIdWithDelete(openId)	
	if err!=nil{
		util.ResponseError400(c.Writer,"删除失败!")
		return
	}	
	util.ResponseSuccess(c.Writer)
}






















