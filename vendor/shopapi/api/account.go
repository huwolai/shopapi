package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/redis"
	"shopapi/comm"
	"strconv"
	"math/rand"
	"time"
)

type AccountPreRechargeDto struct  {
	OpenId string `json:"open_id"`
	Money float64 `json:"money"`
	PayType int `json:"pay_type"`
}

type AccountDetailDto struct  {
	//账户余额 单位元
	Amount float64 `json:"amount"`
	//账户状态 1.正常 0.异常 3.锁定
	Status int `json:"status"`
	//是否设置支付密码
	PasswordIsSet int `json:"password_is_set"`
}

type LoginForSMSParam struct  {
	//手机号
	Mobile string `json:"mobile"`
	//验证码
	Code string `json:"code"`
}

type PayPwdUpdateDto struct  {
	Mobile string `json:"mobile"`
	Code string `json:"code"`
	Password string `json:"password"`

}


func LoginForSMS(c *gin.Context)  {

	var loginSms LoginForSMSParam
	err :=c.BindJSON(&loginSms)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	appId := security.GetAppId2(c.Request)

	resultMap,err :=service.LoginForSMS(loginSms.Mobile,loginSms.Code,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,resultMap)
}

//支付密码修改短信
func PayPwdUpdateSMS(c *gin.Context)  {
	//获取用户openid
	_,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	mobile := c.Param("mobile")
	if mobile=="" {

		util.ResponseError400(c.Writer,"请输入手机号!")
		return
	}

	if len(mobile)!=11 {
		util.ResponseError400(c.Writer,"手机号输入有误!")
		return
	}

	code :=redis.GetString(comm.CODE_PAYPWD_PREFIX+mobile)
	if code== ""{
		//code =GetRandCode()
		code = "1111"
	}
	redis.SetAndExpire(comm.CODE_PAYPWD_PREFIX+mobile,code,comm.CODE_PAYPWD_EXPIRE)

	//err =service.SendCodeSMS(mobile,code)
	//if err!=nil{
	//	log.Error(err)
	//	util.ResponseError400(c.Writer,"短信发送失败!")
	//	return
	//}

	util.ResponseSuccess(c.Writer)
}

func PayPwdUpdate(c *gin.Context)  {
	openId,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	var param PayPwdUpdateDto
	err =c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	mobile :=c.Param("mobile")
	param.Mobile = mobile
	appId := security.GetAppId2(c.Request)

	err = service.PayPwdUpdate(openId,mobile,param.Password,param.Code,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	util.ResponseSuccess(c.Writer)

}

func GetRandCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var code string
	for i:=0; i<4; i++ {
		code+=strconv.Itoa(r.Intn(9))
	}

	return code
}


//账户充值
func AccountPreRecharge(c *gin.Context)  {

	//获取用户openid
	openId,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	var param AccountPreRechargeDto
	err =c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if openId!=c.Param("open_id") {
		util.ResponseError400(c.Writer,"不能跟别人充值!")
		return
	}
	param.OpenId = c.Param("open_id")


	model :=&service.AccountRechargeModel{}
	model.Money = param.Money
	model.OpenId = param.OpenId
	model.PayType = param.PayType
	resultMap,err := service.AccountPreRecharge(model)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,resultMap)
}

func AccountDetail(c *gin.Context)  {

	//获取用户openid
	openId,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	detailModel,err :=service.AccountDetail(openId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,accountDetailModelToDto(detailModel))
}

func accountDetailModelToDto(model *service.AccountDetailModel) *AccountDetailDto {

	dto :=&AccountDetailDto{}
	dto.Amount = float64(model.Amount)/100.0
	dto.PasswordIsSet = model.PasswordIsSet
	dto.Status = model.Status
	return dto
}