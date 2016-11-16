package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Order struct  {
	Id int64
	No string
	Code  string
	//预付款编号(主要针对第三方支付的)
	PrepayNo string
	PayapiNo string
	DbnAmount float64
	//商户应得金额
	MerchantAmount float64
	//优惠金额
	CouponAmount float64
	OpenId string
	MerchantId int64
	MerchantName string
	MOpenId string
	AppId string
	AddressId int64
	Address string
	AddressMobile string
	AddressName string
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
	
	GmOrdernum string
	GmPassnum string
	GmPassway string
	WayStatus int64
	
	DetailTitle []string
}

type OrderDetail struct  {
	Id int64
	No string
	PayapiNo string
	OpenId string
	AppId string
	AddressId int64
	Address string
	AddressMobile string
	AddressName string
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

	GmOrdernum string
	GmPassnum string
	GmPassway string
	WayStatus int64
}

type OrderCount struct  {
	Count int64
}

type OrderSearch struct {
	MerchantName	string
	Title  			string
	OrderNo  	 	string
	PayStatus 	 	uint64
	OrderStatus		uint64
	AddressMobile	string
}

func NewOrder() *Order {

	return &Order{}
}

func NewOrderDetail() *OrderDetail  {

	return &OrderDetail{}
}

func NewOrderCount() *OrderCount  {

	return &OrderCount{}
}

func (self *Order) OrderWithStatusLTTime(payStatus int,orderStatus int,time string) ([]*Order,error)  {
	var orders []*Order
	_,err :=db.NewSession().Select("*").From("`order`").Where("pay_status=?",payStatus).Where("order_status=?",orderStatus).Where("update_time<=?",time).LoadStructs(&orders)

	return orders,err

}

func (self *Order) With(searchs interface{},pageIndex uint64,pageSize uint64,appId string) ([]*Order,error)  {
	var orders []*Order
	//_,err :=db.NewSession().Select("*").From("`order`").Where("app_id=?",appId).Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("create_time",false).LoadStructs(&orders)

	buider :=db.NewSession().Select("*").From("`order`").Where("app_id=?",appId)
	
	search:=searchs.(OrderSearch)
	if search.MerchantName!="" {
		buider = buider.Where("merchant_name like ?","%"+search.MerchantName+"%")
	}
	if search.Title!="" {
		buider = buider.Where("title like ?","%"+search.Title+"%")
	}
	if search.OrderNo!="" {
		buider = buider.Where("no = ?",search.OrderNo)
	}
	if search.AddressMobile!="" {
		buider = buider.Where("address_mobile like ?",search.AddressMobile+"%")
	}
	switch search.PayStatus {
		case 1://1，未付款；
			buider = buider.Where("pay_status = ?",0)
		case 2://2，已付款
			buider = buider.Where("pay_status = ?",1)
		case 3://3，已付款
			buider = buider.Where("pay_status = ?",2)
	}
	switch search.OrderStatus {
		case 1://1，未确认
			buider = buider.Where("order_status = ?",0)
		case 2://2，已确认；
			buider = buider.Where("order_status = ?",1)
		case 3://3，已取消；
			buider = buider.Where("order_status = ?",2)
		case 4://4，无效；
			buider = buider.Where("order_status = ?",3)
		case 5://5，退货
			buider = buider.Where("order_status = ?",4)
	}
	
	_,err :=buider.Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("create_time",false).LoadStructs(&orders)
	return orders,err
}

func (self *Order) WithCount(searchs interface{},appId string) (int64,error)  {
	
	var count int64
	//err :=db.NewSession().Select("count(*)").From("`order`").Where("app_id=?",appId).LoadValue(&count)
	
	buider :=db.NewSession().Select("count(*)").From("`order`").Where("app_id=?",appId)
	
	search:=searchs.(OrderSearch)
	if search.MerchantName!="" {
		buider = buider.Where("merchant_name like ?","%"+search.MerchantName+"%")
	}
	if search.Title!="" {
		buider = buider.Where("title like ?","%"+search.Title+"%")
	}
	if search.OrderNo!="" {
		buider = buider.Where("no = ?",search.OrderNo)
	}
	switch search.PayStatus {
		case 1://1，未付款；
			buider = buider.Where("pay_status = ?",0)
		case 2://2，已付款
			buider = buider.Where("pay_status = ?",1)
		case 3://3，已付款
			buider = buider.Where("pay_status = ?",2)
	}
	switch search.OrderStatus {
		case 1://1，未确认
			buider = buider.Where("order_status = ?",0)
		case 2://2，已确认；
			buider = buider.Where("order_status = ?",1)
		case 3://3，已取消；
			buider = buider.Where("order_status = ?",2)
		case 4://4，无效；
			buider = buider.Where("order_status = ?",3)
		case 5://5，退货
			buider = buider.Where("order_status = ?",4)
	}

	err :=buider.LoadValue(&count)
	
	return count,err
}

//查询没有付款的订单 并且时间小于某个时间的
func (self *Order) OrderWithNoPayAndLTTime(time string) ([]*Order,error) {
	var orders []*Order
	_,err :=db.NewSession().Select("*").From("`order`").Where("pay_status=? or pay_status=?",comm.ORDER_PAY_STATUS_NOPAY,comm.ORDER_PAY_STATUS_PAYING).Where("update_time<=?",time).Where("order_status=?",comm.ORDER_STATUS_WAIT_SURE).LoadStructs(&orders)

	return orders,err
}

