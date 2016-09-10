package api

import (
	"github.com/gin-gonic/gin"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
)

type Flags struct  {
	AppId string
	Name string
	Flag string
	Type string
}

func FlagsWithTypes(c *gin.Context)  {

	stypes := c.Query("types")

	appId :=security.GetAppId2(c.Request)
	flags,err := service.FlagsWithTypes(stypes,appId)
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