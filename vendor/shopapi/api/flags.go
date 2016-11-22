package api

import (
	"github.com/gin-gonic/gin"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
	"strings"
)

type Flags struct  {
	AppId string `json:"app_id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
	Type string `json:"type"`
	//Json string `json:"json"`
}

func FlagsWithTypes(c *gin.Context)  {

	stypes := c.Query("types")
	status :=c.Query("status")
	appId :=security.GetAppId2(c.Request)
	var typesArray []string
	var statusArray []string
	if stypes!="" {
		typesArray = strings.Split(stypes,",")
	}
	if status!="" {
		statusArray = strings.Split(status,",")
	}
	flags,err := service.FlagsWithTypes(typesArray,statusArray,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}
	list :=make([]*Flags,0)
	if flags!=nil{
		for _,flag :=range flags {
			list = append(list,flagsToS(flag))
		}
	}
	c.JSON(http.StatusOK,list)

}

func flagsToS(model *dao.Flags) *Flags  {
	s :=&Flags{}
	s.AppId = model.AppId
	s.Flag = model.Flag
	s.Name = model.Name
	s.Type = model.Type
	//s.Json = model.Json

	return s
}
func FlagsGetJsonWithTypes(c *gin.Context)  {
	stypes := c.Query("types")
	status :=c.Query("status")
	appId :=security.GetAppId2(c.Request)
	var typesArray []string
	var statusArray []string
	if stypes!="" {
		typesArray = strings.Split(stypes,",")
	}
	if status!="" {
		statusArray = strings.Split(status,",")
	}
	flags,err := service.FlagsGetJsonWithTypes(typesArray,statusArray,appId)
	if err!=nil{
		log.Error(err)		
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}
	c.JSON(http.StatusOK,flags)
}
func FlagsSetJsonWithTypes(c *gin.Context)  {
	type Param struct {
		Json	string `json:"json"`
	}

	var json Param
	c.BindJSON(&json)
	//json.Json= c.PostForm("json")
	//json.Json= "test"
	
	types:=c.Param("type")
	
	appId :=security.GetAppId2(c.Request)
	
	err:=service.FlagsSetJsonWithTypes(types,json.Json,appId)	
	if err!=nil {
		util.ResponseError400(c.Writer,"操作失败!")
		return
	}
	util.ResponseSuccess(c.Writer)
}












