package dao

import "github.com/gocraft/dbr"

type OrderItem struct  {
	Id int64
	No string
	AppId string
	OpenId string
	MOpenId string
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
func (self* OrderItem) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("order_item").Columns("no","app_id","open_id","m_open_id","prod_id","num","offer_unit_price","offer_total_price","buy_unit_price","buy_total_price","json").Record(self).Exec()

	return err
}