func (self *Order) InsertTx(tx *dbr.Tx) (int64,error)  {
	result,err :=tx.InsertInto("order").Columns("no","prepay_no","address_id","address","address_name","address_mobile","merchant_id","merchant_name","m_open_id","payapi_no","code","open_id","app_id","title","coupon_amount","dbn_amount","merchant_amount","real_price","pay_price","omit_money","price","order_status","pay_status","flag","json").Record(self).Exec()
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
func (self *Order) OrderWithId(id uint64,appId string) (*Order,error)  {

	sess := db.NewSession()
	var order *Order
	_,err :=sess.Select("*").From("`order`").Where("`id`=?",id).Where("app_id=?",appId).LoadStructs(&order)

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

	builder :=sess.Select("*").From("`order`").Where("merchant_id=?",merchantId).Where("app_id=?",appId)

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

func (self *OrderDetail) OrderDetailWithUser(openId string,orderStatus []int,payStatus []int,appId string,pageIndex uint64,pageSize uint64) ([]*OrderDetail,error)  {

	sess := db.NewSession()
	var orders []*OrderDetail

	builder :=sess.Select("*").From("`order`").Where("open_id=?",openId).Where("app_id=?",appId)

	if orderStatus!=nil&&len(orderStatus)>0{
		builder =builder.Where("order_status in ?",orderStatus)
	}

	if payStatus!=nil&&len(payStatus) >0 {
		builder =builder.Where("pay_status in ?",payStatus)
	}
	log.Error("==========================");
	log.Error(pageSize);
	log.Error(pageIndex);
	log.Error("==========================");
	_,err :=builder.Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("create_time",false).LoadStructs(&orders)
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

//获取用户指定状态订单数量
func (self *OrderCount) OrderWithUserAndStatusCount(openId string,orderStatus []int,payStatus []int,appId string) (*OrderCount,error)  {

	sess := db.NewSession()
	var orders *OrderCount

	builder :=sess.Select("count(id) as count").From("`order`").Where("open_id=?",openId).Where("app_id=?",appId)
	//builder :=sess.Select("count(id) as count").From("`order`")

	if orderStatus!=nil&&len(orderStatus)>0{
		builder =builder.Where("order_status in ?",orderStatus)
	}

	if payStatus!=nil&&len(payStatus) >0 {
		builder =builder.Where("pay_status in ?",payStatus)
	}
	
	//log.Error(builder.ToSql())
	
	_,err :=builder.LoadStructs(&orders)
	if err!=nil{
		return nil,err
	}
	if orders==nil{
		return nil,nil
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
			if len(orderItemDetailMap)>0{
				order.GmOrdernum=orderItemDetailMap[order.No][0].GmOrdernum
				order.GmPassnum=orderItemDetailMap[order.No][0].GmPassnum
				order.GmPassway=orderItemDetailMap[order.No][0].GmPassway
				order.WayStatus=orderItemDetailMap[order.No][0].WayStatus
			}			
		}		
	}

	return nil
}

func (self *Order) OrderPayapiUpdateWithNoAndCodeTx(payapiNo string,addressId int64,address string,addressName string,addressMobile string,code string,orderStatus int,payStatus int,no string,json string,appId string,tx *dbr.Tx) error  {
	builder :=tx.Update("order").Set("payapi_no",payapiNo).Set("address_id",addressId).Set("address_name",addressName).Set("address",address).Set("address_mobile",addressMobile).Set("code",code).Set("order_status",orderStatus).Set("pay_status",payStatus)
	if json!="" {
		builder = builder.Set("json",json)
	}
	_,err := builder.Where("app_id=?",appId).Where("`no`=?",no).Exec()
	return err
}

func (self *Order) UpdateWithStatus(orderStatus int,payStatus int,orderNo string) error {

	_,err :=db.NewSession().Update("order").Set("order_status",orderStatus).Set("pay_status",payStatus).Where("no=?",orderNo).Exec()

	return err
}

func (self *Order) UpdateWithStatusTx(orderStatus int,payStatus int,orderNo string,tx *dbr.Tx) error {
	log.Error("payStatus:",payStatus)
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
//订单删除
func (self *Order) OrderDelete(orderNo string,appId string) error {
	sesson := db.NewSession()
	tx,_  :=sesson.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	
	//_,err :=db.NewSession().DeleteBySql("delete from order where `no`=? and app_id=? limit 1",orderNo,appId).Exec()
	_,err :=tx.DeleteFrom("order").Where("no=?",orderNo).Where("app_id=?",appId).Exec()
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	//_,err =db.NewSession().DeleteBySql("delete from order_item where `no`=? and app_id=?",orderNo,appId).Exec()
	_,err =tx.DeleteFrom("order_item").Where("no=?",orderNo).Where("app_id=?",appId).Exec()
	if err!=nil{
		tx.Rollback()
		return err
	}
	
	if err :=tx.Commit();err!=nil{
		tx.Rollback()

		return err
	}
	
	return nil
}


























