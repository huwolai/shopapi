package service

import (
	"shopapi/dao"
	"errors"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
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
	AddressId int64
	PayType int
	AppId string
	NotifyUrl string
}

type OrderPrePayDto struct  {
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

	address := dao.NewAddress()
	address,err = address.WithId(model.AddressId)
	if address==nil{
		return nil,errors.New("没有找到对应的地址信息!")
	}

	if model.PayType == comm.Pay_Type_Account {//账户支付
		params := map[string]interface{}{
			"open_id":order.OpenId,
			"type": 1,
			"amount": int64(order.ActPrice*100),
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
		return resultImprestMap,nil
	}


	//参数
	params := map[string]interface{}{
		"open_id": order.OpenId,
		"out_trade_no": order.No,
		"amount": int64(order.ActPrice*100),
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
	}
	_,err = RequestPayApi("/pay/payimprest",params)
	if err!=nil{
		return err
	}

	err =order.UpdateWithStatus(comm.ORDER_STATUS_WAIT_SURE,comm.ORDER_PAY_STATUS_SUCCESS,orderNo)
	if err!=nil{
		log.Error("订单号:",orderNo,"状态更新为支付成功的时候失败!")
		return errors.New("订单更新错误!")
	}

	return nil
}

func OrderCancel(orderNo string,appId string) error {

	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil{
		return err
	}
	if order.OrderStatus==comm.ORDER_STATUS_SURED {
		return errors.New("订单已确认,不能取消!")
	}


	if order.OrderStatus==comm.ORDER_PAY_STATUS_SUCCESS { //付款了的订单需要退款
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
	}
	err = order.UpdateWithOrderStatus(comm.ORDER_STATUS_CANCELED,orderNo)
	if err!=nil{
		log.Error("更新订单状态失败! 订单号:",orderNo)
		return err
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
		totalActPrice+=prodSku.DisPrice*float64(item.Num)
		totalPrice += prodSku.Price*float64(item.Num)
		err =orderItemSave(prodSku,item,order.No,tx)
		if err!=nil{
			return nil,err
		}

	}
	order.ActPrice = totalActPrice
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
	model.NotifyUrl = config.GetValue("notify_url").ToString()
	return model
}