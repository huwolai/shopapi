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
	OpenId string
	AppId string
	Title string
	Status int
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
	PayType int
	AppId string
	NotifyUrl string
}

type OrderPrePayDto struct  {
	OrderNo string `json:"order_no"`
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
		err =order.OrderPayapiUpdateWithNoAndCode("",code,comm.ORDER_STATUS_PAY_WAIT,order.No,order.AppId)
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
	log.Info(params)

	resultPrepayMap,err :=RequestPayApi("/pay/makeprepay",params)
	if err!=nil{
		return nil,err
	}
	if resultPrepayMap!=nil{
		payapiNo :=resultPrepayMap["pay_no"].(string)
		code :=resultPrepayMap["code"].(string)
		//将payapi的订单号更新到订单数据里
		err :=order.OrderPayapiUpdateWithNoAndCode(payapiNo,code,comm.ORDER_STATUS_PAY_WAIT,order.No,order.AppId)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	return resultPrepayMap,nil

}

func OrderByUser(openId string,status []int,appId string)  ([]*dao.OrderDetail,error)  {

	orderDetail :=dao.NewOrderDetail()
	orderDetails,err := orderDetail.OrderDetailWithUser(openId,status,appId)

	return orderDetails,err
}

func OrderDetailWithNo(orderNo string,appId string) (*dao.OrderDetail,error)  {
	orderDetail :=dao.NewOrderDetail()
	orderDetail,err := orderDetail.OrderDetailWithNo(orderNo,appId)
	return orderDetail,err
}

func OrderPayForAccount(openId string,orderNo string,appId string) error  {

	order :=dao.NewOrder()
	order,err :=order.OrderWithNo(orderNo,appId)
	if err!=nil {
		log.Error(err)
		return errors.New("订单查询失败!")
	}

	if order==nil{
		return  errors.New("没找到订单信息!")
	}
	if order.Status!=comm.ORDER_STATUS_PAY_WAIT {
		return  errors.New("订单不是待付款状态!")
	}

	if order.Code==""{
		return  errors.New("订单没有预付款代号,订单数据有误!")
	}
	account :=dao.NewAccount()
	account,err =account.AccountWithOpenId(openId,appId)
	if err!=nil {
		return err
	}
	if account==nil{
		return errors.New("没有找到用户的账户信息!请重新登录再试")
	}



	//获取支付token
	params := map[string]interface{}{
		"open_id": order.OpenId,
		"password": account.Password,
	}
	resultPayTokenMap,err := RequestPayApi("/pay/token",params)
	if err!=nil{
		return err
	}
	paytoken :=resultPayTokenMap["token"].(string)

	//支付预付款
	params = map[string]interface{}{
		"pay_token": paytoken,
		"open_id": order.OpenId,
		"code": order.Code,
	}
	_,err = RequestPayApi("/pay/payimprest",params)
	if err!=nil{
		return err
	}

	err =order.UpdateWithStatus(comm.ORDER_STATUS_PAY_SUCCESS,orderNo)
	if err!=nil{
		log.Error("订单号:",orderNo,"状态更新为支付成功的时候失败!")
		return errors.New("订单更新错误!")
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
	order.Status = model.Status
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
		log.Error("-----prodSku=",prodSku.SkuNo)
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
	model.NotifyUrl = config.GetValue("notify_url").ToString()
	return model
}