package service

import (
	"shopapi/dao"
	"errors"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
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
	NotifyUrl string
}

type OrderPrePayDto struct  {
	//优惠token
	CouponTokens []string `json:"coupon_tokens"`
	OrderNo string `json:"order_no"`
	AddressId int64 `json:"address_id"`
	//付款类型(1.支付宝 2.微信 3.现金支付 4.账户)
	PayType int `json:"pay_type"`
	AppId string `json:"app_id"`
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
	if address==nil{
		return nil,errors.New("没有找到对应的地址信息!")
	}

	var couponTotalAmount float64
	//存在优惠信息
	if model.CouponTokens!=nil&&len(model.CouponTokens) >0 {
		couponTotalAmount,err =HandleCoupon(order,model.CouponTokens)
		if err!=nil{
			log.Error(err)
			return nil,errors.New("优惠信息处理错误!")
		}
	}
	payPrice :=order.RealPrice - couponTotalAmount
	if payPrice<0 {
		return nil,errors.New("付款金额不能为负数!")
	}

	if model.PayType == comm.Pay_Type_Account {//账户支付

		if order.PayStatus==comm.ORDER_PAY_STATUS_NOPAY{
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
			code :=resultImprestMap["code"].(string)
			err =order.OrderPayapiUpdateWithNoAndCode("",address.Id,address.Address,code,comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_PAYING,order.No,order.AppId)
			if err!=nil{
				log.Error(err)
				return nil,err
			}

			resultMap := map[string]interface{}{
				"open_id": resultImprestMap["open_id"],
				"code": resultImprestMap["code"],
			}
			return resultMap,nil
		}else {
			resultMap := map[string]interface{}{
				"open_id": order.OpenId,
				"code": order.Code,
			}
			return resultMap,nil
		}
	}
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


	if err!=nil{
		return nil,err
	}

	resultPrepayMap,err :=RequestPayApi("/pay/makeprepay",params)
	if err!=nil{
		return nil,err
	}
	if resultPrepayMap!=nil{
		payapiNo :=resultPrepayMap["pay_no"].(string)
		code :=resultPrepayMap["code"].(string)
		//将payapi的订单号更新到订单数据里
		err :=order.OrderPayapiUpdateWithNoAndCode(payapiNo,address.Id,address.Address,code,comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_PAYING,order.No,order.AppId)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	return resultPrepayMap,nil
}

//处理优惠信息
func HandleCoupon(order *dao.Order,coupotokens []string) (float64,error)  {
	tx,_ :=db.NewSession().Begin()
	orderCoupon := dao.NewOrderCoupon()
	err :=orderCoupon.DeleteWithOrderNoTx(comm.ORDER_COUPON_STATUS_UNACTIVATE,order.No,tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return 0.0,errors.New("删除原有优惠记录失败!")
	}
	//去重
	ncoupontokens := comm.RemoveDuplicatesAndEmpty(coupotokens)
	//凭证校验
	var couponTotalAmount float64
	for _,couponToken :=range ncoupontokens {
		jwtAuth := comm.InitJWTAuthenticationBackend()
		cpToken,err :=jwtAuth.FetchToken(couponToken)
		if err!=nil{
			tx.Rollback()
			return 0.0,err
		}
		if !cpToken.Valid {
			tx.Rollback()
			return 0.0,errors.New("优惠凭证无效!")
		}
		orderNo,isok :=cpToken.Claims["order_no"].(string)
		if !isok {
			tx.Rollback()
			return 0.0,errors.New("优惠券有误[获取order_no失败]!")
		}
		if orderNo!=order.No {
			return 0.0,errors.New("优惠凭证不是当前订单的!")
		}

		couponNo,isok :=cpToken.Claims["coupon_no"].(string)
		if !isok{
			tx.Rollback()
			return 0.0,errors.New("优惠券有误[获取coupon_no失败]!")
		}
		couponAmount,isok :=cpToken.Claims["coupon_amount"].(string)
		if !isok {
			tx.Rollback()
			return 0.0,errors.New("优惠券有误[获取coupon_amount失败]!")
		}
		notifyUrl,isok :=cpToken.Claims["notify_url"].(string)
		if !isok {
			tx.Rollback()
			return 0.0,errors.New("优惠券有误[获取notify_url失败]!")
		}
		trackCode,isok :=cpToken.Claims["track_code"].(string)
		if !isok {
			tx.Rollback()
			return 0.0,errors.New("优惠券有误[获取track_code失败]!")
		}
		fcouponAmount,err :=strconv.ParseFloat(couponAmount,10)
		if err!=nil{
			log.Error(err)
			tx.Rollback()
			return 0.0,errors.New("优惠券金额有误!")
		}

		orderCoupon := dao.NewOrderCoupon()
		orderCoupon.CouponAmount = fcouponAmount
		orderCoupon.CouponCode = couponNo
		orderCoupon.OpenId = order.OpenId
		orderCoupon.TrackCode = trackCode
		orderCoupon.CouponToken = couponToken
		orderCoupon.OrderNo = orderNo
		orderCoupon.AppId = order.AppId
		orderCoupon.NotifyUrl = notifyUrl
		err =orderCoupon.InsertTx(tx)
		if err!=nil{
			tx.Rollback()
			return 0.0,errors.New("插入优惠信息失败!")
		}
		couponTotalAmount += orderCoupon.CouponAmount
	}

	return couponTotalAmount,nil
}

func OrderByUser(openId string,orderStatus []int,payStatus []int,appId string)  ([]*dao.OrderDetail,error)  {

	orderDetail :=dao.NewOrderDetail()
	orderDetails,err := orderDetail.OrderDetailWithUser(openId,orderStatus,payStatus,appId)

	return orderDetails,err
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
	return orderDetail,err
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
	//支付预付款
	params := map[string]interface{}{
		"pay_token": payToken,
		"open_id": order.OpenId,
		"code": order.Code,
		"out_trade_pay":order.No,
	}
	resultMap,err := RequestPayApi("/pay/payimprest",params)
	if err!=nil{
		return err
	}

	subTradeNo :=resultMap["sub_trade_no"].(string)

	tx,_ :=db.NewSession().Begin()

	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
		}
	}()

	//调整商品库存
	err = ProdSKUStockSubWithOrder(orderNo,tx)
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

	err =NotifyCouponServer(orderNo,appId,subTradeNo)
	if err!=nil{
		tx.Rollback()
		return errors.New("通知第三方优惠服务失败!")
	}

	tx.Commit()

	return nil
}

