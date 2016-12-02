package service

import (
	"shopapi/dao"
	"errors"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"encoding/json"
	
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"fmt"
	"net/url"
	"crypto/md5"
    "encoding/hex"
	"encoding/base64"
	"strings"
)

type OrderModel struct  {
	Items []OrderItemModel
	Json string
	AddressId int64
	Address string
	OpenId string
	AppId string
	Title string
	MOpenId string
	MerchantId int64
	OrderStatus int
	PayStatus int
}

type OrderItemModel struct  {
	//分销编号
	DbnNo string
	//商品sku
	SkuNo string
	//商品数量
	Num int
	Json string
}

type OrderPrePayModel struct  {
	OrderNo string
	CouponTokens []string
	AddressId int64
	PayType int
	AppId string
	Json string
	NotifyUrl string
}

type OrderPrePayDto struct  {
	//优惠token
	CouponTokens []string `json:"coupon_tokens"`
	OrderNo string `json:"order_no"`
	AddressId int64 `json:"address_id"`
	Json string `json:"json"`
	//付款类型(1.支付宝 2.微信 3.现金支付 4.账户)
	PayType int `json:"pay_type"`
	AppId string `json:"app_id"`
}

type OrderCouponDto struct  {
	AppId string `json:"app_id"`
	//订单号
	OrderNo string `json:"order_no"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//商户open_id
	MOpenId string `json:"m_open_id"`
	//下单用户
	OpenId string `json:"open_id"`
	//订单标题
	Title string `json:"title"`
	//付款方式
	PayMethod int `json:"pay_method"`
	Flag string `json:"flag"`
	Json string `json:"json"`
	//订单实际金额(此金额为实际付款金额)
	ActPrice float64 `json:"act_price"`
	//订单价格
	Price float64 `json:"price"`
}

type OrderItemCouponDto struct {
	//订单号
	OrderNo string `json:"order_no"`
	ProdId int64 `json:"prod_id"`
	SkuNo string `json:"sku_no"`
	Num int `json:"num"`
	Flag string `json:"flag"`
	Json string `json:"json"`
	BuyTotalPrice float64 `json:"buy_total_price"`

}

type OrderNo struct {
	OrderNo	[]string   `json:"order_no"`
}

func OrderAdd(model *OrderModel) (*dao.Order,error)  {
	sess :=db.NewSession()
	tx,_ := sess.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()

	order,err := orderSave(model,tx)
	if err !=nil {
		tx.Rollback()
		return nil,err
	}
	tx.Commit()
	return order,nil

}

func OrderPrePay(model *OrderPrePayModel) (map[string]interface{},error) {

	if model.PayType==comm.Pay_Type_Cash { //现金支付

		return nil,errors.New("暂不支持此支付方式!")
	}
	order := dao.NewOrder()
	order,err :=order.OrderWithNo(model.OrderNo,model.AppId)
	if err!=nil {
		return nil,err
	}
	if order==nil{
		return nil,errors.New("没有找到对应的订单信息!")
	}
	address := dao.NewAddress()
	address,err = address.WithId(model.AddressId)
	if err!=nil{
		log.Error(err)
		return nil,errors.New("查询地址失败!")
	}
	if address==nil{
		return nil,errors.New("没有找到对应的地址信息!")
	}

	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			log.Error(err)
			tx.Rollback()
		}
	}()

	var couponTotalAmount float64 = 0
	//存在优惠信息
	if model.CouponTokens!=nil&&len(model.CouponTokens) >0 {
		couponTotalAmount,err =HandleCoupon(order,model.CouponTokens,tx)
		if err!=nil{
			log.Error(err)
			tx.Rollback()
			return nil,errors.New("优惠信息处理错误!")
		}
	}
	//实际付款金额
	payPrice :=order.RealPrice - couponTotalAmount
	if payPrice<0 {
		tx.Rollback()
		return nil,errors.New("付款金额不能为负数!")
	}
	//计算订单金额
	err =calOrderAmount(order,payPrice,couponTotalAmount,tx)
	if err!=nil{
		tx.Rollback()
		log.Error(err)
		return nil,errors.New("计算订单金额失败!")
	}

	if model.PayType == comm.Pay_Type_Account {//账户支付
		code :=order.Code
		if order.PayStatus==comm.ORDER_PAY_STATUS_NOPAY || order.PayStatus==comm.ORDER_PAY_STATUS_PAYING {
			//请求预付款
			resultImprestMap,err := makeImprest(order,address,payPrice)
			if err!=nil{
				tx.Rollback()
				log.Error(err)
				return nil,err
			}
			code :=resultImprestMap["code"].(string)
			err =order.OrderPayapiUpdateWithNoAndCodeTx("",address.Id,address.Address,address.Name,address.Mobile,code,comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_PAYING,order.No,model.Json,order.AppId,tx)
			if err!=nil{
				log.Error(err)
				tx.Rollback()
				return nil,err
			}
		}
		resultMap := map[string]interface{}{
			"open_id": order.OpenId,
			"code": code,
		}
		tx.Commit()
		return resultMap,nil
	}else{
		resultPrepayMap,err := makePrepay(order,address,payPrice,model)
		if err!=nil{
			tx.Rollback()
			log.Error(err)
			return nil,err
		}
		if resultPrepayMap!=nil{
			payapiNo :=resultPrepayMap["pay_no"].(string)
			code :=resultPrepayMap["code"].(string)
			//将payapi的订单号更新到订单数据里
			err :=order.OrderPayapiUpdateWithNoAndCodeTx(payapiNo,address.Id,address.Address,address.Name,address.Mobile,code,comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_PAYING,order.No,model.Json,order.AppId,tx)
			if err!=nil{
				log.Error(err)
				tx.Rollback()
				return nil,err
			}
		}
		tx.Commit()
		return resultPrepayMap,nil

	}

}

//确认订单
func OrderSure(orderNo,appId string) error  {
	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return errors.New("订单查询错误!")
	}
	if order==nil{
		return errors.New("订单没找到!")
	}
	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err:=recover();err!=nil{
			log.Error(err)
			tx.Rollback()
		}
	}()
	err =order.UpdateWithOrderStatusTx(comm.ORDER_STATUS_SURED,order.No,tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return err
	}

	//订单的钱结算
	err = allocOrderAmount(order)
	if err!=nil{
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//分订单金额
func allocOrderAmount(order *dao.Order) error  {

	orderItem := dao.NewOrderItem()
	items,err := orderItem.OrderItemWithOrderNo(order.No)
	if err!=nil{
		log.Error(err)
		return err
	}
	if items==nil{
		return errors.New("没有找到订单明细数据!")
	}

	//分配给商户的钱
	imprestsModel := &ImprestsModel{}
	imprestsModel.Code = order.Code
	imprestsModel.Amount = int64(order.MerchantAmount*100)
	imprestsModel.OpenId = order.MOpenId
	imprestsModel.Title="厨师上门服务" //分销系统待区分
	imprestsModel.Remark = "厨师上门服务费用"
	_,err =FetchImprests(imprestsModel)
	if err!=nil{
		log.Error(err)
		log.Error("syserr->订单号[",order.No,"]","商户ID[",order.MOpenId,"]", "结算商户的钱失败!,导致结算给分销者的钱未成功!严重问题")
		return err
	}

	if order.DbnAmount<=0 {
		log.Warn("此订单",order.No,"没有需要分配给分销者的钱!")
		return nil
	}
	//分配给分销者的钱
	distribMap :=make(map[string]float64)
	for _,oItem :=range items {
		if oItem.DbnNo != "" {
			distribution := dao.NewUserDistributionDetail()
			distribution, err := distribution.WithCode(oItem.DbnNo)
			if err != nil {
				log.Error(err)
				return errors.New("查询分销信息失败!")
			}
			if distribution==nil{
				log.Warn("分销编号:",oItem.DbnNo,"没有找到!")
				continue
			}
			disamount := distribMap[distribution.OpenId]
			disamount+=oItem.DbnAmount
			distribMap[distribution.OpenId] = disamount
		}
	}

	if len(distribMap)>0 {
		for key,value :=range distribMap {
			//分配给商户的钱
			imprestsModel := &ImprestsModel{}
			imprestsModel.Code = order.Code
			imprestsModel.Amount =int64(value*100)
			imprestsModel.OpenId = key
			imprestsModel.Title= "分销佣金"
			imprestsModel.Remark = "分销佣金"
			_,err =FetchImprests(imprestsModel)
			if err!=nil{
				log.Error(err)
				log.Error("syserr->订单号[",order.No,"]","分销者ID[",key,"]", "结算商分销者的钱失败!,可能导致此订单的后面的分销者没结算到钱!严重问题")
				return err
			}
		}
	}


	return nil

}

//计算订单金额
func calOrderAmount(order *dao.Order,payPrice float64,couponTotalAmount float64,tx *dbr.Tx)  error {
	orderItem := dao.NewOrderItem()
	orderitems,err :=orderItem.OrderItemWithOrderNo(order.No)
	if err!=nil{
		log.Error(err)
		return errors.New("查询订单明细失败!")
	}
	var totaldbnAmount float64 //总佣金
	var totalOmitMoney float64 //省略的金额
	var totalMerchantAmount float64 //商户应得的金额
	if  orderitems!=nil{
		for _,oItem :=range orderitems {
			if oItem.DbnNo != "" {
				distribution := dao.NewUserDistributionDetail()
				distribution, err := distribution.WithCode(oItem.DbnNo)
				if err != nil {
					log.Error(err)
					return errors.New("查询分销信息失败!")
				}
				if distribution==nil{
					log.Warn("分销编号:",oItem.DbnNo,"没有找到!")
					continue
				}
				dbnAmount := payPrice * distribution.CsnRate
				oItem.DbnAmount = comm.Floor(dbnAmount, 2)

			}
			var couponAmount float64 = 0
			if order.RealPrice!=0 {
				couponAmount = (oItem.BuyTotalPrice / order.RealPrice) * couponTotalAmount

			}
			oItem.CouponAmount = comm.Floor(couponAmount, 2)
			//==========================
			if oItem.BuyTotalPrice==180 {
				oItem.MerchantAmount =150
			}else if oItem.BuyTotalPrice==240 {
				oItem.MerchantAmount =180
			}else{
				oItem.MerchantAmount = oItem.BuyTotalPrice - oItem.CouponAmount - oItem.DbnAmount
			}
			
			
			oItem.OmitMoney = 0
			err = oItem.UpdateAmountWithIdTx(oItem.DbnAmount, oItem.OmitMoney, oItem.CouponAmount, oItem.MerchantAmount, oItem.Id, tx)
			if err != nil {
				log.Error(err)
				return errors.New("更新订单详情失败!")
			}
			totaldbnAmount += oItem.DbnAmount
			totalOmitMoney += oItem.OmitMoney
			totalMerchantAmount += oItem.MerchantAmount
		}
	}
	//更新对应金额
	err =order.UpdateAmountTx(couponTotalAmount,payPrice,totalMerchantAmount,totalOmitMoney,totaldbnAmount,order.No,tx)
	if err!=nil{
		log.Error(err)
		return errors.New("更新订单信息失败!")
	}

	return nil
}

//制作预支付
func makePrepay(order *dao.Order,address *dao.Address,payPrice float64,model *OrderPrePayModel)  (map[string]interface{},error)  {
	//参数
	params := map[string]interface{}{
		"open_id": order.OpenId,
		"out_trade_no": order.No,
		"amount": int64(payPrice*100),
		"trade_type": 2,  //交易类型(1.充值 2.购买)
		"pay_type": model.PayType,
		"title": order.Title,
		"client_ip": "127.0.0.1",
		"notify_url": model.NotifyUrl,
		"remark": order.Title,
	}
	resultPrepayMap,err :=RequestPayApi("/pay/makeprepay",params)

	return resultPrepayMap,err

}

//制作预付款
func makeImprest(order *dao.Order,address *dao.Address,payPrice float64) (map[string]interface{},error)  {
	params := map[string]interface{}{
		"open_id":order.OpenId,
		"type": 1,
		"amount": int64(payPrice*100),
		"title": order.Title,
		"remark": order.Title,
	}
	resultImprestMap,err := RequestPayApi("/pay/makeimprest",params)
	if err!=nil{
		return nil,err
	}


	resultMap := map[string]interface{}{
		"open_id": resultImprestMap["open_id"],
		"code": resultImprestMap["code"],
	}
	return resultMap,nil

}
//处理优惠信息
func HandleCoupon(order *dao.Order,coupotokens []string,tx *dbr.Tx) (float64,error)  {

	orderCoupon := dao.NewOrderCoupon()
	err :=orderCoupon.DeleteWithOrderNoTx(comm.ORDER_COUPON_STATUS_UNACTIVATE,order.No,tx)
	if err!=nil{
		log.Error(err)
		return 0.0,errors.New("删除原有优惠记录失败!")
	}
	//去重
	ncoupontokens := comm.RemoveDuplicatesAndEmpty(coupotokens)
	//凭证校验
	var couponTotalAmount float64
	for _,couponToken :=range ncoupontokens {
		jwtAuth := comm.InitJWTCouponBackend()
		cpToken,err :=jwtAuth.FetchCouponToken(couponToken)
		if err!=nil{
			return 0.0,err
		}
		if !cpToken.Valid {
			return 0.0,errors.New("优惠凭证无效!")
		}
		orderNo,isok :=cpToken.Claims["order_no"].(string)
		if !isok {
			return 0.0,errors.New("优惠券有误[获取order_no失败]!")
		}
		if orderNo!=order.No {
			return 0.0,errors.New("优惠凭证不是当前订单的!")
		}

		couponCode,isok :=cpToken.Claims["coupon_code"].(string)
		if !isok{
			return 0.0,errors.New("优惠券有误[获取coupon_code失败]!")
		}
		couponAmount,isok :=cpToken.Claims["coupon_amount"].(float64)
		if !isok {
			return 0.0,errors.New("优惠券有误[获取coupon_amount失败]!")
		}
		notifyUrl,isok :=cpToken.Claims["notify_url"].(string)
		if !isok {
			return 0.0,errors.New("优惠券有误[获取notify_url失败]!")
		}
		trackCode,isok :=cpToken.Claims["track_code"].(string)
		if !isok {
			return 0.0,errors.New("优惠券有误[获取track_code失败]!")
		}
		if err!=nil{
			log.Error(err)
			return 0.0,errors.New("优惠券金额有误!")
		}

		orderCoupon := dao.NewOrderCoupon()
		orderCoupon.CouponAmount = couponAmount
		orderCoupon.CouponCode = couponCode
		orderCoupon.OpenId = order.OpenId
		orderCoupon.TrackCode = trackCode
		orderCoupon.CouponToken = couponToken
		orderCoupon.OrderNo = orderNo
		orderCoupon.AppId = order.AppId
		orderCoupon.NotifyUrl = notifyUrl
		err =orderCoupon.InsertTx(tx)
		if err!=nil{
			log.Error(err)
			return 0.0,errors.New("插入优惠信息失败!")
		}
		couponTotalAmount += orderCoupon.CouponAmount
	}

	return couponTotalAmount,nil
}

func OrderByUser(openId string,orderStatus []int,payStatus []int,appId string,pageIndex uint64,pageSize uint64)  ([]*dao.OrderDetail,error)  {

	orderDetail :=dao.NewOrderDetail()
	orderDetails,err := orderDetail.OrderDetailWithUser(openId,orderStatus,payStatus,appId,pageIndex,pageSize)	
	
	for k,item :=range orderDetails {
		if len(item.Items)>0 {
			if req2map, err := dao.JsonToMap(item.Items[0].Json); err == nil {
				if req2map["goods_type"]=="mall_yyg"{
					//中奖的号码
					prodPurchaseCode,_:=dao.ProdPurchaseCodeWithProdId(item.Items[0].ProdId)
					orderDetails[k].Itemsyyg=prodPurchaseCode
					//购买的号码
					buyCodes,_:=dao.OrderItemPurchaseCodesWithNo(item.Items[0].No)
					orderDetails[k].ItemsyygBuyCodes=buyCodes
				}
			}
		}		
	}

	return orderDetails,err
}

//获取用户指定状态订单数量
func OrderWithUserAndStatusCount(openId string,orderStatus []int,payStatus []int,appId string)  (*dao.OrderCount,error)  {

	orderCount :=dao.NewOrderCount()
	orderCount,err := orderCount.OrderWithUserAndStatusCount(openId,orderStatus,payStatus,appId)

	return orderCount,err
}

//查询订单信息通过商户ID
func OrderDetailWithMerchantId(merchantId int64,orderStatus []int,payStatus []int,appId string) ([]*dao.OrderDetail,error) {
	orderDetail :=dao.NewOrderDetail()
	orderDetails,err := orderDetail.OrderDetailWithMerchantId(merchantId,orderStatus,payStatus,appId)

	return orderDetails,err
}

func OrderDetailWithNo(orderNo string,appId string) (*dao.OrderDetail,error)  {
	orderDetail :=dao.NewOrderDetail()
	orderDetail,err := orderDetail.OrderDetailWithNo(orderNo,appId)
	
	if len(orderDetail.Items)>0 {
		if req2map, err := dao.JsonToMap(orderDetail.Items[0].Json); err == nil {
			if req2map["goods_type"]=="mall_yyg"{
				prodPurchaseCode,_:=dao.ProdPurchaseCodeWithProdId(orderDetail.Items[0].ProdId)
				orderDetail.Itemsyyg=prodPurchaseCode
				
				buyCodes,_:=dao.OrderItemPurchaseCodesWithNo(orderDetail.Items[0].No)
				orderDetail.ItemsyygBuyCodes=buyCodes
			}
		}
	}
	
	return orderDetail,err
}



func OrderPayForAccount(openId string,orderNo string,payToken string,appId string) error  {

	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil {
		log.Error(err)
		return errors.New("订单查询失败!")
	}

	if order==nil{
		return  errors.New("没找到订单信息!")
	}
	if order.PayStatus!=comm.ORDER_PAY_STATUS_PAYING {
		return  errors.New("订单不是待付款状态!")
	}	

	tx,_ :=db.NewSession().Begin()

	defer func() {
		if err :=recover();err!=nil{
			log.Error(err)
			tx.Rollback()
		}
	}()
	orderItem := dao.NewOrderItem()
	orderItems,err :=orderItem.OrderItemWithOrderNo(orderNo)
	if err!=nil{
		return  errors.New("查询订单明细失败!")
	}
	if orderItems==nil{
		return  errors.New("订单明细为空!")
	}
	//调整商品库存
	err = ProdSKUStockSubWithOrder(orderItems,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	//商品累计销量增加
	err = ProdSoldNumAdd(orderItems,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	//商户权重加1
	err =MerchantWeightAdd(1,order.MerchantId,tx)
	if err!=nil{
		tx.Rollback()
		return errors.New("商户权重添加失败!")
	}
	//修改订单状态
	err =order.UpdateWithStatusTx(comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_SUCCESS,orderNo,tx)
	if err!=nil{
		log.Error("订单号:",orderNo,"状态更新为支付成功的时候失败!")
		tx.Rollback()
		return errors.New("订单更新错误!")
	}
	
	//一元购
	err = purchaseCodes(orderItems,appId,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	//支付预付款
	params := map[string]interface{}{
		"pay_token": payToken,
		"open_id": order.OpenId,
		"code": order.Code,
		"out_trade_no":order.No,
	}
	resultMap,err := RequestPayApi("/pay/payimprest",params)
	if err!=nil{
		return err
	}
	subTradeNo :=resultMap["sub_trade_no"].(string)	
	
	err =NotifyCouponServer(orderNo,appId,subTradeNo)
	if err!=nil{
		tx.Rollback()
		return errors.New("通知第三方优惠服务失败!")
	}

	err =tx.Commit()
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return errors.New("数据提交失败!")
	}

	//发布事件
	go func(){
		err =PublishOrderPaidEvent(orderItems,order)
		if err!=nil{
			log.Warn("发送订单事件失败:", err)
		}
	}()

	return nil
}

/**
  发布订单支付事件
 */
func PublishOrderPaidEvent(items []*dao.OrderItem,order *dao.Order) error {

	merchant,err := dao.NewMerchant().MerchantWithId(order.MerchantId)
	if err!=nil{
		return err
	}
	if merchant==nil{

		return  errors.New("商户信息未找到！")
	}
	if merchant.Mobile=="" {
		return errors.New("商户没有填写手机号！")
	}
	
	//是否代购买 start
	type IsDiy struct {
		IsDiy bool `json:"isDIY"`
	}	
	var isdiy IsDiy
	var diy string
	json.Unmarshal([]byte(items[0].Json), &isdiy)
	if isdiy.IsDiy {
		diy="需要"
	}else{
		diy="不需要"
	}
	//是否代购买 end
	
	orderEvent := queue.NewOrderEvent()
	//订单已付款
	orderEvent.EventKey = queue.ORDER_EVENT_PAID
	orderEvent.EventName="订单已付款事件"

	orderEventContent :=queue.NewOrderEventContent()
	orderEventContent.OpenId = order.OpenId
	orderEventContent.Amount = order.RealPrice
	orderEventContent.OrderNo = order.No
	orderEventContent.CreateTime = qtime.ToyyyyMMddHHmm(order.CreateTime)
	orderEventContent.Title = order.Title
	orderEventContent.Json = order.Json
	orderEventContent.Flag = order.Flag
	orderEventContent.ExtData = map[string]interface{}{
		// 商户手机号
		"m_mobile":merchant.Mobile,
		// 商户名称
		"m_name":merchant.Name,
		// 联系人名字
		"name": order.AddressName,
		// 用户配送地址
		"address":order.Address,
		// 联系人手机号
		"mobile": order.AddressMobile,
		// 是否代购买
		"isdiy": diy,
	}
	if items!=nil&&len(items)>0 {
		orderEventItems := make([]*queue.OrderEventItem,0)
		for _,oitem :=range items {
			eventItem :=queue.NewOrderEventItem()
			eventItem.Title = oitem.Title
			eventItem.Json = oitem.Json
			eventItem.Num = oitem.Num
			eventItem.OrderNo = oitem.No
			eventItem.Price = oitem.BuyUnitPrice
			eventItem.TotalPrice = oitem.BuyTotalPrice

			orderEventItems = append(orderEventItems,eventItem)
		}
		orderEventContent.Items = orderEventItems
	}

	orderEvent.Content = orderEventContent
	return queue.PublishOrderEvent(orderEvent)
}

//通知优惠服务
func NotifyCouponServer(orderNo string,appId string,subTradeNo string) error {

	orderCoupon :=dao.NewOrderCoupon()
	orderCoupons,err := orderCoupon.WithOrderNo(orderNo,appId)
	if err!=nil{
		return err
	}
	if orderCoupons==nil{
		log.Warn("没有优惠券信息!")
		return nil
	}

	for _,ordercn :=range orderCoupons {
		if ordercn.NotifyUrl=="" {
			log.Warn("优惠券没有通知地址!")
			continue
		}
		if ordercn.Status == comm.ORDER_COUPON_STATUS_ACTIVATED {
			log.Warn("优惠券",ordercn.CouponCode,"已使用")
			continue
		}
		requestModel :=queue.NewRequestModel()
		requestModel.NotifyUrl = ordercn.NotifyUrl
		requestModel.Data = map[string]interface{}{
			"coupon_code": ordercn.CouponCode,
			"track_code": ordercn.TrackCode,
			"open_id": ordercn.OpenId,
			"sub_trade_no": subTradeNo,
		}
		log.Warn("优惠券放入队列")
		err = queue.PublishRequestMsg(requestModel)
		if err!=nil{
			log.Error(err)
			return err
		}
	}

	return nil
}

//商户权重增加
func MerchantWeightAdd(num int,merchantId int64,tx *dbr.Tx) error {
	merchant :=dao.NewMerchant()
	return merchant.IncrWeightWithIdTx(num,merchantId,tx)
}

//商品累计销售数量递增
func ProdSoldNumAdd(orderItems []*dao.OrderItem,tx *dbr.Tx) error  {

	for _,oItem :=range orderItems {
		err :=dao.NewProduct().SoldNumInc(oItem.Num,oItem.ProdId,tx)
		if err!=nil{
			log.Error(err)
			return err
		}
		err =dao.NewProdSku().SoldNumInc(oItem.Num,oItem.SkuNo,oItem.AppId,tx)
		if err!=nil{
			log.Error(err)
			return err
		}
	}

	return nil
}

//减商品sku的库存
func ProdSKUStockSubWithOrder(orderItems []*dao.OrderItem,tx *dbr.Tx) error  {

	for _,oItem :=range orderItems {
		prodSku := dao.NewProdSku()
		prodSku,err :=prodSku.WithSkuNo(oItem.SkuNo)
		if err!=nil{
			log.Error(err)
			return  errors.New("查询订单SKU失败!")
		}
		if prodSku==nil{
			return  errors.New("没有找到对应的商品信息!")
		}

		if  prodSku.Stock < oItem.Num {
			return  errors.New("库存数量不足!")
		}
		err =prodSku.UpdateStockWithSkuNoTx(prodSku.Stock-oItem.Num,oItem.SkuNo,tx)
		if err!=nil{
			log.Error(err)
			return  errors.New("修改库存失败!")
		}
	}
	return nil
}

//减商品购买码
func purchaseCodes(orderItems []*dao.OrderItem,appId string,tx *dbr.Tx) error {
	ProdPurchaseCode := &dao.ProdPurchaseCode{}	
			
	if len(orderItems)>1 {
		return errors.New("商品购买错误!")
	}
	
	for _,oItem :=range orderItems {
		goodsType,_:=dao.JsonToMap(oItem.Json);
		
		if goodsType["goods_type"]!="mall_yyg"{
			return nil
		}
	
		ProdPurchaseCode.AppId	=oItem.AppId
		ProdPurchaseCode.Id	=oItem.Id
		ProdPurchaseCode.ProdId=oItem.ProdId
		ProdPurchaseCode.Sku	=oItem.SkuNo
		ProdPurchaseCode.Num	=oItem.Num
		
		//productDao:=dao.NewProduct()
		
		codes,err:=dao.ProductAndPurchaseCodesTx(ProdPurchaseCode,tx)
		if err!=nil || codes==nil{
			return errors.New("数据库错误!")
		}
		
		if(codes.Num<ProdPurchaseCode.Num){
			return errors.New("购买数量大于库存数量!")
		}
		//购买完成 设置开奖时间
		if(codes.Num==ProdPurchaseCode.Num){
			err=dao.ProductAndPurchaseCodesOpening(tx,ProdPurchaseCode,fmt.Sprintf("%d",time.Now().Unix()+300))//5分钟
			if err!=nil{
				return err
			}
		}
		//============================
		s:=strings.Split(codes.Codes, ",")
		ns:=s[ProdPurchaseCode.Num:]	
		ls:=s[0:ProdPurchaseCode.Num]		
		
		err=dao.ProductAndPurchaseCodesMinus(tx,codes.Id,codes.Num,len(ns),strings.Join(ns,","))
		if err!=nil{
			return err
		}
		
		var index int64 = 0
		for _,codesItem :=range ls {
			index++
			err=dao.OrderItemPurchaseCodesAdd(tx,oItem.Id,oItem.No,oItem.ProdId,codesItem,index)//strings.Join(ls,",")
			if err!=nil{
				return err
			}
		}
		
		break
	}
	return nil
}



func OrderAutoCancel(orderNo string,appId string)error  {
	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}

	if order.PayStatus!=comm.ORDER_PAY_STATUS_NOPAY &&
	    order.PayStatus!=comm.ORDER_PAY_STATUS_PAYING {
		return nil
	}
	err = order.UpdateWithOrderStatus(comm.ORDER_STATUS_CANCELED,orderNo)
	if err!=nil{
		log.Error("更新订单状态失败! 订单号:",orderNo)
		return err
	}
	log.Error("订单状态为:",order.PayStatus,"不能取消!")
	return errors.New("订单状态错误!")
}

//商户同意取消订单
func OrderAgreeCancel(orderNo string,appId string) error {
	order := dao.NewOrder()
	order, err := order.OrderWithNo(orderNo, appId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("没有找到此订单!")
	}
	if order.OrderStatus != comm.ORDER_STATUS_CANCELED_WAIT_SURE {
		return errors.New("订单状态不是等待取消确认状态!")
	}

	if order.PayStatus == comm.ORDER_PAY_STATUS_SUCCESS {
		if order.Code=="" {
			return errors.New("订单不存在预付款code!")
		}
		params :=map[string]interface{}{
			"code":order.Code,
		}
		_,err =RequestPayApi("/imprest/refund",params)
		if err!=nil{
			return err
		}
		tx,_ :=db.NewSession().Begin()
		defer func() {
			if err :=recover();err!=nil{
				tx.Rollback()
				log.Error(err)
			}
		}()
		err = order.UpdateWithOrderStatusTx(comm.ORDER_STATUS_CANCELED,orderNo,tx)
		if err!=nil{
			tx.Rollback()
			log.Error("更新订单状态失败! 订单号:",orderNo)
			return err
		}

		tx.Commit()
	}else {
		err = order.UpdateWithOrderStatus(comm.ORDER_STATUS_CANCELED,orderNo)
		if err!=nil{
			log.Error("更新订单状态失败! 订单号:",orderNo)
			return err
		}
	}

	return nil

}

func OrderRefuseCancel(orderNo string,reason string,appId string) error {
	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}
	if order==nil{
		return errors.New("没有找到此订单!")
	}
	if order.OrderStatus!=comm.ORDER_STATUS_CANCELED_WAIT_SURE {
		return errors.New("订单状态不是等待取消确认状态!")
	}

	tx,_ :=db.NewSession().Begin()

	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			return
		}
	}()
	//更新为拒绝取消订单状态
	err = order.UpdateWithOrderStatusTx(comm.ORDER_STATUS_CANCELED_REJECTED,orderNo,tx)
	if err!=nil{
		log.Error("更新订单状态失败! 订单号:",orderNo)
		tx.Rollback()
		return err
	}
	err =order.UpdateWithRefuseCancelReasonTx(reason,orderNo,tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func OrderCancel(orderNo string,reason string,appId string) error {

	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}
	if order.OrderStatus==comm.ORDER_STATUS_SURED {
		return errors.New("订单已确认,不能取消!")
	}


	if order.PayStatus==comm.ORDER_PAY_STATUS_SUCCESS { //付款了的订单需要退款

		if time.Now().Unix() > order.UpdateTime.Unix() + 60*10 {

			return errors.New("订单已超过10分钟.不能取消!")
		}

		if order.Code=="" {
			return errors.New("订单不存在预付款code!")
		}

		tx,_ :=db.NewSession().Begin()
		defer func() {
			if err:=recover();err!=nil {
				tx.Rollback()
			}
		}()
		err = order.UpdateWithOrderStatusTx(comm.ORDER_STATUS_CANCELED_WAIT_SURE,orderNo,tx)
		if err!=nil{
			log.Error("更新订单状态失败! 订单号:",orderNo)
			tx.Rollback()
			return err
		}

		err :=order.UpdateWithCancelReasonTx(reason,orderNo,tx)
		if err!=nil{
			log.Error(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		err = order.UpdateWithOrderStatus(comm.ORDER_STATUS_CANCELED,orderNo)
		if err!=nil{
			log.Error("更新订单状态失败! 订单号:",orderNo)
			return err
		}
	}
	return nil
}

func OrdersGet(searchs interface{},pageIndex uint64,pageSize uint64,appId string)([]*dao.Order,error) {

	return dao.NewOrder().With(searchs,pageIndex,pageSize,appId)
}
func OrdersGetCount(searchs interface{},appId string)(int64,error) {

	return dao.NewOrder().WithCount(searchs,appId)
}


func orderSave(model *OrderModel,tx *dbr.Tx) (*dao.Order,error)  {

	merchant := dao.NewMerchant()
	merchant,err  :=merchant.MerchantWithId(model.MerchantId)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	if merchant==nil{
		log.Error("商户不存在!")
		return nil,errors.New("商户不存在!")
	}
	order := dao.NewOrder()
	order.Json = model.Json
	order.OpenId = model.OpenId
	order.No = NewInOrderNo()
	order.AppId = model.AppId
	order.Title = model.Title
	order.OrderStatus = model.OrderStatus
	order.PayStatus = model.PayStatus
	order.AddressId = model.AddressId
	order.MerchantId = model.MerchantId
	order.MOpenId = merchant.OpenId
	order.MerchantName = merchant.Name


	items := model.Items
	if items==nil || len(items)<=0 {
		return nil,errors.New("订单项不能为空!")
	}
	totalActPrice := 0.0
	totalPrice :=0.0
	for _,item :=range items  {
		prodSkuDetail := dao.NewProdSkuDetail()
		prodSkuDetail,err :=prodSkuDetail.WithSkuNo(item.SkuNo)
		if err!=nil{
			return nil,err
		}
		if prodSkuDetail==nil{
			return nil,errors.New("没有找到此商品")
		}

		if prodSkuDetail.Stock<=0 {
			return nil,errors.New("此商品已没有库存!")
		}
		totalActPrice+=prodSkuDetail.DisPrice*float64(item.Num)
		totalPrice += prodSkuDetail.Price*float64(item.Num)
		err =orderItemSave(prodSkuDetail,item,order.No,tx)
		if err!=nil{
			return nil,err
		}

	}
	order.RealPrice = totalActPrice
	order.Price = totalPrice

	orderId,err := order.InsertTx(tx)
	if err!=nil{
		return nil,err
	}
	order.Id = orderId

	return order,err
}

func orderItemSave(prodSkuDetail *dao.ProdSkuDetail,item OrderItemModel,orderNo string,tx *dbr.Tx) error  {
	orderItem :=dao.NewOrderItem()
	orderItem.No = orderNo
	orderItem.Title = prodSkuDetail.Title
	orderItem.ProdId = prodSkuDetail.ProdId
	orderItem.SkuNo = prodSkuDetail.SkuNo
	orderItem.DbnNo = item.DbnNo
	orderItem.AppId = prodSkuDetail.AppId
	orderItem.Num = item.Num
	orderItem.OfferUnitPrice = prodSkuDetail.Price
	orderItem.OfferTotalPrice = prodSkuDetail.Price*float64(item.Num)
	orderItem.BuyUnitPrice = prodSkuDetail.DisPrice
	orderItem.BuyTotalPrice = prodSkuDetail.DisPrice*float64(item.Num)
	orderItem.Json = item.Json
	return  orderItem.InsertTx(tx)
}

func OrderPrePayDtoToModel(dto OrderPrePayDto ) *OrderPrePayModel  {

	model :=&OrderPrePayModel{}
	model.AppId = dto.AppId
	model.OrderNo = dto.OrderNo
	model.PayType = dto.PayType
	model.AddressId = dto.AddressId
	model.CouponTokens = dto.CouponTokens
	model.Json = dto.Json
	model.NotifyUrl = config.GetValue("notify_url").ToString()
	return model
}
//第三方支付
func UpdateToPayed(orderNo string,appId string) error  {
	order :=dao.NewOrder()
	//order
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil {
		return  errors.New("paycall:订单查询失败!")
	}
	if order==nil{
		return  errors.New("paycall:没找到订单信息!")
	}
	if order.PayStatus!=comm.ORDER_PAY_STATUS_PAYING {
		return  errors.New("paycall:订单不是待付款状态!")
	}
	//支付预付款
	/*params := map[string]interface{}{
		"pay_token": payToken,
		"open_id": order.OpenId,
		"code": order.Code,
		"out_trade_pay":order.No,
	}
	 resultMap,err := RequestPayApi("/pay/payimprest",params)
	if err!=nil{
		return err
	}
	subTradeNo :=resultMap["sub_trade_no"].(string) */
	//tx
	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			log.Error(err)
			tx.Rollback()
		}
	}()
	//orderItem
	orderItem := dao.NewOrderItem()
	orderItems,err :=orderItem.OrderItemWithOrderNo(orderNo)
	if err!=nil{
		return  errors.New("paycall:查询订单明细失败!")
	}
	if orderItems==nil{
		return  errors.New("paycall:订单明细为空!")
	}
	//调整商品库存
	err = ProdSKUStockSubWithOrder(orderItems,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//商品累计销量增加
	err = ProdSoldNumAdd(orderItems,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//商户权重加1
	/* err =MerchantWeightAdd(1,order.MerchantId,tx)
	if err!=nil{
		tx.Rollback()
		return errors.New("paycall:商户权重添加失败!")
	} */
	//修改订单状态
	err =order.UpdateWithStatusTx(comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_SUCCESS,orderNo,tx)
	if err!=nil{
		tx.Rollback()
		return errors.New("paycall:订单号:状态更新为支付成功的时候失败!")
	}
	/* err =NotifyCouponServer(orderNo,appId,subTradeNo)
	if err!=nil{
		tx.Rollback()
		return errors.New("paycall:通知第三方优惠服务失败!")
	} */
	err =tx.Commit()
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return errors.New("paycall:数据提交失败!")
	}
	//发布事件
	go func(){
		err =PublishOrderPaidEvent(orderItems,order)
		if err!=nil{
			log.Warn("paycall:发送订单事件失败:", err)
		}
	}()
	//============
	return nil
}
//订单删除
func OrderDelete(orderNo string,appId string) error {

	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}
	
	if order.OrderStatus!=comm.ORDER_STATUS_CANCELED {
		return errors.New("订单未取消!")
	}
	
	log.Error("删除订单",orderNo)	
	err = order.OrderDelete(orderNo,appId)
	if err!=nil{		
		return err
	}
	
	return nil
}
//订单删除 批量
func OrderDeleteBatch(orderNo OrderNo,appId string) error {

	orderDao :=dao.NewOrder()
	var order *dao.Order
	var err	error
	for _,orderNo :=range orderNo.OrderNo {
		log.Error("删除订单",orderNo)
		
		order,err =orderDao.OrderWithNo(orderNo,appId)
		if err!=nil || order==nil{
			continue
		}
		
		if order.OrderStatus!=comm.ORDER_STATUS_CANCELED {
			continue
		}		
		err = orderDao.OrderDelete(orderNo,appId)
		if err!=nil{		
			continue
		}
	}
	
	return nil
}
//单增加购买订单号
func OrdersAddNum(orderNo string,appId string,ordernum string) error {
	/* order :=dao.NewOrder()
	order,err :=order.OrderWithId(orderId,appId)
	if err!=nil{
		return err
	} */
	
	orderItem := dao.NewOrderItem()
	err :=orderItem.OrdersAddNumWithNo(orderNo,appId,ordernum)
	if err!=nil{
		return err
	}
	
	return nil
}
//订单增加购买运单号
func OrdersAddPassnum(orderNo string,appId string,passnum string,passway string) error {
	/* order :=dao.NewOrder()
	order,err :=order.OrderWithId(orderId,appId)
	if err!=nil{
		return err
	} */
	
	orderItem := dao.NewOrderItem()
	err :=orderItem.OrdersAddNumWithPassnum(orderNo,appId,passnum,passway)
	if err!=nil{
		return err
	}
	
	return nil
}
func OrderItems(orderNo string) ([]*dao.OrderItem,error)  {
	orderItem := dao.NewOrderItem()
	orderItems,err :=orderItem.OrderItemWithOrderNo(orderNo)
	if err!=nil{
		return nil,err
	}
	return orderItems,nil
}
//订单快递查询
func ExpressDelivery(logisticCode string,shipperCode string) (string,error)  {	
	//header参数
	headers := map[string]string{
		"Content-Type"	: "application/x-www-form-urlencoded",
		"Host"			: "api.kdniao.cc",
	}
	
	//参数 
	params := map[string]string{
		"EBusinessID": "1269064",
		"DataType"	 : "2",
		"RequestType": "1002",
	}
	
	h := md5.New()
    h.Write([]byte("{'OrderCode':'','ShipperCode':'"+shipperCode+"','LogisticCode':'"+logisticCode+"'}84b6cbd8-605b-4d60-96c6-0622a8a4328b"))
    md5s:=fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))	
	base64s := base64.StdEncoding.EncodeToString([]byte(md5s))	
	
	params["RequestData"] = url.QueryEscape("{'OrderCode':'','ShipperCode':'"+shipperCode+"','LogisticCode':'"+logisticCode+"'}")
	
	params["DataSign"] 	  = url.QueryEscape(base64s)
	
	request:=""
	for k, v := range params {  
		request+=fmt.Sprintf("%s=%s&", k, v)  
    }
	
	paramData:= []byte("")
	
	response,err := network.Post("http://api.kdniao.cc/Ebusiness/EbusinessOrderHandle.aspx?"+request,paramData,headers)
	if err!=nil{
		return "",err
	}
	if response.StatusCode==http.StatusOK {
		str:=response.Body
		str =strings.Replace(str, "\n", "", -1)
		
		/* str:=strings.Replace(response.Body, "\n", "", -1)
		str =strings.Replace(str, "\r", "", -1)
		str =strings.Replace(str, "\\", "", -1)
	
		type Res struct {
			State string
		}
		var res Res 
		json.Unmarshal([]byte(str), &res)
	
		
		if res.State=="3" {
			fmt.Println("=@#$#%$%^$%")
		} */
		
		
		
		return str,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		return "",errors.New(resultMap["err_msg"].(string))
	}
	return "",errors.New("访问接口失败")	
}

//changeshowstate
func OrderChangeShowState(appId string,no string,show int64) error  {	
	return dao.NewOrder().OrderChangeShowState(appId,no,show)
}


















