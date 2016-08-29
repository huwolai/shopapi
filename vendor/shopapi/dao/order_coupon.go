package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type OrderCoupon struct  {
	Id int64
	AppId string
	OpenId string
	OrderNo string
	CouponCode string
	TrackCode string
	CouponAmount float64
	CouponToken string
	NotifyUrl string
	Status int
}

func NewOrderCoupon() *OrderCoupon  {

	return &OrderCoupon{}
}

func (self *OrderCoupon) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("order_coupon").Columns("app_id","open_id","notify_url","order_no","coupon_code","track_code","coupon_amount","coupon_token","status").Exec()

	return err
}

func (self *OrderCoupon) DeleteWithOrderNoTx(status int,orderNo string,tx *dbr.Tx) error  {

	_,err :=tx.DeleteFrom("order_coupon").Where("status=?",status).Where("order_no=?",orderNo).Exec()

	return err
}

func (self *OrderCoupon) WithOrderNo(orderNo string,appId string) ([]*OrderCoupon,error)  {

	var list []*OrderCoupon
	_,err :=db.NewSession().Select("*").From("order_coupon").Where("order_no=?",orderNo).Where("app_id=?",appId).LoadStructs(&list)

	return list,err
}