//通知优惠服务
func NotifyCouponServer(orderNo string,appId string,subTradeNo string) error {

	orderCoupon :=dao.NewOrderCoupon()
	orderCoupons,err := orderCoupon.WithOrderNo(orderNo,appId)
	if err!=nil{
		return err
	}
	if orderCoupons==nil{
		return nil
	}

	for _,ordercn :=range orderCoupons {
		if orderCoupon.NotifyUrl=="" {
			continue
		}
		requestModel :=queue.NewRequestModel()
		requestModel.NotifyUrl = orderCoupon.NotifyUrl
		requestModel.Data = map[string]interface{}{
			"coupon_code": ordercn.CouponCode,
			"track_code": ordercn.TrackCode,
			"open_id": ordercn.OpenId,
			"sub_trade_no": subTradeNo,
		}
		err = queue.PublishRequestMsg(requestModel)
		if err!=nil{
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

//减商品sku的库存
func ProdSKUStockSubWithOrder(orderNo string,tx *dbr.Tx) error  {
	orderItem := dao.NewOrderItem()
	orderItems,err :=orderItem.OrderItemWithOrderNo(orderNo)
	if err!=nil{
		return  errors.New("查询订单明细失败!")
	}
	if orderItems!=nil&&len(orderItems)>0{
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
	}
	return nil
}

func OrderAutoCancel(orderNo string,appId string)error  {
	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}

	if order.PayStatus!=comm.ORDER_PAY_STATUS_NOPAY ||
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
		if order.Code=="" {
			return errors.New("订单不存在预付款code!")
		}

		tx,_ :=db.NewSession().Begin()
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



func orderSave(model *OrderModel,tx *dbr.Tx) (*dao.Order,error)  {

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
	order.MOpenId = model.MOpenId


	items := model.Items
	if items==nil || len(items)<=0 {
		return nil,errors.New("订单项不能为空!")
	}
	totalActPrice := 0.0
	totalPrice :=0.0
	for _,item :=range items  {
		prodSku := dao.NewProdSku()
		prodSku,err :=prodSku.WithSkuNo(item.SkuNo)
		if err!=nil{
			return nil,err
		}
		if prodSku==nil{
			return nil,errors.New("没有找到此商品")
		}

		if prodSku.Stock<=0 {
			return nil,errors.New("此商品已没有库存!")
		}
		totalActPrice+=prodSku.DisPrice*float64(item.Num)
		totalPrice += prodSku.Price*float64(item.Num)
		err =orderItemSave(prodSku,item,order.No,tx)
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

func orderItemSave(prodSku *dao.ProdSku,item OrderItemModel,orderNo string,tx *dbr.Tx) error  {
	orderItem :=dao.NewOrderItem()
	orderItem.No = orderNo
	orderItem.ProdId = prodSku.ProdId
	orderItem.SkuNo = prodSku.SkuNo
	orderItem.AppId = prodSku.AppId
	orderItem.Num = item.Num
	orderItem.OfferUnitPrice = prodSku.Price
	orderItem.OfferTotalPrice = prodSku.Price*float64(item.Num)
	orderItem.BuyUnitPrice = prodSku.DisPrice
	orderItem.BuyTotalPrice = prodSku.DisPrice*float64(item.Num)
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
	model.NotifyUrl = config.GetValue("notify_url").ToString()
	return model
}
