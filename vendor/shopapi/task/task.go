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
	"fmt"
)

const (
	//订单自动取消时间
	//ORDER_AUTO_CANCEL_TIME = 5 //单位分钟
	ORDER_AUTO_CANCEL_TIME = 1440 //单位分钟
	ORDER_AUTO_CANCEL_TIME_MALL = 1440 //商城订单 单位分钟
	ORDER_AUTO_CANCEL_TIME_CHEF = 30 //厨师服务订单 单位分钟
	//订单结算时间
	//ORDER_CAL_MAX_TIME = 30 //单位分钟
	ORDER_CAL_MAX_TIME = 21600 //单位分钟
	ORDER_CAL_MAX_TIME_MALL = 21600 //商城订单 单位分钟
	ORDER_CAL_MAX_TIME_CHEF = 10080 //厨师服务订单 单位分钟
)

func  StartCron() {

	c := cron.New()

	c.AddFunc("0 0/6 * * * ?", OrderFetchMoney)

	c.AddFunc("0 0/2 * * * ?", OrderAutoCancel)
	
	//厨师随机增加服务数量 0到2
	c.AddFunc("0 0 9 * * ?", MerchantServiceAdd)
	//商品初始化售出数量
	c.AddFunc("0 0 1 * * ?", ProductInitNum)
	//商品 售出数量 定时增加
	c.AddFunc("0 0 1 * * ?", ProductAddNum)

	//c.AddFunc("0 0/1 * * * ?", PurchaseCodesOpen)
	
	c.Start()
}

// 订单结算
func OrderFetchMoney() {
	order := dao.NewOrder()
	tm := time.Now().Add(-time.Minute * ORDER_CAL_MAX_TIME_CHEF)
	stm := qtime.ToyyyyMMddHHmm(tm)
	log.Info("-----------时间--------", stm)
	orders, err := order.OrderWithStatusLTTime(comm.ORDER_PAY_STATUS_SUCCESS, comm.ORDER_STATUS_WAIT_SURE, stm)
	if err != nil {
		log.Error(err)
		return
	}
	if orders != nil && len(orders) > 0 {
		cancel:=0
		for _, order := range orders {
			//========================================
			//orderType,err:=order.OrderTypes(order.No,order.AppId)
			orderType,err:=dao.NewOrderItem().OrderItemWithOrderNo(order.No)
			if err != nil {
				log.Error(err)
				continue
			}
			if req2map, err := dao.JsonToMap(orderType[0].Json); err == nil {
				if req2map["goods_type"]=="mall"{
					//商城订单
					//log.Info("mall")					
					if order.UpdateTimeUnix+ORDER_CAL_MAX_TIME_MALL*60>time.Now().Unix() {
						continue
					}					
				}else if req2map["goods_type"]=="chef"{
					//厨师订单
					//log.Info("chef")
					if order.UpdateTimeUnix+ORDER_CAL_MAX_TIME_CHEF*60>time.Now().Unix() {
						continue
					}
				}
			}
			//========================================		
			//结算商户
			err = calMerchant(order)
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
			cancel=cancel+1
		}
		log.Info("结算订单数:",cancel)	
	} else {
		log.Warn("没有要结算的订单！")
	}

	/* order := dao.NewOrder()
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
	} */
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
				log.Error(order.DbnAmount)
				log.Error(int64(order.DbnAmount * 100))
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
		"out_trade_no": order.No,
	}
	_, err := service.RequestPayApi("/v2/imprests/fetch", params)

	if err==nil {
		accountRecharge 		:= dao.NewAccountRecharge()
		accountRecharge.Amount   = calMoney * 100
		accountRecharge.No 		 = order.No + "-" + mark
		accountRecharge.AppId 	 = order.AppId
		accountRecharge.Status 	 = comm.ACCOUNT_RECHARGE_STATUS_WAIT
		accountRecharge.OpenId 	 = openId
		accountRecharge.Froms 	 = 0
		accountRecharge.Opt 	 = "system"
		//accountRecharge.audit 	 = ""
		accountRecharge.Remark   = order.Title
		accountRecharge.Type     = 1
		accountRecharge.Insert()
		/* err :=accountRecharge.Insert()
		if err!=nil{
			log.Error(err)
			return nil,errors.New("充值记录插入失败!")
		} */
	}
	
	
	return err
}

