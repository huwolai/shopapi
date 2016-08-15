package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type OrderItem struct  {
	Id int64
	No string
	AppId string
	OpenId string
	ProdId int64
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	BuyTotalPrice float64
	Json string
}

type OrderItemDetail struct  {
	Id int64
	No string
	AppId string
	OpenId string
	ProdId int64
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	BuyTotalPrice float64
	Json string

}

func NewOrderItem()  *OrderItem {

	return &OrderItem{}
}

func NewOrderItemDetail() *OrderItemDetail {

	return &OrderItemDetail{}
}
func (self* OrderItem) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("order_item").Columns("no","app_id","open_id","prod_id","num","offer_unit_price","offer_total_price","buy_unit_price","buy_total_price","json").Record(self).Exec()

	return err
}

func (self *OrderItemDetail) OrderItemWithOrderNo(orderNo []string) ([]*OrderItemDetail,error)  {
	sess := db.NewSession()
	var orderItems []*OrderItemDetail
	_,err :=sess.SelectBySql("select * from order_item where order_no ?",orderNo).LoadStructs(&orderItems)

	return orderItems,err

}