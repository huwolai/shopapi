package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"net/http"
	"strconv"
	"shopapi/dao"
)

type UserBank struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	AccountName string `json:"account_name"`
	BankName string `json:"bank_name"`
	BankNo string `json:"bank_no"`
	BankFullName string `json:"branch"`
}

func UserBankGet(c *gin.Context)  {
	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	list,err := service.UserBankGet(openId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}
	dtos :=make([]*UserBank,0)
	if list!=nil{
		for _,uBank :=range list {
			dtos = append(dtos,userBankToA(uBank))
		}
	}

	c.JSON(http.StatusOK,dtos)

}

func UserBankAdd(c *gin.Context)   {
	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	var uBank *UserBank
	err :=c.BindJSON(&uBank)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误!")
		return
	}
	uBank.OpenId = openId
	uBank.AppId = appId

	ub,err :=service.UserBankAdd(userBankToS(uBank))
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加失败!")
		return
	}
	
	c.JSON(http.StatusOK,userBankToA(ub))
	
}

func UserBankDel(c *gin.Context)  {
	id := c.Param("id")
	iid,err :=strconv.ParseInt(id,10,64)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"id格式有误!")
		return
	}
	err =service.UserBankDel(iid)
	if err!=nil{
		util.ResponseError400(c.Writer,"删除失败!")
		return
	}

	util.ResponseSuccess(c.Writer)
}

func UserBankUpdate(c *gin.Context)  {
	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	var uBank *UserBank
	err :=c.BindJSON(&uBank)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误!")
		return
	}
	uBank.OpenId = openId
	uBank.AppId = appId

	ub,err :=service.UserBankUpdate(userBankToS(uBank))
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"修改失败!")
		return
	}

	c.JSON(http.StatusOK,userBankToA(ub))
}

func userBankToA(s *dao.UserBank) *UserBank  {
	a := &UserBank{}
	a.Id = s.Id
	a.OpenId = s.OpenId
	a.AppId = s.AppId
	a.AccountName = s.AccountName
	a.BankName = s.BankName
	a.BankNo = s.BankNo
	a.BankFullName = s.BankFullName
	return a
}

func userBankToS(a *UserBank) *service.UserBank  {

	s :=&service.UserBank{}
	s.Id = a.Id
	s.BankNo = a.BankNo
	s.AccountName = a.AccountName
	s.BankName = a.BankName
	s.AppId = a.AppId
	s.OpenId = a.OpenId
	s.BankFullName = a.BankFullName
	return s

}
