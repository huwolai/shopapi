package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"time"
	//"fmt"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
)

type OrderItemPurchaseCode struct  {
	Id 			int64	`json:"id"`
	OrderItemId	int64	`json:"order_item_id"`
	No			string	`json:"no"`
	Codes 		string	`json:"codes"`
	ProdId 		int64	`json:"prod_id"`
	BuyTime 	int64	`json:"buy_time"`
}
type OrderItemPurchaseCodeRrecord struct  {
	OrderItemPurchaseCode
	Mobile	string	`json:"mobile"`
	OpenId	string	`json:"open_id"`
}

type ProdPurchaseCode struct  {
	AppId 		string `json:"app_id"`
	Id 			int64 `json:"id"`
	ProdId		int64 `json:"prod_id"`
	Sku 		string `json:"sku"`
	Codes 		string `json:"codes"`
	Num 		int `json:"num"`
	Nums 		int `json:"nums"`
	OpenId 		string `json:"open_id"`
	OpenStatus 	int64 `json:"open_status"`
	OpenTime 	int64 `json:"Open_time"`
	OpenCode	string `json:"Open_code"`
	OpenMobile	string `json:"Open_mobile"`
}
type UserOpen struct  {
	OpenId string
	Mobile string
}
type OrdersYygSearch struct {
	ProdTitle string
}
//==
func OrderItemPurchaseCodesAdd(tx *dbr.Tx,orderItemId int64,no string,prodId int64,codes string,index int64) error {
	var orderItemPurchaseCode OrderItemPurchaseCode
	orderItemPurchaseCode.OrderItemId	=orderItemId
	orderItemPurchaseCode.No			=no
	orderItemPurchaseCode.Codes			=codes
	orderItemPurchaseCode.ProdId		=prodId
	//orderItemPurchaseCode.BuyTime		=fmt.Sprintf("%d",time.Now().UnixNano()/1e6)
	orderItemPurchaseCode.BuyTime		=time.Now().UnixNano()/1e6+index
	
	_,err :=tx.InsertInto("order_item_purchase_codes").Columns("order_item_id","no","codes","prod_id","buy_time").Record(orderItemPurchaseCode).Exec()
	return err
}
//==
func OrderItemPurchaseCodesWithTime(time int64,limit int64)  ([]*OrderItemPurchaseCode,error)  {
	var orderItemPurchaseCode []*OrderItemPurchaseCode
	//buy_time 毫秒
	_,err :=db.NewSession().SelectBySql("select * from order_item_purchase_codes where buy_time <= ? order by id desc limit ?",time*1000+999,limit).LoadStructs(&orderItemPurchaseCode)
	return  orderItemPurchaseCode,err
}
//==
func OrderItemPurchaseCodesRrecordWithTime(time int64,limit int64)  ([]*OrderItemPurchaseCodeRrecord,error)  {
	var items []*OrderItemPurchaseCodeRrecord
	//buy_time 毫秒
	_,err :=db.NewSession().SelectBySql("select c.*,order.address_mobile as mobile,order.open_id from order_item_purchase_codes as c left join `order` on c.`no`=`order`.`no` where c.buy_time <= ? order by id desc limit ?",time*1000+999,limit).LoadStructs(&items)
	return  items,err
}
func OrderItemPurchaseCodesWithProdId(prodId int64)  (int64,error)  {
	var count int64
	//_,err :=db.NewSession().SelectBySql("SELECT count(*) from (select * from order_item_purchase_codes where prod_id =? GROUP BY codes) c",prodId).LoadStructs(&count)
	_,err :=db.NewSession().SelectBySql("SELECT nums from prod_purchase_codes where prod_id=?",prodId).LoadStructs(&count)
	return  count,err
}
func OrderItemPurchaseCodesWithNo(orderNo string)  ([]string,error)  {
	var codes []string
	_,err :=db.NewSession().SelectBySql("select codes from order_item_purchase_codes where no = ?",orderNo).LoadStructs(&codes)
	return  codes,err
}
//一元购生成购买码
func ProductAndPurchaseCodesAdd(prodPurchaseCode *ProdPurchaseCode) error {
	prodPurchaseCode.Nums=prodPurchaseCode.Num
	_,err :=db.NewSession().InsertInto("prod_purchase_codes").Columns("sku","app_id","prod_id","codes","num","nums").Record(prodPurchaseCode).Exec()
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
func ProdPurchaseCodes(search OrdersYygSearch,pIndex uint64,pSize uint64) ([]*ProdPurchaseCode,error){
	var prodPurchaseCode []*ProdPurchaseCode
	builder:=db.NewSession().Select("*").From("prod_purchase_codes")
	
	if len(search.ProdTitle)>0 {
		builder = builder.Where("prod_id in (select id from product where title like ?)","%"+search.ProdTitle+"%")
	}
	
	_,err :=builder.OrderDir("id",false).Limit(pSize).Offset((pIndex-1)*pSize).LoadStructs(&prodPurchaseCode)
	//_,err :=db.NewSession().SelectBySql("select * from prod_purchase_codes order by id desc limit ?,? ",(pIndex-1)*pSize,pSize).LoadStructs(&prodPurchaseCode)
	return  prodPurchaseCode,err
}
func ProdPurchaseCodesCount(search OrdersYygSearch) (int64,error){
	builder:=db.NewSession().Select("count(id)").From("prod_purchase_codes")
	
	if len(search.ProdTitle)>0 {
		builder = builder.Where("prod_id in (select id from product where title like ?)","%"+search.ProdTitle+"%")
	}
	
	count,err :=builder.ReturnInt64()
	return  count,err
}
//开奖中
func ProductAndPurchaseCodesOpening(tx *dbr.Tx,prodPurchaseCode *ProdPurchaseCode,openTime string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=tx.Update("prod_purchase_codes")	
	builder = builder.Set("open_status",1)
	builder = builder.Set("open_time",openTime)	
	builder = builder.Where("sku=?",prodPurchaseCode.Sku)
	builder = builder.Where("app_id=?",prodPurchaseCode.AppId)
	builder = builder.Where("prod_id=?",prodPurchaseCode.ProdId)		
	_,err :=builder.Exec()
	return err
}
//开奖
func ProductAndPurchaseCodesOpened(prodId int64,openId string,mobile string,openCode string) error  {
	//_,err :=tx.UpdateBySql("update prod_purchase_codes set open_id=?,open_time=? where sku=? and num=?",openId,openTime,id,num).Exec()	
	builder:=db.NewSession().Update("prod_purchase_codes")	
	builder = builder.Set("open_status",2)
	builder = builder.Set("open_code",openCode)
	builder = builder.Set("open_id",openId)
	builder = builder.Set("open_mobile",mobile)
	builder = builder.Where("prod_id=?",prodId)		
	_,err :=builder.Exec()
	return err
}
//开奖状态
func ProductAndPurchaseCodesOpenedStatus() error  {
	
	/* count,_ :=db.NewSession().SelectBySql("select count(id) from prod_purchase_codes where open_time<=UNIX_TIMESTAMP()+52 and open_status=1 and open_time>?",0).ReturnInt64()
	
	if count<1 {
		return nil
	} */
	_,err :=db.NewSession().UpdateBySql("update prod_purchase_codes set open_status=? where open_time<=UNIX_TIMESTAMP()+52 and open_status=1 and open_time>0",2).Exec()	
	return err
}
//开奖
func GetOpenIdbyOpenCode(prodId int64,openCode string) (*UserOpen,error)  {
	var user *UserOpen
	
	_,err :=db.NewSession().SelectBySql("SELECT order.open_id,account.mobile FROM order_item_purchase_codes JOIN `order` ON order_item_purchase_codes.no=`order`.no  JOIN `account` ON `order`.open_id=account.open_id WHERE (order_item_purchase_codes.prod_id =?) AND (order_item_purchase_codes.codes =?)",prodId,openCode).LoadStructs(&user)
	
	return user,err	
}










