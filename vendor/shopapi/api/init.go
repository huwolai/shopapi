package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/dao"
)
//判断token是否过期
func TokenWithExpired(c *gin.Context)  {
	_,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError400(c.Writer,"重新登入!")
		return
	}
	util.ResponseSuccess(c.Writer)
}
//厨师面试登记表
func MerchantResumesWithAdd(c *gin.Context)  {
	var resumes dao.MerchantResume
	c.BindJSON(&resumes)
	//resumes.Name= c.PostForm("name")
	
	err:=service.MerchantResumesWithAdd(resumes)	
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}