package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"time"
	"fmt"
)

type OrderItemPurchaseCode struct  {
	Id 			int64
	orderItemId	int64
	No			string
	Codes 		string
	ProdId 		int64
	BuyTime 	string
}

type ProdPurchaseCode struct  {
	AppId 		string `json:"app_id"`
	Id 			int64 `json:"id"`
	ProdId		int64 `json:"prod_id"`
	Sku 		string `json:"sku"`
	Codes 		string `json:"codes"`
	Num 		int `json:"num"`
	OpenId 		string `json:"open_id"`
	OpenStatus 	int64 `json:"open_status"`
	OpenTime 	int64 `json:"Open_time"`
}
//==
func OrderItemPurchaseCodesAdd(tx *dbr.Tx,orderItemId int64,no string,prodId int64,codes string) error {
	var orderItemPurchaseCode *OrderItemPurchaseCode
	orderItemPurchaseCode.orderItemId	=orderItemId
	orderItemPurchaseCode.No			=no
	orderItemPurchaseCode.Codes			=codes
	orderItemPurchaseCode.ProdId		=prodId
	orderItemPurchaseCode.BuyTime		=fmt.Sprintf("%d",time.Now().UnixNano()/1e6)
	
	_,err :=tx.InsertInto("order_item_purchase_codes").Columns("order_item_id","no","codes","prod_id").Record(orderItemPurchaseCode).Exec()
	return err
}
//==
func OrderItemPurchaseCodes(prodId int64,limit int64)  ([]*OrderItemPurchaseCode,error)  {
	var orderItemPurchaseCode []*OrderItemPurchaseCode
	_,err :=db.NewSession().SelectBySql("select * from order_item_purchase_codes where prod_id = ? order by id desc limit ?",prodId,limit).LoadStructs(&orderItemPurchaseCode)
	return  orderItemPurchaseCode,err
}
//一元购生成购买码
func ProductAndPurchaseCodesAdd(prodPurchaseCode *ProdPurchaseCode) error {
	_,err :=db.NewSession().InsertInto("prod_purchase_codes").Columns("sku","app_id","prod_id","codes","num").Record(prodPurchaseCode).Exec()
	return err
}
//一元购减去购买码
func ProductAndPurchaseCodesMinus(tx *dbr.Tx,id int64,num int,newNum int,newCodes string) error  {

	_,err :=tx.UpdateBySql("update prod_purchase_codes set codes=?,num=? where id=? and num=?",newCodes,newNum,id,num).Exec()

	return err
}
//一元购购买码
func ProductAndPurchaseCodesTx(prodPurchaseCode *ProdPurchaseCode,tx *dbr.Tx) (*ProdPurchaseCode,error)  {
	var codes *ProdPurchaseCode
	_,err :=tx.Select("*").From("prod_purchase_codes").Where("sku=?",prodPurchaseCode.Sku).Where("app_id=?",prodPurchaseCode.AppId).Where("prod_id=?",prodPurchaseCode.ProdId).LoadStructs(&codes)
	return codes,err
}
func ProdPurchaseCodeWithProdIds(prodIds []int64) ([]*ProdPurchaseCode,error){
	var prodPurchaseCode []*ProdPurchaseCode
	_,err :=db.NewSession().SelectBySql("select * from prod_purchase_codes where prod_id in ?",prodIds).LoadStructs(&prodPurchaseCode)
	return  prodPurchaseCode,err
}
func ProdPurchaseCodeWithProdId(prodId int64) (*ProdPurchaseCode,error){
	var prodPurchaseCode *ProdPurchaseCode
	_,err :=db.NewSession().SelectBySql("select * from prod_purchase_codes where prod_id = ?",prodId).LoadStructs(&prodPurchaseCode)
	return  prodPurchaseCode,err
}
func ProdPurchaseCodeWithOpenStatus(openStatus int64) ([]*ProdPurchaseCode,error){
	var prodPurchaseCode []*ProdPurchaseCode
	_,err :=db.NewSession().SelectBySql("select * from prod_purchase_codes where open_status = ?",openStatus).LoadStructs(&prodPurchaseCode)
	return  prodPurchaseCode,err
}
//开奖中
func ProductAndPurchaseCodesOpening(tx *dbr.Tx,prodPurchaseCode *ProdPurchaseCode,openTime string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=tx.Update("prod_purchase_codes")	
	builder = builder.Set("status",1)
	builder = builder.Set("open_time",openTime)	
	builder = builder.Where("sku=?",prodPurchaseCode.Sku)
	builder = builder.Where("app_id=?",prodPurchaseCode.AppId)
	builder = builder.Where("prod_id=?",prodPurchaseCode.ProdId)		
	_,err :=builder.Exec()
	return err
}
//开奖
func ProductAndPurchaseCodesOpened(prodId int64,openId string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=db.NewSession().Update("prod_purchase_codes")	
	builder = builder.Set("status",2)
	builder = builder.Set("open_id",openId)
	builder = builder.Where("prod_id=?",prodId)		
	_,err :=builder.Exec()
	return err
}
//开奖
func GetOpenIdbyOpenCode(prodId int64,openCode string) (string,error)  {
	builder :=db.NewSession().Select("open_id").From("order_item_purchase_codes").Join("order","order_item_purchase_codes.no=order.no")
	builder = builder.Where("order_item_purchase_codes.prod_id =",prodId)
	builder = builder.Where("order_item_purchase_codes.codes =",openCode)
	
	var openId string
	err :=builder.LoadValue(&openId)

	return openId,err	
}










