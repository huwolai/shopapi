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

	return s
}