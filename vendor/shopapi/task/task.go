package task

import (
	"errors"
	"github.com/robfig/cron"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"shopapi/comm"
	"shopapi/dao"
	"shopapi/service"
	"time"
	"strconv"
)

const (
	//订单自动取消时间
	//ORDER_AUTO_CANCEL_TIME = 5 //单位分钟
	ORDER_AUTO_CANCEL_TIME = 1440 //单位分钟
	//订单结算时间
	//ORDER_CAL_MAX_TIME = 30 //单位分钟
	ORDER_CAL_MAX_TIME = 21600 //单位分钟
)

func  StartCron() {

	c := cron.New()

	c.AddFunc("0 0/6 * * * ?", OrderFetchMoney)

	c.AddFunc("0 0/2 * * * ?", OrderAutoCancel)
	//厨师随机增加服务数量 0到2
	c.AddFunc("0 0 9 * * ?", MerchantServiceAdd)

	c.Start()
}

// 订单结算
func OrderFetchMoney() {

	order := dao.NewOrder()
	tm := time.Now().Add(-time.Minute * ORDER_CAL_MAX_TIME)
	stm := qtime.ToyyyyMMddHHmm(tm)
	log.Info("-----------时间--------", stm)
	orders, err := order.OrderWithStatusLTTime(comm.ORDER_PAY_STATUS_SUCCESS, comm.ORDER_STATUS_WAIT_SURE, stm)
	if err != nil {
		log.Error(err)
		return
	}
	if orders != nil && len(orders) > 0 {
		log.Error("订单数量:", len(orders))

		for _, order := range orders {
			//结算商户
			err := calMerchant(order)
			if err != nil {
				log.Error("商户结算失败,订单号:", order.No)
				return
			}

			if order.DbnAmount > 0 {
				//结算分销者
				err = calDBN(order)
				if err != nil {
					log.Error("分销者结算失败,订单号:", order.No)
					return
				}
			}

			err = order.UpdateWithOrderStatus(comm.ORDER_STATUS_SURED, order.No)
			if err != nil {
				log.Error("订单号:", order.No, err)
				continue
			}

		}
	} else {
		log.Warn("没有要结算的订单！")
	}
}

//结算商户
func calMerchant(order *dao.Order) error {

	return cal(order, order.MerchantAmount, order.MOpenId, "merchant")
}

//结算分销商
func calDBN(order *dao.Order) error {

	if order.DbnAmount <= 0 {
		return nil
	}

	orderItems, err := dao.NewOrderItem().OrderItemWithOrderNo(order.No)
	if err != nil {
		return err
	}
	if orderItems == nil && len(orderItems) <= 0 {
		return nil
	}
	for _, orderItem := range orderItems {
		dbnNo := orderItem.DbnNo
		if dbnNo != "" {
			userdbn, err := dao.NewUserDistribution().WithCode(dbnNo)
			if err != nil {
				log.Error(err)
				return err
			}
			if userdbn == nil {
				return errors.New("分销者信息没有找到！")
			}


			err = cal(order, order.DbnAmount, userdbn.OpenId, "dbn"+strconv.FormatInt(orderItem.Id,10))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//订单结算
//calMoney 结算金额
//mark 标记
func cal(order *dao.Order, calMoney float64, openId string, mark string) error {
	params := map[string]interface{}{
		"open_id":  openId, //需要领取金额的用户ID
		"code":     order.Code, // 预支付code
		"amount":   int64(calMoney * 100), // 领取金额
		"title":    order.Title, //标题
		"remark":   order.Title, //备注
		"out_code": order.No + "-" + mark, //第三方唯一标示
	}
	_, err := service.RequestPayApi("/v2/imprests/fetch", params)

	return err
}

//订单自动取消
func OrderAutoCancel() {
	log.Info("开始查询待取消的订单！")
	order := dao.NewOrder()
	tm := time.Now().Add(-time.Minute * ORDER_AUTO_CANCEL_TIME)
	stm := qtime.ToyyyyMMddHHmm(tm)
	orders, err := order.OrderWithNoPayAndLTTime(stm)
	if err != nil {
		log.Error(err)
		return
	}

	if orders != nil && len(orders) > 0 {

		for _, order := range orders {
			err := service.OrderAutoCancel(order.No, order.AppId)
			if err != nil {
				log.Error(err)
				continue
			}
		}
		log.Info("取消订单数:",len(orders))
	} else {
		log.Warn("没有需要取消的订单")
	}

}

//厨师随机增加服务数量 0到2
func MerchantServiceAdd() {
	log.Info("厨师随机增加服务数量 0~2")
	dao.MerchantServiceAdd()
}
