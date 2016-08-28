package dao

import "github.com/gocraft/dbr"

type OrderCoupon struct  {
	Id int64
	OrderNo string
	CouponNo string
	CouponAmount float64
	CouponToken string
	NotifyUrl string
	AppId string
	Status int
}

func NewOrderCoupon() *OrderCoupon  {

	return &OrderCoupon{}
}

func (self *OrderCoupon) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("order_coupon").Columns("app_id","notify_url","order_no","coupon_no","coupon_amount","coupon_token","status").Exec()

	return err
}

func (self *OrderCoupon) DeleteWithOrderNoTx(status int,orderNo string,tx *dbr.Tx) error  {

	_,err :=tx.DeleteFrom("order_coupon").Where("status=?",status).Where("order_no=?",orderNo).Exec()

	return err
}