package service

import (
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"net/http"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/dao"
	"shopapi/comm"
	"shopapi/redis"
)

type AccountRechargeModel struct  {
	AppId string
	//充值账户
	OpenId string
	//充值金额
	Money float64
	PayType int  //付款类型(1.支付宝 2.微信)
	Remark string
}

type AccountDetailModel struct  {

	//账户余额 单位分
	Amount int64 `json:"amount"`
	CashoutAmount int64 `json:"cashout_amount"`
	//冻结金额
	FreezeAmount int64 `json:"freeze_amount"`
	//账户状态 1.正常 0.异常 3.锁定
	Status int `json:"status"`
	CreateTime string `json:"create_time"`
	//是否设置支付密码
	PasswordIsSet int `json:"password_is_set"`

}

//账户预充值
func AccountPreRecharge(model *AccountRechargeModel) (map[string]interface{},error) {

	accountRecharge :=dao.NewAccountRecharge()
	accountRecharge.Amount = model.Money
	accountRecharge.No = util.GenerUUId()
	accountRecharge.AppId = model.AppId
	accountRecharge.Status = comm.ACCOUNT_RECHARGE_STATUS_WAIT
	err :=accountRecharge.Insert()
	if err!=nil{
		log.Error(err)
		return nil,errors.New("充值记录插入失败!")
	}
	
	remark:="充值"
	if len(model.Remark)>0{
		remark=model.Remark
	}

	//参数
	params := map[string]interface{}{
		"out_trade_no": accountRecharge.No,
		"out_trade_type": comm.Trade_Type_Recharge,
		"open_id": model.OpenId,
		"amount": int64(model.Money*100),
		"trade_type": comm.Trade_Type_Recharge,
		"pay_type": model.PayType,
		"title": remark,
		"client_ip": "127.0.0.1",
		"notify_url": config.GetValue("notify_url").ToString(),
		"remark": remark,
	}
	resultMap,err :=RequestPayApi("/pay/makeprepay",params)

	return resultMap,err
}



func AccountDetail(openId string) (*AccountDetailModel,error)  {
	//参数
	params := map[string]interface{}{
		"open_id": openId,
	}
	log.Info(params)

	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetPayapiSign(params)
	log.Info(fmt.Sprintf("%s.%s",basesign,sign))
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(params);

	response,err := network.Post(config.GetValue("payapi_url").ToString()+"/account/detail",paramData,headers)
	if err!=nil{
		return nil,err
	}
	if response.StatusCode==http.StatusOK {
		var resultModel *AccountDetailModel
		err =util.ReadJsonByByte([]byte(response.Body),&resultModel)
		if err!=nil{
			return nil,err
		}

		return resultModel,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}

	return nil,errors.New("充值失败")
}

//通过短信登录
func LoginForSMS(mobile string,code string,appId string) (map[string]interface{},error)  {

	commusermap,err  :=loginSMSOfCommuser(mobile,code,appId)
	if err!=nil {

		return nil,err
	}
	openId := commusermap["open_id"].(string)

	account := dao.NewAccount()
	account,err =account.AccountWithOpenId(openId,appId)
	if err!=nil {
		return nil,err
	}
	if account==nil {
		//password :=util.GenerUUId()
		account = dao.NewAccount()
		account.AppId = appId
		account.Money = 0
		account.Mobile = mobile
		account.OpenId = openId
		//account.Password = password
		account.Status =comm.ACCOUNT_STATUS_WAIT_BINDPAY //等待开通支付
		err =account.Insert()
		if err!=nil {
			log.Error(err)
			return nil,err
		}
	}

	if account.Status==comm.ACCOUNT_STATUS_WAIT_BINDPAY {
		//绑定支付
		err :=payBind(openId,account.Password)
		if err!=nil {
			return nil,err
		}

		err = account.AccountUpdateStatus(comm.ACCOUNT_STATUS_NORMAL,openId,appId)
		if err!=nil {
			return nil,err
		}
	}

	merchant := dao.NewMerchant()
	merchant,err =merchant.MerchantWithOpenId(openId,appId)
	if err!=nil{
		return nil,err
	}

	commusermap["is_merchant"] = 0
	if merchant!=nil&&merchant.Status==1 {
		commusermap["is_merchant"]=1
		commusermap["merchant_id"] = merchant.Id
	}

	return commusermap,nil

}

