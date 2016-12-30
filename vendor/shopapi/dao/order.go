package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Order struct  {
	Id 					int64		`json:"id"`	
	No 					string		`json:"no"`
	Code  				string		`json:"code"`
	//预付款编号(主要针对第三方支付的)
	PrepayNo 			string		`json:"prepay_no"`
	PayapiNo 			string		`json:"payapi_no"`
	DbnAmount 			float64		`json:"dbn_amount"`
	//商户应得金额
	MerchantAmount		float64		`json:"merchant_amount"`
	//优惠金额	
	CouponAmount 		float64		`json:"coupon_amount"`
	OpenId 				string		`json:"open_id"`
	MerchantId 			int64		`json:"merchant_id"`
	MerchantName 		string		`json:"merchant_name"`
	MOpenId 			string		`json:"mopen_id"`
	AppId 				string		`json:"app_id"`
	AddressId 			int64		`json:"address_id"`
	Address 			string		`json:"address"`
	AddressMobile 		string		`json:"address_mobile"`
	AddressName 		string		`json:"address_name"`
	Title 				string		`json:"title"`
	Price 				float64		`json:"price"`
	RealPrice 			float64		`json:"real_price"`
	PayPrice 			float64		`json:"pay_price"`
	OmitMoney 			float64		`json:"omit_money"`
	Flag 				string		`json:"flag"`
	RejectCancelReason 	string		`json:"reject_cancel_reason"`
	CancelReason 		string		`json:"cancel_reason"`
	OrderStatus 		int			`json:"order_status"`
	PayStatus 			int			`json:"pay_status"`
	Json 				string		`json:"json"`
	BaseDModel	
	
	GmOrdernum 			string		`json:"gm_ordernum"`
	GmPassnum 			string		`json:"gm_passnum"`
	GmPassway 			string		`json:"gm_passway"`
	WayStatus 			int64		`json:"way_status"`
	
	DetailTitle 		[]string	`json:"detail_title"`
	
	UpdateTimeUnix 		int64		`json:"update_time_unix"`
	
	Show 				int			`json:"show"`
	Mobile 				string		`json:"mobile"`
	YdgyName  			string		`json:"tdgy_name"`
	OrderType	  		string		`json:"order_type"`
	
	ProdId 				int64		`json:"prod_id"`
	SkuNo 				string		`json:"sku_no"`
}

