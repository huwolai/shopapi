package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"time"
	"fmt"
)

type orderItemPurchaseCodes struct  {
	Id 				int64
	order_item_id	int64
	No				string
	Codes 			string
	ProdId 			int64
	BuyTime 		string
}

type ProdPurchaseCodes struct  {
	AppId 		string `json:"app_id"`
	Id 			int64 `json:"id"`
	ProdId		int64 `json:"prod_id"`
	Sku 		string `json:"sku"`
	Codes 		string `json:"codes"`
	Num 		int `json:"num"`
	OpenId 		int `json:"open_id"`
	OpenStatus 	int64 `json:"open_status"`
}

func OrderItemPurchaseCodesAdd(tx *dbr.Tx,orderItemId int64,no string,prodId int64,codes string) error {
	var purchaseCodes *orderItemPurchaseCodes
	purchaseCodes.order_item_id	=orderItemId
	purchaseCodes.No			=no
	purchaseCodes.Codes			=codes
	purchaseCodes.ProdId		=prodId
	purchaseCodes.BuyTime		=fmt.Sprintf("%d",time.Now().UnixNano()/1e6)
	
	_,err :=tx.InsertInto("order_item_purchase_codes").Columns("order_item_id","no","codes","prod_id").Record(purchaseCodes).Exec()
	return err
}
//一元购生成购买码
func ProductAndPurchaseCodesAdd(prodPurchaseCodes *ProdPurchaseCodes) error {
	_,err :=db.NewSession().InsertInto("prod_purchase_codes").Columns("sku","app_id","prod_id","codes","num").Record(prodPurchaseCodes).Exec()
	return err
}
//一元购减去购买码
func ProductAndPurchaseCodesMinus(tx *dbr.Tx,id int64,num int,newNum int,newCodes string) error  {

	_,err :=tx.UpdateBySql("update prod_purchase_codes set codes=?,num=? where id=? and num=?",newCodes,newNum,id,num).Exec()

	return err
}
//一元购购买码
func ProductAndPurchaseCodesTx(prodPurchaseCodes *ProdPurchaseCodes,tx *dbr.Tx) (*ProdPurchaseCodes,error)  {
	var codes *ProdPurchaseCodes
	_,err :=tx.Select("*").From("prod_purchase_codes").Where("sku=?",prodPurchaseCodes.Sku).Where("app_id=?",prodPurchaseCodes.AppId).Where("prod_id=?",prodPurchaseCodes.ProdId).LoadStructs(&codes)
	return codes,err
}
func ProductAndPurchaseCodes(prodPurchaseCodes *ProdPurchaseCodes) (*ProdPurchaseCodes,error)  {
	var codes *ProdPurchaseCodes
	_,err :=db.NewSession().Select("*").From("prod_purchase_codes").Where("sku=?",prodPurchaseCodes.Sku).Where("app_id=?",prodPurchaseCodes.AppId).Where("prod_id=?",prodPurchaseCodes.ProdId).LoadStructs(&codes)
	return codes,err
}
//开奖中
func ProductAndPurchaseCodesOpening(tx *dbr.Tx,prodPurchaseCodes *ProdPurchaseCodes,openTime string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=tx.Update("prod_purchase_codes")	
	builder = builder.Set("status",1)
	builder = builder.Set("open_time",openTime)	
	builder = builder.Where("sku=?",prodPurchaseCodes.Sku)
	builder = builder.Where("app_id=?",prodPurchaseCodes.AppId)
	builder = builder.Where("prod_id=?",prodPurchaseCodes.ProdId)		
	_,err :=builder.Exec()
	return err
}
//开奖
func ProductAndPurchaseCodesOpened(tx *dbr.Tx,prodPurchaseCodes *ProdPurchaseCodes,openId string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=tx.Update("prod_purchase_codes")	
	builder = builder.Set("open_id",openId)
	builder = builder.Set("status",2)
	builder = builder.Where("sku=?",prodPurchaseCodes.Sku)
	builder = builder.Where("app_id=?",prodPurchaseCodes.AppId)
	builder = builder.Where("prod_id=?",prodPurchaseCodes.ProdId)		
	_,err :=builder.Exec()
	return err
}



















