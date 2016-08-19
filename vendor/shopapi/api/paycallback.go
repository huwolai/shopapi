package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/comm"
	"shopapi/service"
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

	log.Debug("支付回调....")

	var params CallbackDto
	err :=c.BindJSON(&params)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数错误!")
		return
	}

	if params.TradeType == comm.Trade_Type_Recharge {
		log.Debug("交易充值")
		model := &service.ImprestsModel{}
		model.Amount = params.Amount
		model.OpenId = params.OpenId
		model.Code = params.Code
		model.Remark = params.Remark
		model.Title = params.Title
		err :=service.FetchImprests(model)
		if err!=nil{
			log.Error(err)
			util.ResponseError400(c.Writer,"领取失败!")
			return
		}else{
			log.Debug("预付款领取成功")
			util.ResponseSuccess(c.Writer)
		}
	}else if params.TradeType == comm.Trade_Type_Buy { //购买交易
		log.Debug("购买交易")
	}
}
