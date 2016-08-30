package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
)

type Order struct  {
	Id int64
	No string
	Code  string
	//预付款编号(主要针对第三方支付的)
	PrepayNo string
	PayapiNo string
	DbnAmount float64
	MerchantAmount float64
	CouponAmount float64
	OpenId string
	MerchantId int64
	MOpenId string
	AppId string
	AddressId int64
	Address string
	Title string
	Price float64
	RealPrice float64
	PayPrice float64
	OmitMoney float64
	Flag string
	RejectCancelReason string
	CancelReason string
	OrderStatus int
	PayStatus int
	Json string
	BaseDModel
}

type OrderDetail struct  {
	Id int64
	No string
	PayapiNo string
	OpenId string
	AppId string
	Name *string
	Mobile *string
	AddressId int64
	Address string
	Title string
	DbnAmount float64
	MerchantAmount float64
	CouponAmount float64
	Price float64
	RealPrice float64
	PayPrice float64
	OmitMoney float64
	RejectCancelReason string
	CancelReason string
	OrderStatus int
	PayStatus int
	Items []*OrderItemDetail
	Flag string
	Json string
	BaseDModel
}

func NewOrder() *Order {

	return &Order{}
}

func NewOrderDetail() *OrderDetail  {

	return &OrderDetail{}
}

func (self *Order) OrderWithStatusLTTime(payStatus int,orderStatus int,time string) ([]*Order,error)  {
	var orders []*Order
	_,err :=db.NewSession().Select("*").From("`order`").Where("pay_status=?",payStatus).Where("order_status=?",orderStatus).Where("update_time<=?",time).LoadStructs(&orders)

	return orders,err

}

//查询没有付款的订单 并且时间小于某个时间的
func (self *Order) OrderWithNoPayAndLTTime(time string) ([]*Order,error) {
	var orders []*Order
	_,err :=db.NewSession().Select("*").From("`order`").Where("pay_status=? or pay_status=?",comm.ORDER_PAY_STATUS_NOPAY,comm.ORDER_PAY_STATUS_PAYING).Where("update_time<=?",time).LoadStructs(&orders)

	return orders,err
}

