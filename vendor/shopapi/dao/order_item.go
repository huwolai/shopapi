package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type OrderItem struct  {
	Id int64
	No string
	Title string
	AppId string
	OpenId string
	ProdId int64
	DbnNo string
	SkuNo string
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	//购买金额
	BuyTotalPrice float64
	//佣金
	DbnAmount float64
	//优惠金额
	CouponAmount float64
	//省略的金额
	OmitMoney float64
	//商家所得金额
	MerchantAmount float64
	Json string
	
	GmOrdernum string
	GmPassnum string
	GmPassway string
	WayStatus int64
}

type OrderItemDetail struct  {
	Id int64
	No string
	Title string
	AppId string
	OpenId string
	//商户名称
	MerchantName string
	//商户ID
	MerchantId int64
	//商品ID
	ProdId int64
	//分销编号
	DbnNo string
	//商品SKU
	SkuNo string
	//商品标题
	ProdTitle string
	//商品cover 封面图 url
	ProdCoverImg string
	//商品标识
	ProdFlag string
	//购买数量
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	BuyTotalPrice float64
	DbnAmount float64
	MerchantAmount float64
	Json string

	GmOrdernum string
	GmPassnum string
	GmPassway string
	WayStatus int64

}

func NewOrderItem()  *OrderItem {

	return &OrderItem{}
}

func NewOrderItemDetail() *OrderItemDetail {

	return &OrderItemDetail{}
}
func (self* OrderItem) InsertTx(tx *dbr.Tx) error {
	_,err :=tx.InsertInto("order_item").Columns("no","title","dbn_no","coupon_amount","omit_money","dbn_amount","merchant_amount","app_id","prod_id","sku_no","num","offer_unit_price","offer_total_price","buy_unit_price","buy_total_price","json").Record(self).Exec()

	return err
}

func (self *OrderItem) OrderItemWithOrderNo(orderNo string) ([]*OrderItem,error)  {

	var items []*OrderItem
	_,err :=db.NewSession().Select("*").From("order_item").Where("no=?",orderNo).LoadStructs(&items)

	return items,err
}

func (self *OrderItem) UpdateAmountWithIdTx(dbnAmount float64,omitMoney float64,couponAmount float64,merchantAmount float64,id int64,tx *dbr.Tx)  error {
	_,err :=tx.Update("order_item").Set("dbn_amount",dbnAmount).Set("omit_money",omitMoney).Set("coupon_amount",couponAmount).Set("merchant_amount",merchantAmount).Where("id=?",id).Exec()

	return err
}

func (self *OrderItemDetail) OrderItemDetailWithOrderNo(orderNo []string) ([]*OrderItemDetail,error)  {
	sess := db.NewSession()
	var orderItems []*OrderItemDetail
	_,err :=sess.SelectBySql("select od.*,pt.title prod_title,mt.id merchant_id,mt.`name` merchant_name,mt.open_id,pt.flag prod_flag from order_item od,product pt,merchant_prod mpd,merchant mt where od.prod_id=pt.id and mpd.prod_id=pt.id and mpd.merchant_id=mt.id and  `no` in ?",orderNo).LoadStructs(&orderItems)
	if err !=nil {
		return nil,err
	}
	if orderItems!=nil{
		err :=fillOrderItemImg(orderItems)
		if err!=nil{
			return nil,err
		}
	}
	return orderItems,err
}

func fillOrderItemImg(orderItems []*OrderItemDetail) error  {
	prodids := make([]int64,0)
	for _,orderItem :=range orderItems {
		prodids = append(prodids,orderItem.ProdId)
	}
	imgDetail := NewProdImgsDetail()
	prodImgDetailList,err  := imgDetail.ProdImgsWithProdIds(prodids)
	if err!=nil{
		return err
	}
	imgDetailMap  :=make(map[int64][]*ProdImgsDetail)
	if prodImgDetailList!=nil {
		for _,prodimgDetal :=range prodImgDetailList {
			prodimgDetals := imgDetailMap[prodimgDetal.ProdId]
			if prodimgDetals==nil {
				prodimgDetals =make([]*ProdImgsDetail,0)
			}
			prodimgDetals = append(prodimgDetals,prodimgDetal)
			imgDetailMap[prodimgDetal.ProdId]=prodimgDetals
		}
	}
	log.Debug(imgDetailMap)
	for _,orderItem :=range orderItems {
		imgs := imgDetailMap[orderItem.ProdId]
		if imgs!=nil&&len(imgs)>0{
			orderItem.ProdCoverImg = imgs[0].Url
		}
	}

	return nil
}

func (self *OrderItem)OrdersAddNumWithNo(orderNo string,appId string,ordernum string) error {
	_,err :=db.NewSession().Update("order_item").Set("gm_ordernum",ordernum).Where("no=?",orderNo).Where("app_id=?",appId).Exec()
	return err
}
func (self *OrderItem)OrdersAddNumWithPassnum(orderNo string,appId string,passnum string,passway string) error {
	_,err :=db.NewSession().Update("order_item").Set("gm_passnum",passnum).Set("gm_passway",passway).Where("no=?",orderNo).Where("app_id=?",appId).Exec()
	return err
}

//商品已购买数量
func ProdOrderCountWithId(prodId int64,OpenId string,Date string) (int64,error) {
	//count, err :=db.NewSession().Select("sum(order_item.num)").From("order_item").LeftJoin("order","order_item.no=order.no").Where("order_item.prod_id = ?", prodId).Where("order.open_id = ?", OpenId).ReturnInt64()
	
	buider :=db.NewSession().Select("sum(order_item.num)").From("order_item").LeftJoin("order","order_item.no=order.no")
	
	if len(Date)>0 {
		buider = buider.Where("order_item.create_time like ?",Date+"%")
	}
	
	count, err :=buider.Where("order.order_status <= ?", 1).Where("order_item.prod_id = ?", prodId).Where("order.open_id = ?", OpenId).ReturnInt64()
	
	
	return count,err
}

























