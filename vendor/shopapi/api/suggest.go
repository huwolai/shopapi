package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
)

type Suggest struct  {
	OpenId string `json:"open_id"`
	Contact string `json:"contact"`
	Content string `json:"content"`
}

// 添加建议
func SuggestAdd(c *gin.Context)  {
	var param *Suggest
	err :=c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	err =service.SuggestAdd(param.Content,param.Contact,param.OpenId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加失败!")
		return
	}
	util.ResponseSuccess(c.Writer)
}