//修改支付密码
func PayPwdUpdate(openId string,mobile string,newpwd string,code string,appId string) error  {

	paycode :=redis.GetString(comm.CODE_PAYPWD_PREFIX+mobile)
	if paycode==""{

		return errors.New("请先获取验证码!")
	}

	if paycode!=code {
		return errors.New("验证码输入不对!")
	}

	account :=dao.NewAccount()
	account,err :=account.AccountWithMobile(mobile,appId)
	if err!=nil{
		return err
	}

	params :=map[string]interface{}{
		"newpwd":newpwd,
		"oldpwd": account.Password,
		"open_id": openId,
	}
	_,err = RequestPayApi("/pay/password",params)
	if err!=nil{
		return err
	}
	err =account.AccountUpdatePwd(newpwd,openId,appId)

	return err
}

//绑定支付
func payBind(openId string,password string) error  {
	header := GetPayAuthHeader(openId)
	params :=map[string]interface{}{
		"open_id":openId,
	}
	paramData,_:= json.Marshal(params);
	response,err :=network.Post(config.GetValue("payapi_url").ToString()+"/pay/bind",paramData,header)
	if err!=nil {
		return err
	}

	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		if err!=nil{
			return err
		}

		return nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return errors.New(resultMap["err_msg"].(string))
	}

	return errors.New("开通支付失败!")
}

func loginSMSOfCommuser(mobile,code,appId string) ( map[string]interface{},error)  {
	param := map[string]interface{}{
		"mobile":mobile,
		"code": code,
	}
	paramData,_:= json.Marshal(param);

	header :=map[string]string{
		"app_id":appId,
	}
	response,err :=network.Post(config.GetValue("commuser_url").ToString()+"/loginSMS",paramData,header)
	if err!=nil {
		return nil,err
	}

	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		if err!=nil{
			return nil,err
		}

		return resultMap,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}

	return nil,errors.New("调用统一用户中心登录失败!")
}

func AccountsWith(pageIndex uint64,pageSize uint64,mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string) ([]*dao.Account,error)   {

	accounts,err := dao.NewAccount().AccountsWith(pageIndex,pageSize,mobile,appId,userName,ydgyId,ydgyName,ydgyStatus)
	
	if err==nil {
		merchant := dao.NewMerchant()
		for i,account :=range accounts  {
			merchant,_ = merchant.MerchantWithOpenId(account.OpenId,appId)
			if merchant!=nil {
				accounts[i].Name=merchant.Name
			}
		}
	}	
	

	return accounts,err
}

func AccountsWithCount(mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string) (int64,error) {

	return dao.NewAccount().AccountsWithCount(mobile,appId,userName,ydgyId,ydgyName,ydgyStatus)
}
//配置登入界面
func GetOnKey() (*dao.GetOnKey,error)   {
	GetOnKey :=dao.NewGetOnKey()
	GetOnKey,err :=GetOnKey.GetOnKey()
	return GetOnKey,err
}
//账户预充值 减去 AccountPreRechargeMinus
func AccountPreRechargeMinus(model *AccountRechargeModel) (map[string]interface{},error) {
	//制作预付款
	payparams := map[string]interface{}{
		"open_id":model.OpenId,
		"type": 1,
		"amount": int64(model.Money*100),
		"title": model.Remark,
		"remark": model.Remark,
	}
	resultImprestMap,err := RequestPayApi("/pay/makeimprest",payparams)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	log.Info("后台领取预付款")
	//领取预付款
	payparams =map[string]interface{}{
		"open_id": model.OpenId,
		"code": resultImprestMap["code"],
		"amount": int64(model.Money*100),
		"title": model.Remark,
		"remark": model.Remark,
	}
	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetPayapiSign(payparams)
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(payparams);
	response,err := network.Post(config.GetValue("payapi_url").ToString()+"/imprests/fetch",paramData,headers)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	
	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		return resultMap,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		log.Error(err)
		return nil,errors.New(resultMap["err_msg"].(string))
	}else{
		return nil,errors.New("请求API失败!")
	}
}




























