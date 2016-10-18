package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/comm"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
)

type CallbackDto struct  {
	//交易号
	TradeNo string `json:"trade_no"`
	//交易类型 1.充值 2.普通支出
	TradeType int `json:"trade_type"`
	//第三方系统中的交易号
	OutTradeNo string `json:"out_trade_no"`
	//预付款代号
	Code string `json:"code"`
	//第三方系统中的交易类型
	OutTradeType int `json:"out_trade_type"`
	//应用ID
	AppId string  `json:"app_id"`
	//用户openID
	OpenId string `json:"open_id"`
	//交易时间
	TradeTime string `json:"trade_time"`
	//交易金额
	Amount int64  `json:"amount"`
	//交易标题
	Title string `json:"title"`
	//交易备注
	Remark string `json:"remark"`
	//交易通知地址
	NotifyUrl string `json:"notify_url"`

	//是否必须一次付清
	NoOnce int `json:"no_once"`

}

func CallbackForPayapi(c *gin.Context)  {

	log.Info("支付回调....")

	var params CallbackDto
	err :=c.BindJSON(&params)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数错误!")
		return
	}

	if params.TradeType == comm.Trade_Type_Recharge {
		log.Debug("交易充值")
		//#bug 提取预付款接口需要做暝等
		model := &service.ImprestsModel{}
		model.Amount = params.Amount
		model.OpenId = params.OpenId
		model.Code = params.Code
		model.Remark = params.Remark
		model.Title = params.Title
		subTradeNo,err :=service.FetchImprests(model)
		if err!=nil{
			log.Error(err)
			util.ResponseError400(c.Writer,"领取失败!")
			return
		}else{
			log.Debug("预付款领取成功")
			err :=PublishAccountRechargeEvent(subTradeNo,model.OpenId,params.AppId,params.Amount)
			if err!=nil{
				log.Error(err)
			}			
			util.ResponseSuccess(c.Writer)
			return
		}
	}else if params.TradeType == comm.Trade_Type_Buy { 	   //购买交易
		log.Debug("购买交易")
		
		//更新订单状态
		log.Debug("OutTradeNo")
		err =service.UpdateToPayed(params.OutTradeNo,"shopapi")
		if err!=nil{
			log.Error(err)
		}	
		
		util.ResponseSuccess(c.Writer)
		return		
	}
}

//发布账户充值事件
func PublishAccountRechargeEvent(subTradeNo string,openId string,appId string,changeAmount int64) error  {

	accountEvent :=queue.NewAccountEvent()
	accountEvent.EventName="账户金额发送变化"
	accountEvent.EventKey = queue.ACCOUNT_AMOUNT_EVENT_CHANGE
	eventContent :=queue.NewAccountEventContent()
	eventContent.Action = "ACCOUNT_RECHARGE"
	eventContent.SubTradeNo = subTradeNo
	eventContent.AppId = appId
	eventContent.OpenId =openId
	eventContent.ChangeAmount = changeAmount
	accountEvent.Content = eventContent
	err :=queue.PublishAccountEvent(accountEvent)

	return err
}