func (self *Order) InsertTx(tx *dbr.Tx) (int64,error)  {
	result,err :=tx.InsertInto("order").Columns("no","prepay_no","address_id","address","merchant_id","m_open_id","payapi_no","code","open_id","app_id","title","coupon_amount","dbn_amount","merchant_amount","real_price","pay_price","omit_money","price","order_status","pay_status","flag","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}

	lastId,err :=result.LastInsertId()

	return lastId,err
}

func (self *Order) OrderWithNo(no string,appId string) (*Order,error)  {

	sess := db.NewSession()
	var order *Order
	_,err :=sess.Select("*").From("`order`").Where("`no`=?",no).Where("app_id=?",appId).LoadStructs(&order)

	return order,err
}

func (self *OrderDetail) OrderDetailWithNo(no string,appId string) (*OrderDetail,error)  {
	sess := db.NewSession()
	var orders []*OrderDetail
	_,err :=sess.Select("*").From("`order`").Where("app_id=?",appId).Where("no=?",no).LoadStructs(&orders)
	if err!=nil {
		return nil,err
	}
	if len(orders)<=0 {

		return nil,nil
	}
	err =fillOrderItemDetail(orders)
	if err!=nil {
		return nil,err
	}

	return orders[0],err
}

func (self *OrderDetail) OrderDetailWithMerchantId(merchantId int64,orderStatus []int,payStatus []int,appId string) ([]*OrderDetail,error) {
	sess := db.NewSession()
	var orders []*OrderDetail

	builder :=sess.Select("`order`.*","address.name","address.mobile").From("`order`").LeftJoin("address","`order`.address_id=address.id").Where("`order`.merchant_id=?",merchantId).Where("`order`.app_id=?",appId)

	if orderStatus!=nil&&len(orderStatus)>0{
		builder =builder.Where("order_status in ?",orderStatus)
	}

	if payStatus!=nil&&len(payStatus) >0 {
		builder =builder.Where("pay_status in ?",payStatus)
	}
	_,err :=builder.OrderDir("create_time",false).LoadStructs(&orders)
	if err!=nil{

		return nil,err
	}
	if orders==nil{
		return nil,nil
	}

	err = fillOrderItemDetail(orders)
	if err!=nil {
		return nil,err
	}


	return orders,err
}

func (self *OrderDetail) OrderDetailWithUser(openId string,orderStatus []int,payStatus []int,appId string) ([]*OrderDetail,error)  {

	sess := db.NewSession()
	var orders []*OrderDetail

	builder :=sess.Select("*").From("`order`").Where("open_id=?",openId).Where("app_id=?",appId)

	if orderStatus!=nil&&len(orderStatus)>0{
		builder =builder.Where("order_status in ?",orderStatus)
	}

	if payStatus!=nil&&len(payStatus) >0 {
		builder =builder.Where("pay_status in ?",payStatus)
	}
	_,err :=builder.OrderDir("create_time",false).LoadStructs(&orders)
	if err!=nil{

		return nil,err
	}
	if orders==nil{
		return nil,nil
	}

	err = fillOrderItemDetail(orders)
	if err!=nil {
		return nil,err
	}


	return orders,err
}

func fillOrderItemDetail(orders []*OrderDetail)  error {
	ordernos :=make([]string,0)
	for _,orderDetail :=range orders {
		ordernos = append(ordernos,orderDetail.No)
	}

	if len(ordernos)>0 {
		orderItemDetail := NewOrderItemDetail()
		orderItemDetails,err :=orderItemDetail.OrderItemDetailWithOrderNo(ordernos)
		if err!=nil{
			return err
		}

		orderItemDetailMap :=make(map[string][]*OrderItemDetail)
		if len(orderItemDetails)>0 {
			for _,orderItemDetail :=range orderItemDetails {
				odDetailList := orderItemDetailMap[orderItemDetail.No]
				if odDetailList==nil {
					odDetailList = make([]*OrderItemDetail,0)
				}
				odDetailList = append(odDetailList,orderItemDetail)
				orderItemDetailMap[orderItemDetail.No] = odDetailList
			}
		}

		for _,order :=range orders {
			order.Items = orderItemDetailMap[order.No]
		}

	}

	return nil
}

func (self *Order) OrderPayapiUpdateWithNoAndCodeTx(payapiNo string,addressId int64,address string,code string,orderStatus int,payStatus int,no string,appId string,tx *dbr.Tx) error  {
	_,err :=tx.Update("order").Set("payapi_no",payapiNo).Set("address_id",addressId).Set("address",address).Set("code",code).Set("order_status",orderStatus).Set("pay_status",payStatus).Where("app_id=?",appId).Where("`no`=?",no).Exec()
	return err
}

func (self *Order) UpdateWithStatus(orderStatus int,payStatus int,orderNo string) error {

	_,err :=db.NewSession().Update("order").Set("order_status",orderStatus).Set("pay_status",payStatus).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithStatusTx(orderStatus int,payStatus int,orderNo string,tx *dbr.Tx) error {

	_,err :=tx.Update("order").Set("order_status",orderStatus).Set("pay_status",payStatus).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithOrderStatus(orderStatus int,orderNo string) error  {

	_,err :=db.NewSession().Update("order").Set("order_status",orderStatus).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithOrderStatusTx(orderStatus int,orderNo string,tx *dbr.Tx) error  {

	_,err :=tx.Update("order").Set("order_status",orderStatus).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateOrderPayInfoWithOrderNoTX(couponAmount float64,payPrice float64,orderNo string,appId string,tx *dbr.Tx) error  {
	_,err :=tx.Update("order").Set("coupon_amount",couponAmount).Set("pay_price",payPrice).Where("app_id=?",appId).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithRefuseCancelReasonTx(refuseCancelReason string,orderNo string,tx *dbr.Tx) error  {
	_,err :=tx.Update("order").Set("reject_cancel_reason",refuseCancelReason).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithCancelReasonTx(cancelReason string,orderNo string,tx *dbr.Tx) error  {
	_,err :=tx.Update("order").Set("cancel_reason",cancelReason).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateAmountTx(couponAmount,payPrice,merchantAmount,omitMoney,dbnAmount float64,orderNo string,tx *dbr.Tx) error  {
	_,err :=tx.Update("order").Set("coupon_amount",couponAmount).Set("pay_price",payPrice).Set("merchant_amount",merchantAmount).Set("omit_money",omitMoney).Set("dbn_amount",dbnAmount).Where("no=?",orderNo).Exec()

	return err
}