type OrderYyg struct  {
	BuyCode string	`json:"buy_code"`
	YdgyName string	`json:"ydgy_name"`
	Order
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
	
	Itemsyyg *ProdPurchaseCode
	
	ItemsyygBuyCodes []string
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
	Show			uint64
	OrderType		[]string
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
	_,err :=db.NewSession().Select("*,UNIX_TIMESTAMP(create_time) as update_time_unix").From("`order`").Where("pay_status=?",payStatus).Where("order_status=?",orderStatus).Where("create_time<=?",time).LoadStructs(&orders)
	//_,err :=db.NewSession().Select("*,UNIX_TIMESTAMP(update_time) as update_time_unix").From("`order`").Where("pay_status=?",payStatus).Where("order_status=?",orderStatus).Where("update_time<=?",time).LoadStructs(&orders)

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
		//buider = buider.Where("address_mobile like ?",search.AddressMobile+"%")
		buider = buider.Where("open_id in ( select open_id from account where mobile like ?)",search.AddressMobile+"%")
	}
	if len(search.OrderType)>0 {
		buider = buider.Where("order_type in ?",search.OrderType)
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
	
	//show
	if search.Show>0 {
		if search.Show==1 {
			buider = buider.Where("`show` = ?",1)
		}else{
			buider = buider.Where("`show` = ?",0)
		}		
	}	
	//log.Error(buider.ToSql())
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
	if search.AddressMobile!="" {
		//buider = buider.Where("address_mobile like ?",search.AddressMobile+"%")
		buider = buider.Where("open_id in ( select open_id from account where mobile like ?)",search.AddressMobile+"%")
	}
	if len(search.OrderType)>0 {
		buider = buider.Where("order_type in ?",search.OrderType)
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
	
	//show
	if search.Show>0 {
		if search.Show==1 {
			buider = buider.Where("`show` = ?",1)
		}else{
			buider = buider.Where("`show` = ?",0)
		}		
	}	
	
	err :=buider.LoadValue(&count)
	
	return count,err
}

//查询没有付款的订单 并且时间小于某个时间的
func (self *Order) OrderWithNoPayAndLTTime(time string) ([]*Order,error) {
	var orders []*Order
	_,err :=db.NewSession().Select("*, UNIX_TIMESTAMP(create_time) as update_time_unix").From("`order`").Where("pay_status=? or pay_status=?",comm.ORDER_PAY_STATUS_NOPAY,comm.ORDER_PAY_STATUS_PAYING).Where("create_time<=?",time).Where("order_status=?",comm.ORDER_STATUS_WAIT_SURE).LoadStructs(&orders)
	//_,err :=db.NewSession().Select("*, UNIX_TIMESTAMP(update_time) as update_time_unix").From("`order`").Where("pay_status=? or pay_status=?",comm.ORDER_PAY_STATUS_NOPAY,comm.ORDER_PAY_STATUS_PAYING).Where("update_time<=?",time).Where("order_status=?",comm.ORDER_STATUS_WAIT_SURE).LoadStructs(&orders)

	return orders,err
}

func (self *Order) InsertTx(tx *dbr.Tx) (int64,error)  {
	result,err :=tx.InsertInto("order").Columns("no","prepay_no","address_id","address","address_name","address_mobile","merchant_id","merchant_name","m_open_id","payapi_no","code","open_id","app_id","title","coupon_amount","dbn_amount","merchant_amount","real_price","pay_price","omit_money","price","order_status","pay_status","flag","json","order_type").Record(self).Exec()
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
				if(len(orderItemDetailMap[order.No])>0){
					order.GmOrdernum=orderItemDetailMap[order.No][0].GmOrdernum
					order.GmPassnum=orderItemDetailMap[order.No][0].GmPassnum
					order.GmPassway=orderItemDetailMap[order.No][0].GmPassway
					order.WayStatus=orderItemDetailMap[order.No][0].WayStatus
				}
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
/* func (self *Order) OrderTypes(orderNo string,appId string)([]*OrderItem,error)  {
	return NewOrderItem().OrderItemWithOrderNo(orderNo)
} */
//changeshowstate
func (self *Order) OrderChangeShowState(appId string,no string,show int64) error  {	
	_,err :=db.NewSession().Update("order").Set("show",show).Where("no=?",no).Where("app_id=?",appId).Exec()
	return err
}
func (self *Order) OrderWithPordYyg(OpenCode string,prodId int64) (*Order,error)  {
	var orders *Order

	builder :=db.NewSession().Select("*").From("`order`")
	
	builder =builder.Where("`no` in (select `no` from order_item_purchase_codes where codes =? and prod_id=?)",OpenCode,prodId)
	
	_,err :=builder.LoadStructs(&orders)
	return orders,err
}
func (self *Order) OrdersWithPordYyg(prodId string,pIndex uint64,pSize uint64) ([]*OrderYyg,int64,error)  {
	var orders []*OrderYyg	
	
	count, _ := db.NewSession().SelectBySql("select count(a.codes) from `order_item_purchase_codes` a left join `order` b on a.`no`=b.`no` where  a.`prod_id` = ?",prodId).ReturnInt64()
	
	_,err :=db.NewSession().SelectBySql("select c.ydgy_name,a.codes as buy_code,b.* from `order_item_purchase_codes` a left join `order` b on a.`no`=b.`no` left join `account` c on c.`open_id`=b.`open_id` where  a.`prod_id` = ? limit ?,?",prodId,(pIndex-1)*pSize,pSize).LoadStructs(&orders)
		
	return orders,count,err
}


func (self *Order) WithItems(searchs interface{},pageIndex uint64,pageSize uint64,appId string) ([]*Order,error)  {
	var orders []*Order
	//_,err :=db.NewSession().Select("*").From("`order`").Where("app_id=?",appId).Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("create_time",false).LoadStructs(&orders)

	buider :=db.NewSession().Select("`order`.*,order_item.prod_id,order_item.sku_no").From("`order_item`").LeftJoin("order","`order`.`no`=order_item.`no`").Where("`order_item`.app_id=?",appId)
	
	search:=searchs.(OrderSearch)
	if search.MerchantName!="" {
		buider = buider.Where("`order`.merchant_name like ?","%"+search.MerchantName+"%")
	}
	if search.Title!="" {
		buider = buider.Where("`order`.title like ?","%"+search.Title+"%")
	}
	if search.OrderNo!="" {
		buider = buider.Where("`order`.`no` = ?",search.OrderNo)
	}
	if search.AddressMobile!="" {
		//buider = buider.Where("address_mobile like ?",search.AddressMobile+"%")
		buider = buider.Where("`order`.open_id in ( select open_id from account where mobile like ?)",search.AddressMobile+"%")
	}
	if len(search.OrderType)>0 {
		buider = buider.Where("`order`.order_type in ?",search.OrderType)
	}
	switch search.PayStatus {
		case 1://1，未付款；
			buider = buider.Where("`order`.pay_status = ?",0)
		case 2://2，已付款
			buider = buider.Where("`order`.pay_status = ?",1)
		case 3://3，已付款
			buider = buider.Where("`order`.pay_status = ?",2)
	}
	switch search.OrderStatus {
		case 1://1，未确认
			buider = buider.Where("`order`.order_status = ?",0)
		case 2://2，已确认；
			buider = buider.Where("`order`.order_status = ?",1)
		case 3://3，已取消；
			buider = buider.Where("`order`.order_status = ?",2)
		case 4://4，无效；
			buider = buider.Where("`order`.order_status = ?",3)
		case 5://5，退货
			buider = buider.Where("`order`.order_status = ?",4)
	}
	
	//show
	if search.Show>0 {
		if search.Show==1 {
			buider = buider.Where("`order`.`show` = ?",1)
		}else{
			buider = buider.Where("`order`.`show` = ?",0)
		}		
	}	
	
	_,err :=buider.Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("create_time",false).LoadStructs(&orders)
	return orders,err
}
func (self *Order) WithItemsCount(searchs interface{},appId string) (int64,error)  {
	var count int64
	
	buider :=db.NewSession().Select("count(`order_item`.id)").From("`order_item`").LeftJoin("order","`order`.`no`=order_item.`no`").Where("`order_item`.app_id=?",appId)
	
	search:=searchs.(OrderSearch)
	if search.MerchantName!="" {
		buider = buider.Where("`order`.merchant_name like ?","%"+search.MerchantName+"%")
	}
	if search.Title!="" {
		buider = buider.Where("`order`.title like ?","%"+search.Title+"%")
	}
	if search.OrderNo!="" {
		buider = buider.Where("`order`.`no` = ?",search.OrderNo)
	}
	if search.AddressMobile!="" {
		//buider = buider.Where("address_mobile like ?",search.AddressMobile+"%")
		buider = buider.Where("`order`.open_id in ( select open_id from account where mobile like ?)",search.AddressMobile+"%")
	}
	if len(search.OrderType)>0 {
		buider = buider.Where("`order`.order_type in ?",search.OrderType)
	}
	switch search.PayStatus {
		case 1://1，未付款；
			buider = buider.Where("`order`.pay_status = ?",0)
		case 2://2，已付款
			buider = buider.Where("`order`.pay_status = ?",1)
		case 3://3，已付款
			buider = buider.Where("`order`.pay_status = ?",2)
	}
	switch search.OrderStatus {
		case 1://1，未确认
			buider = buider.Where("`order`.order_status = ?",0)
		case 2://2，已确认；
			buider = buider.Where("`order`.order_status = ?",1)
		case 3://3，已取消；
			buider = buider.Where("`order`.order_status = ?",2)
		case 4://4，无效；
			buider = buider.Where("`order`.order_status = ?",3)
		case 5://5，退货
			buider = buider.Where("`order`.order_status = ?",4)
	}
	
	//show
	if search.Show>0 {
		if search.Show==1 {
			buider = buider.Where("`order`.`show` = ?",1)
		}else{
			buider = buider.Where("`order`.`show` = ?",0)
		}		
	}	
	//log.Error(buider.ToSql())
	err :=buider.LoadValue(&count)	
	return count,err
}

















































