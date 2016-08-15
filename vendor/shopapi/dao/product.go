package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Product struct  {
	Id int64
	AppId string
	//商品标题
	Title string
	//商品描述
	Description string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//商品状态
	Status int
	//附加数据
	Json string
}

func NewProduct() *Product  {

	return &Product{}
}

func (self *Product) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("product").Columns("title","app_id","description","price","dis_price","json","status").Record(self).Exec()
	if err !=nil {

		return 0,err
	}
	pid,err :=  result.LastInsertId()
	return pid,err
}

func (self *Product) ProductListWithCategory(appId string,categoryId int64) ([]*Product,error)  {
	session := db.NewSession()
	var prodList []*Product
	_,err :=session.SelectBySql("select pt.* from prod_category pc,product pt where pc.app_id=pt.app_id and pc.prod_id=pt.id and pc.category_id=? and pt.app_id=?",categoryId,appId).LoadStructs(&prodList)
	return prodList,err
}


func (self *Product) ProductWithId(appId string,id int64) (*Product,error)  {
	sess :=db.NewSession()
	var prod *Product
	_,err :=sess.Select("*").From("product").Where("app_id=?",appId).Where("id=?",id).LoadStructs(&prod)

	return prod,err
}