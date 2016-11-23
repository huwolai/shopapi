package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
	"github.com/gocraft/dbr"
)

type Cart struct {
	Id			int64 `json:"id"`
	ProdId		int64 `json:"prod_id"`
	SkuNo		string `json:"sku_no"`
	Num			int64  `json:"num"`
	OpenId		string `json:"open_id"`
	
	Mainimg		string `json:"mainimg"`
	Title		string `json:"title"`
	Price		float64 `json:"price"`
	DisPrice	float64 `json:"dis_price"`
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
func (self *Cart)CartNumInLIst(cart Cart,tx *dbr.Tx)(int64,error) {
	var num int64
	_,err :=tx.Select("num").From("cart").Where("open_id=?",cart.OpenId).Where("sku_no=?",cart.SkuNo).Where("prod_id=?",cart.ProdId).LoadStructs(&num)
	return num,err
}
func (self *Cart)CartExistInList(cart Cart,tx *dbr.Tx)(int64,error) {
	var count int64
	_,err :=tx.Select("count(id)").From("cart").Where("open_id=?",cart.OpenId).Where("sku_no=?",cart.SkuNo).Where("prod_id=?",cart.ProdId).LoadStructs(&count)
	return count,err
}
func (self *Cart)CartAddNumToList(cart Cart,tx *dbr.Tx) error {
	_,err :=tx.UpdateBySql("update cart set num=num+? where open_id=? and sku_no=? and prod_id=?",cart.Num,cart.OpenId,cart.SkuNo,cart.ProdId).Exec()
	return err
}
func (self *Cart)CartMinusFromList(cart Cart,tx *dbr.Tx) error {
	_,err :=tx.UpdateBySql("update cart set num=num-? where open_id=? and sku_no=? and prod_id=?",cart.Num,cart.OpenId,cart.SkuNo,cart.ProdId).Exec()
	return err
}
func (self *Cart)CartAddToList(cart Cart,tx *dbr.Tx) error {
	_,err :=tx.InsertInto("cart").Columns("prod_id","sku_no","num","open_id","mainimg","title","price","dis_price").Record(cart).Exec()
	return err
}
func (self *Cart)CartDelFromList(openId string,id uint64) error {
	_,err :=db.NewSession().DeleteFrom("cart").Where("id=?",id).Where("open_id=?",openId).Limit(1).Exec()
	return err
}
func (self *Cart)CartUpdateList(cart Cart,tx *dbr.Tx) error {
	_,err :=db.NewSession().Update("cart").Set("num",cart.Num).Where("open_id=?",cart.OpenId).Where("sku_no=?",cart.SkuNo).Where("prod_id=?",cart.ProdId).Exec()
	return err
}





















