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
	"shopapi/setting"
	"gitlab.qiyunxin.com/tangtao/utils/page"
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
)

type AccountPreRechargeDto struct  {
	OpenId string `json:"open_id"`
	Money float64 `json:"money"`
	PayType int `json:"pay_type"`
	Remark string `json:"content"`
	Opt string `json:"opt"`
}

type AccountDetailDto struct  {
	//账户余额 单位元
	Amount float64 `json:"amount"`
	CashoutAmount float64 `json:"cashout_amount"`
	FreezeAmount float64 `json:"freeze_amount"`
	//账户状态 1.正常 0.异常 3.锁定
	Status int `json:"status"`
	//是否设置支付密码
	PasswordIsSet int `json:"password_is_set"`
}

type Account struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Mobile string `json:"mobile"`
	Money float64 `json:"money"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	Status int `json:"status"`
	YdgyId string `json:"ydgy_id"`
	YdgyName string `json:"ydgy_name"`
	YdgyStatus int64 `json:"ydgy_status"`
	Name string `json:"username"`
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

	code :=GetRandCode()
	redis.SetAndExpire(comm.CODE_PAYPWD_PREFIX+mobile,code,comm.CODE_PAYPWD_EXPIRE)
	configMap :=setting.GetYunTongXunSetting()
	err =service.SendSMSOfYunTongXun(mobile,configMap["paypwdcode_template_id"],[]string{code})
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"短信发送失败!")
		return
	}

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
	appId := security.GetAppId2(c.Request)

	model :=&service.AccountRechargeModel{}
	model.Money = param.Money
	model.OpenId = param.OpenId
	model.PayType = param.PayType
	model.AppId = appId
	resultMap,err := service.AccountPreRecharge(model,2,"")
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

func AccountsGet(c *gin.Context)  {
	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	appId :=security.GetAppId2(c.Request)
	mobile := c.Query("mobile")
	
	userName 	:= c.Query("username") // 用户名
	ydgyId 	 	:= c.Query("ydgy_id")	//一点公益 ID
	ydgyName 	:= c.Query("ydgy_name") //一点公益名字
	ydgyStatus  := c.Query("ydgy_status") //一点公益状态
	
	
	accounts,err := service.AccountsWith(pIndex,pSize,mobile,appId,userName,ydgyId,ydgyName,ydgyStatus)
	if err!=nil{
		util.ResponseError400(c.Writer,"查询失败！")
		return
	}
	
	results :=make([]*Account,0)
	if accounts!=nil{
		var detailModel *service.AccountDetailModel
		log.Info(detailModel)
		for _,account :=range accounts  {			
			detailModel,_ =service.AccountDetail(account.OpenId)
			account.Money=float64(detailModel.Amount)/100.0
			//account.Money=0
			results = append(results,accountToA(account))
		}
	}

	total,err :=service.AccountsWithCount(mobile,appId,userName,ydgyId,ydgyName,ydgyStatus)
	if err!=nil{
		util.ResponseError400(c.Writer,"查询总数量失败！")
		return
	}

	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(total),results))
}

func accountToA(model *dao.Account) *Account  {
	a := &Account{}
	a.AppId = model.AppId
	a.Id = model.Id
	a.Mobile = model.Mobile
	a.Money = model.Money
	a.Status = model.Status
	a.OpenId = model.OpenId
	a.CreateTime = qtime.ToyyyyMMddHHmm(model.CreateTime)
	a.UpdateTime = qtime.ToyyyyMMddHHmm(model.UpdateTime)
	a.YdgyId = model.YdgyId
	a.YdgyName = model.YdgyName
	a.YdgyStatus = model.YdgyStatus
	a.Name = model.Name

	return a
}

func accountDetailModelToDto(model *service.AccountDetailModel) *AccountDetailDto {

	dto :=&AccountDetailDto{}
	dto.Amount = float64(model.Amount)/100.0
	dto.PasswordIsSet = model.PasswordIsSet
	dto.Status = model.Status
	dto.FreezeAmount = float64(model.FreezeAmount/100.0)
	dto.CashoutAmount = float64(model.CashoutAmount/100.0)
	return dto
}
//配置登入界面
func GetOnKey(c *gin.Context)  {
	GetOnKey,err :=service.GetOnKey()
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	param :=map[string]interface{}{
		"res":GetOnKey.Status,
	}
	
	c.JSON(http.StatusOK,param)
}

//管理员后台-账户充值
func AccountPreRechargeByAdmin(c *gin.Context)  {
	openId:=c.Param("open_id")
	if openId!="wesdfsfsdf23323" {
		util.ResponseError400(c.Writer,"权限不足！")
		return
	}
	password := c.Query("sepassword")
	if password!="180181" {
		util.ResponseError400(c.Writer,"权限不足.！")
		return
	}	
	
	var param AccountPreRechargeDto
	err :=c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	if len(param.Remark)<1 {
		util.ResponseError400(c.Writer,"充值或扣款缘由未填！")
		return
	}
	
	if param.Money==0 {
		util.ResponseError400(c.Writer,"金额错误！")
		return
	}
	
	//log.Info(param)
	//param.OpenId = c.Param("open_id")
		
	appId := security.GetAppId2(c.Request)

	model :=&service.AccountRechargeModel{}	
	model.OpenId = param.OpenId
	model.PayType = 3 //param.PayType
	model.AppId = appId
	model.Remark = param.Remark
	
	var resultMap map[string]interface{}
	if param.Money<0 {
		model.Money = 0-param.Money
		resultMap,err = service.AccountPreRechargeMinus(model,1,param.Opt)
		if err!=nil {
			log.Error(err)
			util.ResponseError400(c.Writer,err.Error())
			return
		}		
	}else{
		model.Money = param.Money
		resultMap,err = service.AccountPreRecharge(model,1,param.Opt)
		if err!=nil {
			log.Error(err)
			util.ResponseError400(c.Writer,err.Error())
			return
		}
	}
	c.JSON(http.StatusOK,resultMap)
}
//账户充值记录  后台
func RechargeRecordByAdmin(c *gin.Context)  {
	appId := security.GetAppId2(c.Request)
	openId:=c.Param("open_id")
	froms,err :=strconv.ParseInt(c.Query("froms"),10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"充值来源有误")
		return
	}
	
	rechargeRecord,err:=service.RechargeRecordByAdmin(openId,appId,froms)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	dto 	:=make([]service.AccountRecharge,0)
	for _,item :=range rechargeRecord  {
		dto = append(dto,service.RechargeRecordFormat(item))
	}
	
	c.JSON(http.StatusOK,dto)
}
//账户充值记录  后台 列表
func RechargeRecordByAdmins(c *gin.Context)  {
	appId := security.GetAppId2(c.Request)
	froms,err :=strconv.ParseInt(c.Query("froms"),10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"充值来源有误")
		return
	}
	
	pIndex,pSize := page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	
	var search dao.AccountRechargeSearch
	search.No		=c.Query("no")
	search.YdgyId	=c.Query("ydgy_id")
	search.YdgyName	=c.Query("ydgy_name")
	search.Mobile	=c.Query("mobile")
	
	rechargeRecord,total,err:=service.RechargeRecordByAdmins(appId,froms,pIndex,pSize,search)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	dto 	:=make([]service.AccountRecharge,0)
	var dtoItem	service.AccountRecharge
	account := dao.NewAccount()
	for _,item :=range rechargeRecord  {	
		dtoItem=service.RechargeRecordFormat(item)

		if len(dtoItem.OpenId)>0 {
			account,_ =account.AccountWithOpenId(dtoItem.OpenId,appId)		
			dtoItem.Mobile=account.Mobile
			dtoItem.YdgyId=account.YdgyId
			dtoItem.YdgyName=account.YdgyName
		}
		dto = append(dto,dtoItem)
	}

	//c.JSON(http.StatusOK,dto)
	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(total),dto))
}







































































