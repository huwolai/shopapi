package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Order struct  {
	Id int64
	No string
	Code  string
	PayapiNo string
	OpenId string
	AppId string
	AddressId int64
	Title string
	ActPrice float64
	OmitMoney float64
	Price float64
	Flag string
	Status int
	Json string
}

type OrderDetail struct  {
	Id int64
	No string
	PayapiNo string
	OpenId string
	AppId string
	AddressId int64
	Address string
	Title string
	ActPrice float64
	OmitMoney float64
	Price float64
	Status int
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

func (self *Order) InsertTx(tx *dbr.Tx) (int64,error)  {
	result,err :=tx.InsertInto("order").Columns("no","address_id","payapi_no","code","open_id","app_id","title","act_price","omit_money","price","status","flag","json").Record(self).Exec()
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

func (self *OrderDetail) OrderDetailWithUser(openId string,status []int,appId string) ([]*OrderDetail,error)  {

	sess := db.NewSession()
	var orders []*OrderDetail
	_,err :=sess.SelectBySql("select od.*,ad.address from `order` od left join address ad on od.address_id=ad.id ad where od.open_id=? and od.app_id=? and status in ?",openId,appId,status).OrderDir("create_time",false).LoadStructs(&orders)
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
		orderItemDetails,err :=orderItemDetail.OrderItemWithOrderNo(ordernos)
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

func (self *Order) OrderPayapiUpdateWithNoAndCode(payapiNo string,code string,status int,no string,appId string) error  {
	sess := db.NewSession()
	_,err :=sess.Update("order").Set("payapi_no",payapiNo).Set("code",code).Set("status",status).Where("app_id=?",appId).Where("`no`=?",no).Exec()
	return err
}

func (self *Order) UpdateWithStatus(status int,orderNo string) error {

	_,err :=db.NewSession().Update("order").Set("status",status).Where("no=?",orderNo).Exec()

	return err
}