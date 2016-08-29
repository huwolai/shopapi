package task

import (
	"github.com/robfig/cron"
	"shopapi/dao"
	"shopapi/comm"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
)

func StartCron()  {

	c :=cron.New()
	//没半个小时执行一次
	c.AddFunc("0 0/6 * * * ?", OrderFetchMoney)

	c.AddFunc("0 0/5 * * * ?", OrderAutoCancel)

	c.Start()
}

func OrderFetchMoney()  {

	order :=dao.NewOrder()
	tm :=time.Now().Add(-time.Minute*30)
	stm :=qtime.ToyyyyMMddHHmm(tm)
	log.Error("-----------时间--------",stm)
	orders,err :=order.OrderWithStatusLTTime(comm.ORDER_PAY_STATUS_SUCCESS,comm.ORDER_STATUS_WAIT_SURE,stm)
	if err!=nil{
		log.Error(err)
		return
	}
	if orders!=nil&&len(orders) >0 {
		log.Error("订单数量:",len(orders))

		for _,order :=range orders  {
			 params :=map[string]interface{}{
				"open_id":order.OpenId,
				 "code": order.Code,
				 "amount": int64(order.PayPrice*100),
				 "title": order.Title,
				 "remark": order.Title,
			}
			_,err :=service.RequestPayApi("/imprests/fetch",params)
			if err!=nil{
				log.Error("订单号:",order.No,err)
				continue
			}
			err =order.UpdateWithOrderStatus(comm.ORDER_STATUS_SURED,order.No)
			if err!=nil{
				log.Error("订单号:",order.No,err)
				continue
			}

		}
	}else{
		log.Warn("木有获取到订单")
	}
}

func OrderAutoCancel()  {

	order :=dao.NewOrder()
	tm :=time.Now().Add(-time.Minute*30)
	stm :=qtime.ToyyyyMMddHHmm(tm)
	orders,err :=order.OrderWithNoPayAndLTTime(stm)
	if err!=nil{
		log.Error(err)
		return
	}

	if orders!=nil&&len(orders)>0 {
		for _,order :=range orders  {
			err :=service.OrderCancel(order.No,"",order.AppId)
			if err!=nil{
				log.Error(err)
				continue
			}
		}
	}else{
		log.Warn("OrderAutoCancel","木有获取到订单")
	}

}