//订单自动取消
func OrderAutoCancel() {	
	order := dao.NewOrder()
	tm := time.Now().Add(-time.Minute * ORDER_AUTO_CANCEL_TIME_CHEF)
	stm := qtime.ToyyyyMMddHHmm(tm)
	orders, err := order.OrderWithNoPayAndLTTime(stm)
	if err != nil {
		log.Error(err)
		return
	}
	cancel:=0
	if orders != nil && len(orders) > 0 {
		for _, order := range orders {
			//========================================
			//orderType,err:=order.OrderTypes(order.No,order.AppId)
			orderType,err:=dao.NewOrderItem().OrderItemWithOrderNo(order.No)
			if err != nil {
				log.Error(err)
				continue
			}
			if req2map, err := dao.JsonToMap(orderType[0].Json); err == nil {
				if req2map["goods_type"]=="mall"{
					//商城订单
					//log.Info("mall")					
					if order.UpdateTimeUnix+ORDER_AUTO_CANCEL_TIME_MALL*60>time.Now().Unix() {
						continue
					}					
				}else if req2map["goods_type"]=="chef"{
					//厨师订单
					//log.Info("chef")
					if order.UpdateTimeUnix+ORDER_AUTO_CANCEL_TIME_CHEF*60>time.Now().Unix() {
						continue
					}
				}
			}
			//========================================
			err = service.OrderAutoCancel(order.No, order.AppId)
			if err != nil {
				log.Error(err)
				continue
			}
			cancel=cancel+1
			log.Info("取消订单:",order.No)	
			/* err := service.OrderAutoCancel(order.No, order.AppId)
			if err != nil {
				log.Error(err)
				continue
			} */
		}
		log.Info("取消订单数:",cancel)		
	}else{
		log.Warn("没有需要取消的订单")
	}
	/* //log.Info("开始查询待取消的订单！")
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
		//log.Warn("没有需要取消的订单")
	} */
}

//厨师随机增加服务数量 0到2
func MerchantServiceAdd() {
	log.Info("厨师随机增加服务数量 0~2")
	dao.MerchantServiceAdd()
}
//商品初始化售出数量
func ProductInitNum()  {
	log.Info("商品初始化售出数量！")
	dao.ProductInitNum()
}
//商品 售出数量 定时增加
func ProductAddNum()  {
	log.Info("商品 售出数量 定时增加!")
	dao.ProductAddNum()
}
//定时开奖
func PurchaseCodesOpen()  {
	//公式 
	log.Info("定时开奖")
	prodOpening,_	:=dao.ProdPurchaseCodeWithOpenStatus(1)
	if prodOpening!=nil {
		for _, prod := range prodOpening {	
			if len(prod.OpenMobile)>0 {
				return 
			}
			orderItem,_:=dao.OrderItemPurchaseCodesWithTime(prod.OpenTime,comm.PRODUCT_YYG_BUY_CODES)
			codeCount,_:=dao.OrderItemPurchaseCodesWithProdId(prod.ProdId)//商品份数
			//codeCount++
			
			var c  int64 = 0
			for _, tSum := range orderItem {
				s1	:=fmt.Sprintf("%s%s",time.Unix(int64(tSum.BuyTime/1e3), 0).Format("150405"),dao.Right(fmt.Sprintf("%d",tSum.BuyTime),3))
				i1,_:=strconv.ParseInt(s1,10,64)			
				c	 =c+i1
			}
			openCode:=fmt.Sprintf("%d",c%codeCount+10000001)
			
			user,err:=dao.GetOpenIdbyOpenCode(prod.ProdId,openCode)		
			if err!=nil {
				log.Error(err)
				return
			}
			if user==nil {
				return
			}
			
			dao.ProductAndPurchaseCodesOpened(prod.ProdId,user.OpenId,user.Mobile,openCode)
		}
	}
	//开奖状态更新
	dao.ProductAndPurchaseCodesOpenedStatus()
}
//定时开奖 ****





























