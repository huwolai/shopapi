package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Order struct  {
	Id int64
	No string
	PayapiNo string
	OpenId string
	AppId string
	Title string
	ActPrice float64
	OmitMoney float64
	Price float64
	Status int
	Json string
}

func NewOrder() *Order {

	return &Order{}
}

func (self *Order) InsertTx(tx *dbr.Tx) (int64,error)  {
	result,err :=tx.InsertInto("order").Columns("no","payapi_no","open_id","app_id","title","act_price","omit_money","price","status","json").Record(self).Exec()
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

func (self *Order) OrderPayapiUpdateWithNo(payapiNo string,status int,no string,appId string) error  {
	sess := db.NewSession()
	_,err :=sess.Update("order").Set("payapi_no",payapiNo).Set("status",status).Where("app_id=?",appId).Where("`no`=?",no).Exec()
	return err
}