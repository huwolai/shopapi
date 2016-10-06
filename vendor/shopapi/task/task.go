package task

import (
	"github.com/robfig/cron"
	"shopapi/dao"
	"shopapi/comm"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"errors"
)

func StartCron()  {

	c :=cron.New()

	c.AddFunc("0 0/6 * * * ?", OrderFetchMoney)

	c.AddFunc("0 0/5 * * * ?", OrderAutoCancel)

	c.Start()
}

// 订单结算
func OrderFetchMoney()  {

	order :=dao.NewOrder()
	tm :=time.Now().Add(-time.Minute*30)
	stm :=qtime.ToyyyyMMddHHmm(tm)
	log.Info("-----------时间--------",stm)
	orders,err :=order.OrderWithStatusLTTime(comm.ORDER_PAY_STATUS_SUCCESS,comm.ORDER_STATUS_WAIT_SURE,stm)
	if err!=nil{
		log.Error(err)
		return
	}
	if orders!=nil&&len(orders) >0 {
		log.Error("订单数量:",len(orders))

		for _,order :=range orders  {
			//结算商户
			err := calMerchant(order)
			if err!=nil{
				log.Error("商户结算失败,订单号:",order.No)
				return
			}

			if order.DbnAmount > 0 {
				//结算分销者
				err =calDBN(order)
				if err!=nil{
					log.Error("分销者结算失败,订单号:",order.No)
					return
				}
			}

			err =order.UpdateWithOrderStatus(comm.ORDER_STATUS_SURED,order.No)
			if err!=nil{
				log.Error("订单号:",order.No,err)
				continue
			}

		}
	}else{
		log.Warn("没有要结算的订单！")
	}
}

//结算商户
func calMerchant(order *dao.Order) error {

	return cal(order,order.MerchantAmount,order.MOpenId,"merchant")
}

//结算分销商
func calDBN(order *dao.Order) error  {

	if order.DbnAmount<=0 {
		return nil
	}

	orderItems,err :=dao.NewOrderItem().OrderItemWithOrderNo(order.No)
	if err!=nil{
		return err
	}
	if orderItems==nil&&len(orderItems)<=0{
		return nil
	}
	for _,orderItem :=range orderItems  {
		dbnNo :=orderItem.DbnNo
		if dbnNo!="" {
			userdbn,err :=dao.NewUserDistribution().WithCode(dbnNo)
			if err!=nil{
				log.Error(err)
				return err
			}
			if userdbn==nil{
				return errors.New("分销者信息没有找到！")
			}
			err = cal(order,order.DbnAmount,userdbn.OpenId,"dbn"+orderItem.Id)
			if err!=nil{
				return err
			}
		}
	}

	return nil
}

//订单结算
//calMoney 结算金额
//mark 标记
func cal(order *dao.Order,calMoney float64,openId string,mark string) error {
	params :=map[string]interface{}{
		"open_id":openId,
		"code": order.Code,
		"amount": int64(calMoney*100),
		"title": order.Title,
		"remark": order.Title,
		"out_code": order.No + "-" +mark,
	}
	_,err :=service.RequestPayApi("/v2/imprests/fetch",params)

	return err
}

//订单自动取消
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
		log.Warn("没有需要取消的订单")
	}

}
