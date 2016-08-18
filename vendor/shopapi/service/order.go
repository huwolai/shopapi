package service

import (
	"shopapi/dao"
	"errors"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"encoding/json"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
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
	//商品ID
	ProdId int64
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

	order := dao.NewOrder()
	order,err :=order.OrderWithNo(model.OrderNo,model.AppId)
	if err!=nil {
		return nil,err
	}
	//参数
	params := map[string]interface{}{
		"open_id": order.OpenId,
		"out_trade_no": order.No,
		"amount": int64(order.ActPrice*100),
		"trade_type": 1,  //交易类型(1.充值 2.购买)
		"pay_type": model.PayType,
		"title": order.Title,
		"client_ip": "127.0.0.1",
		"notify_url": model.NotifyUrl,
		"remark": order.Title,
	}
	log.Info(params)

	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetPayapiSign(params)
	log.Info(fmt.Sprintf("%s.%s",basesign,sign))
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(params);

	response,err := network.Post(config.GetValue("payapi_url").ToString()+"/pay/makeprepay",paramData,headers)
	if err!=nil{
		return nil,err
	}

	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		if resultMap!=nil{
			payapiNo :=resultMap["pay_no"].(string)
			//将payapi的订单号更新到订单数据里
			err :=order.OrderPayapiUpdateWithNo(payapiNo,comm.ORDER_STATUS_PAY_WAIT,order.No,order.AppId)
			if err!=nil{
				log.Error(err)
				return nil,err
			}
		}

		return resultMap,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}else{
		return nil,errors.New("请求支付中心失败!")
	}

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
		product :=dao.NewProduct()
		product,err :=product.ProductWithId(item.ProdId,model.AppId)
		if err!=nil{
			return nil,err
		}
		totalActPrice+=product.DisPrice*float64(item.Num)
		totalPrice += product.Price*float64(item.Num)

		err =orderItemSave(product,item,order.No,tx)
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

func orderItemSave(product *dao.Product,item OrderItemModel,orderNo string,tx *dbr.Tx) error  {
	orderItem :=dao.NewOrderItem()
	orderItem.No = orderNo
	orderItem.ProdId = product.Id
	orderItem.AppId = product.AppId
	orderItem.Num = item.Num
	orderItem.OfferUnitPrice = product.Price
	orderItem.OfferTotalPrice = product.Price*float64(item.Num)
	orderItem.BuyUnitPrice = product.DisPrice
	orderItem.BuyTotalPrice = product.DisPrice*float64(item.Num)
	orderItem.Json = item.Json
	return  orderItem.InsertTx(tx)
}