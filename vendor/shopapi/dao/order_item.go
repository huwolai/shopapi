package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type OrderItem struct  {
	Id int64
	No string
	AppId string
	OpenId string
	ProdId int64
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	BuyTotalPrice float64
	Json string
}

type OrderItemDetail struct  {
	Id int64
	No string
	AppId string
	OpenId string
	//商户名称
	MerchantName string
	//商户ID
	MerchantId int64
	//商品ID
	ProdId int64
	//商品标题
	ProdTitle string
	//商品cover 封面图 url
	ProdCoverImg string
	//购买数量
	Num int
	OfferUnitPrice float64
	OfferTotalPrice float64
	BuyUnitPrice float64
	BuyTotalPrice float64
	Json string


}

func NewOrderItem()  *OrderItem {

	return &OrderItem{}
}

func NewOrderItemDetail() *OrderItemDetail {

	return &OrderItemDetail{}
}
func (self* OrderItem) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("order_item").Columns("no","app_id","open_id","prod_id","num","offer_unit_price","offer_total_price","buy_unit_price","buy_total_price","json").Record(self).Exec()

	return err
}

func (self *OrderItemDetail) OrderItemWithOrderNo(orderNo []string) ([]*OrderItemDetail,error)  {
	sess := db.NewSession()
	var orderItems []*OrderItemDetail
	_,err :=sess.SelectBySql("select od.*,pt.title prod_title,mt.id merchant_id,mt.`name` merchant_name from order_item od,product pt,merchant_prod mpd,merchant mt where od.prod_id=pt.id and mpd.prod_id=pt.id and mpd.merchant_id=mt.id and  `no` in ?",orderNo).LoadStructs(&orderItems)
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

