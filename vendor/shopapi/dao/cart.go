package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Cart struct {
	Id			int64 `json:"id"`
	ProdId		int64 `json:"prod_id"`
	SkuNo		string `json:"sku_no"`
	Num			int64  `json:"num"`
	OpenId		string `json:"open_id"`
	//BaseDModel	
}

func NewCart() *Cart  {
	return &Cart{}
}

func (self *Cart)CartList(openId string) ([]Cart,error) {
	var cart []Cart
	_,err :=db.NewSession().Select("*").From("cart").Where("open_id=?",openId).LoadStructs(&cart)
	return cart,err
}

func (self *Cart)CartAddToList(cart Cart) error {
	_,err :=db.NewSession().InsertInto("cart").Columns("prod_id","sku_no","num","open_id").Record(cart).Exec()
	return err	
}
func (self *Cart)CartDelFromList(openId string,id uint64) error {
	_,err :=db.NewSession().DeleteFrom("cart").Where("id=?",id).Where("open_id=?",openId).Limit(1).Exec()
	return err
}