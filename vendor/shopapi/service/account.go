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
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/dao"
	"shopapi/comm"
	"shopapi/redis"
	"time"
	"strconv"
	
)

type AccountRechargeModel struct  {
	AppId 	string
	//充值账户
	OpenId 	string
	//充值金额
	Money 	float64
	PayType int  //付款类型(1.支付宝 2.微信)
	Remark 	string
	Type 	int
	/* 账户余额类型:
1.    系统自动结算的厨师服务费用
2.	系统结算的厨师小店佣金费用
3.	后台给厨师的充值打款(例如食材费和厨师推荐费用等)
4.	后台给厨师充值的通过一点公益的消费金额
每个月15号厨师可提现,  提现金额仅限类型1、2、3属于可提现,    厨师消费时(在app内购买厨师服务和商城商品时)优先使用类型4的金额 */

}

type AccountRecharge struct  {
	Id 			int64 	`json:"id"`
	No 			string	`json:"no"`	
	OpenId 		string	`json:"open_id"`	
	Amount 		float64	`json:"amount"`
	Status 		int		`json:"status"`
	Flag 		string	`json:"flag"`
	Json 		string	`json:"json"`
	Froms 		int		`json:"froms"`
	CreateTime  string  `json:"create_time"`
	Mobile 	    string  `json:"mobile"`
	YdgyId 		string	`json:"ydgy_id"`
	YdgyName 	string	`json:"ydgy_name"`
	Opt 		string	`json:"opt"`
	Remark 		string	`json:"remark"`
	FailRes		string	`json:"fail_res"`
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
	
	FreezeMoney int64 `json:"freeze_money"`
}

//账户预充值
func AccountPreRecharge(model *AccountRechargeModel,from int,opt string) (map[string]interface{},error) {

	accountRecharge :=dao.NewAccountRecharge()
	accountRecharge.Amount = model.Money
	accountRecharge.No = util.GenerUUId()
	accountRecharge.AppId = model.AppId
	accountRecharge.Status = comm.ACCOUNT_RECHARGE_STATUS_WAIT
	accountRecharge.OpenId = model.OpenId
	accountRecharge.Froms = from
	accountRecharge.Opt = opt
	accountRecharge.Remark = model.Remark
	
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
		
		/* account,_ :=dao.NewAccount().AccountWithOpenId(openId,appid)
		if account!=nil {
			resultModel.FreezeMoney=account.FreezeMoney
		} */
		

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

func AccountsWith(pageIndex uint64,pageSize uint64,mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string,openId string,hasMoney string) ([]*dao.Account,error)   {

	accounts,err := dao.NewAccount().AccountsWith(pageIndex,pageSize,mobile,appId,userName,ydgyId,ydgyName,ydgyStatus,openId,hasMoney)
	
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

func AccountsWithCount(mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string,openId string,hasMoney string) (int64,error) {

	return dao.NewAccount().AccountsWithCount(mobile,appId,userName,ydgyId,ydgyName,ydgyStatus,openId,hasMoney)
}
//配置登入界面
func GetOnKey() (*dao.GetOnKey,error)   {
	GetOnKey :=dao.NewGetOnKey()
	GetOnKey,err :=GetOnKey.GetOnKey()
	return GetOnKey,err
}
//账户预充值 减去 AccountPreRechargeMinus
func AccountPreRechargeMinus(model *AccountRechargeModel,from int,opt string) (map[string]interface{},error) {
	
	accountRecharge :=dao.NewAccountRecharge()
	accountRecharge.Amount = 0-model.Money
	accountRecharge.No = util.GenerUUId()
	accountRecharge.AppId = model.AppId
	accountRecharge.Status = comm.ACCOUNT_RECHARGE_STATUS_WAIT
	accountRecharge.OpenId = model.OpenId
	accountRecharge.Froms = from
	accountRecharge.Opt = opt
	accountRecharge.Remark = model.Remark
	
	err :=accountRecharge.Insert()
	if err!=nil{
		return nil,errors.New("充值记录插入失败!")
	}
	
	account := dao.NewAccount()
	account,err =account.AccountWithOpenId(model.OpenId,model.AppId)
	if err!=nil{
		return nil,err
	}
	
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
	//payToken
	payparams = map[string]interface{}{
		"open_id": model.OpenId,
		"password": account.Password,
	}
	resultMap,err := RequestPayApi("/pay/token",payparams)
	if err!=nil{
		return nil,err
	}
	//支付预付款
	payparams = map[string]interface{}{
		"pay_token": resultMap["pay_token"].(string),
		"open_id": model.OpenId,
		"code": resultImprestMap["code"],
		"out_trade_no":"",
	}
	resultMap,err = RequestPayApi("/pay/payimprest",payparams)
	if err!=nil{
		return nil,err
	}
	return resultMap,nil
	/* log.Info(resultMap)
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
	} */
}
//账户充值记录  后台
func RechargeRecordByAdmin(openId string,appId string,froms int64) ([]*dao.AccountRecharge,error) {
	accountRecharge :=dao.NewAccountRecharge()
	return accountRecharge.WithOpenId(openId,appId,froms)
}
func RechargeRecordByAdmins(appId string,froms int64,pageIndex uint64,pageSize uint64,search dao.AccountRechargeSearch ) ([]*dao.AccountRecharge,int64,error) {
	accountRecharge :=dao.NewAccountRecharge()
	
	rechargeRecord,err:=accountRecharge.RecordWithUser(appId,froms,pageIndex,pageSize,search)
	
	count,err:=accountRecharge.RecordWithUserCount(appId,froms,pageIndex,pageSize,search)
	
	//log.Info(err)
	
	return rechargeRecord,count,err
}
//账户充值记录 格式化
func RechargeRecordFormat(model *dao.AccountRecharge) AccountRecharge  {
	//dto 	:=make([]AccountRecharge,0)
	dtoItem		:=AccountRecharge{}
	//for _,item :=range model {
	dtoItem.Id		=model.Id
	dtoItem.No		=model.No
	dtoItem.OpenId	=model.OpenId
	dtoItem.Amount	=model.Amount
	dtoItem.Status	=model.Status
	dtoItem.Flag	=model.Flag
	dtoItem.Json	=model.Json
	dtoItem.Froms	=model.Froms		
	dtoItem.Mobile	=model.Mobile	
	
	dtoItem.Opt		=model.Opt		
	dtoItem.FailRes	=model.FailRes	
	
	dtoItem.Remark	=model.Remark		
	
	dtoItem.CreateTime=time.Unix(model.CreateTimeUnix, 0).Format("2006-01-02 15:04:05")
		
		//dto = append(dto,dtoItem)
	//}
	return dtoItem
}

//账户变动记录
func AccountChangeRecord(model *AccountRechargeModel,from int,opt string) error {
	
	log.Info(opt)
	
	accountRecharge :=dao.NewAccountRecharge()
	accountRecharge.Amount = model.Money
	accountRecharge.No = util.GenerUUId()
	accountRecharge.AppId = model.AppId
	accountRecharge.Status = comm.ACCOUNT_RECHARGE_STATUS_WAIT
	accountRecharge.OpenId = model.OpenId
	accountRecharge.Froms = from
	accountRecharge.Opt = opt
	accountRecharge.Remark = model.Remark
	accountRecharge.Type = model.Type
	
	err :=accountRecharge.Insert()
	if err!=nil{
		return errors.New("充值记录插入失败!")
	}
	return nil
}
//账户变动记录
func AccountChangeRecordOK(model map[string]interface{},appId string) (map[string]interface{},error) {	
	rechargeRecord:=dao.NewAccountRecharge()
	
	record,err:=rechargeRecord.WithNo(model["no"].(string),appId)
	if err!=nil || record==nil{
		return nil,errors.New("无记录!")
	}	
	
	if record.Status>0 {
		return nil,errors.New("已审核!!") 
	}
	
	amount,_ :=strconv.ParseFloat(model["amount"].(string),20)
	if record.Amount!=amount {
		return nil,errors.New("金额错误!!") 
	}
	//审核不通过
	failRes  :=model["fail_res"].(string)
	if len(failRes)>1 {
		rechargeRecord.UpdateStatusAuditWithNoFail(model["audit"].(string),model["no"].(string),appId,failRes)
		return nil,nil
	}
	
	log.Info(len(failRes))
	rechargeRecord.UpdateStatusAuditWithNo(model["audit"].(string),model["no"].(string),appId)
	
	if record.Amount>0 {
		params := map[string]interface{}{
			"out_trade_no": record.No,
			"out_trade_type": comm.Trade_Type_Recharge,
			"open_id": record.OpenId,
			"amount": int64(record.Amount*100),
			"trade_type": comm.Trade_Type_Recharge,
			"pay_type": 3,
			"title": record.Remark,
			"client_ip": "127.0.0.1",
			"notify_url": config.GetValue("notify_url").ToString(),
			"remark": record.Remark,
		}
		resultMap,err :=RequestPayApi("/pay/makeprepay",params)
		//冻结金额增加 freeze_money
		if err==nil && record.Type == 4 {
			dao.NewAccount().AccountAddFreezeMoney(record.OpenId,int64(record.Amount*100))
		}
		//推送
		PushSingle(record.OpenId,appId,"充值成功",fmt.Sprintf("充值%.2f成功",record.Amount),"paySucceed")
		return resultMap,err
	}else{
		account := dao.NewAccount()
		account,err =account.AccountWithOpenId(record.OpenId,appId)
		if err!=nil{
			return nil,err
		}
		//制作预付款
		payparams := map[string]interface{}{
			"open_id":record.OpenId,
			"type": 1,
			"amount": 0-int64(record.Amount*100),
			"title": record.Remark,
			"remark": record.Remark,
		}
		resultImprestMap,err := RequestPayApi("/pay/makeimprest",payparams)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
		//payToken
		payparams = map[string]interface{}{
			"open_id": record.OpenId,
			"password": account.Password,
		}
		resultMap,err := RequestPayApi("/pay/token",payparams)
		if err!=nil{
			return nil,err
		}
		//支付预付款
		payparams = map[string]interface{}{
			"pay_token": resultMap["pay_token"].(string),
			"open_id": record.OpenId,
			"code": resultImprestMap["code"],
			"out_trade_no":"",
		}
		resultMap,err = RequestPayApi("/pay/payimprest",payparams)
		//冻结金额减少 freeze_money
		if err==nil && record.Type == 6 {
			money:=account.FreezeMoney+int64(record.Amount*100)
			if money<0 {
				money=0
			}	
			dao.NewAccount().AccountMinusFreezeMoney(record.OpenId,money)
		}
		//推送
		PushSingle(record.OpenId,appId,"充值成功",fmt.Sprintf("充值%.2f成功",record.Amount),"paySucceed")
		return resultMap,err
	}	
	return nil,errors.New("操作错误!") 
}
//zsum  统计充值的金额
func RechargeRecordZsum(appId string,froms int64,search dao.AccountRechargeSearch) (string,error) {
	return dao.NewAccountRecharge().RechargeRecordZsum(appId,froms,search)
}
//fsum  统计扣款的金额
func RechargeRecordFsum(appId string,froms int64,search dao.AccountRechargeSearch) (string,error) {
	return dao.NewAccountRecharge().RechargeRecordFsum(appId,froms,search)
}
//获取全部用户
func Accounts(appId string) ([]*dao.Account,error)  {
	return dao.NewAccount().Accounts(appId)
}
//
func MakeCashout(appId string,cashOut interface{}) error  {
	sesson := db.NewSession()
	tx,_   :=sesson.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	
	cashoutId,err:=dao.MakeCashout(appId,cashOut,tx)	
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	 cashout:=cashOut.(dao.Cashout)

	resultMap,err :=RequestPayApi("/makecashout",map[string]interface{}{
		"open_id"	: cashout.OpenId,
		"amount"	: cashout.Amount,
		"title"		: cashout.Title,
		"remark"	: cashout.Remark,
	})
	
	if err!=nil{
		tx.Rollback()
		return err
	}	
	err=dao.CashoutcodeUpdateTx(cashoutId,resultMap["cashout_code"].(string),tx)	
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	err = tx.Commit()	
	if err!=nil{
		tx.Rollback()
		return err
	}	

	return nil
}
func Cashout(cashoutId string) error  {
	sesson := db.NewSession()
	tx,_   :=sesson.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	
	record,err:=dao.CashoutRecordTx(cashoutId,tx)
	
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	if record.Status>0{
		tx.Rollback()
		return errors.New("已审核!") 
	}
	
	/* resultMap,err :=RequestPayApi("/makecashout",map[string]interface{}{
		"open_id"	: record.OpenId,
		"amount"	: record.Amount,
		"title"		: record.Title,
		"remark"	: record.Remark,
	})
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	err=dao.CashoutcodeUpdateTx(cashoutId,resultMap["cashout_code"].(string),tx)	
	if err!=nil{
		tx.Rollback()
		return err
	} */
	
	resultMap,err :=RequestPayApi("/cashout",map[string]interface{}{
		"open_id"	: record.OpenId,
		"code"		: string(record.Cashoutcode),
	})
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	err=dao.CashoutStatusUpdateTx(cashoutId,resultMap,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}	
	
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	return nil	
}
func CashoutRecord(appId string,pageIndex uint64,pageSize uint64,mobile string,openId string,status string) ([]*dao.Cashout,int64,error)  {
	items,err:=dao.CashoutRecord(appId,pageIndex,pageSize,mobile,openId,status)
	count	 :=dao.CashoutRecordCount(appId,mobile,openId,status)
	return items,count,err
}
func UpdateGetui(openId string,cid string,devicetoken string) error  {
	return dao.NewAccount().UpdateGetui(openId,fmt.Sprintf("{\"cid\":\"?\",\"devicetoken\":\"?\"}",cid,devicetoken))
}
func AccountSyncMoney(openId string,money int64) error {
	return dao.NewAccount().UpdateMoney(openId,money)
}


































