package service

import (
	"shopapi/dao"	
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"fmt"
	"shopapi/comm"
	"time"
	"strconv"
	//"strings"
)

type ProdPurchaseCode struct  {
	Id 			int64  `json:"id"`
	ProdId		int64  `json:"prod_id"`
	Sku 		string `json:"sku"`
	OpenStatus 	int64  `json:"open_status"`	
	PordTitle 	string `json:"pord_title"`	
		
	OrderDto	 	
	
	OpenMobile	string  `json:"open_mobile"`
	YdgyName 	string `json:"ydgy_name"`
}
type OrderDto struct  {
	Json string `json:"json"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	MOpenId string
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	AddressId int64 `json:"address_id"`
	Title string `json:"title"`
	OrderNo string `json:"order_no"`
	RejectCancelReason string `json:"reject_cancel_reason"`
	CancelReason string `json:"cancel_reason"`
	OrderStatus int `json:"order_status"`
	PayStatus int `json:"pay_status"`
	CouponAmount float64 `json:"coupon_amount"`
	RealPrice float64 `json:"real_price"`
	PayPrice float64 `json:"pay_price"`
	CreateTime string `json:"create_time"`
	
	GmOrdernum string `json:"ordernum"`
	GmPassnum string `json:"passnum"`
	GmPassway string `json:"passway"`
	WayStatus int64 `json:"way_status"`
	
	DetailTitle []string `json:"detailtitle"`
	
	Address string	 `json:"address"`
	AddressMobile string	 `json:"address_mobile"`
	AddressName string	 `json:"address_name"`
	
	Show int `json:"show"`
	Mobile 	string	 `json:"mobile"`
	YdgyName  	 string	 `json:"ydgy_name"`
}

func OrdersYygWin(search dao.OrdersYygSearch,appId string,pIndex uint64,pSize uint64) ([]*ProdPurchaseCode,int64,error) {
	count,err:=dao.ProdPurchaseCodesCount(search)
	if err!=nil {	
		return nil,0,err
	}	
	
	prods,err:=dao.ProdPurchaseCodes(search,pIndex,pSize)
	if err!=nil {		
		return nil,0,err
	}
	
	log.Info(prods==nil)
	
	itemsDto :=make([]*ProdPurchaseCode,0)
	orderDao :=dao.NewOrder()
	accountDao :=dao.NewAccount()
	prodDao	 :=dao.NewProduct()
	var orderItem []*dao.OrderItem
	for _,item :=range prods {
		dto :=&ProdPurchaseCode{}		
		dto.Id			=item.Id
		dto.ProdId		=item.ProdId
		dto.Sku			=item.Sku
		dto.OpenStatus	=item.OpenStatus
		dto.OpenMobile	=item.OpenMobile
		//OpenCode
		
		prod,_:=prodDao.ProductWithId(item.ProdId,appId)
		if prod!=nil {		
			dto.PordTitle 			= prod.Title
		}		
		
		if item.OpenStatus==2 {			
			order,err:=orderDao.OrderWithPordYyg(item.OpenCode,item.ProdId)			
			if err!=nil {		
				return nil,0,err
			}
			if order!=nil {
				account,err :=accountDao.AccountWithOpenId(item.OpenId,appId)
				if err==nil && account!=nil {
					order.YdgyName	=account.YdgyName	
				}
				
				orderItem,_=OrderItems(order.No);			
				if len(orderItem)>0 {
					order.GmOrdernum	=orderItem[0].GmOrdernum
					order.GmPassnum	=orderItem[0].GmPassnum
					order.GmPassway	=orderItem[0].GmPassway
					order.WayStatus	=orderItem[0].WayStatus								
				}
				orderToA(dto,order)
			}
		}
		itemsDto = append(itemsDto,dto)
	}
	
	return itemsDto,count,nil
}
func OrdersYygRecord(prodId string,appId string,pIndex uint64,pSize uint64) ([]*dao.OrderYyg,int64,error) {
	order,count,err:=dao.NewOrder().OrdersWithPordYyg(prodId,pIndex,pSize)
	if err!=nil {		
		return nil,0,err
	}
	return order,count,nil
}
func orderToA(dto *ProdPurchaseCode,order *dao.Order) *ProdPurchaseCode {

	dto.AddressId 			= order.AddressId
	dto.AppId 				= order.AppId
	dto.CancelReason 		= order.CancelReason
	dto.Json				= order.Json
	dto.MerchantId 			= order.MerchantId
	dto.MerchantName 		= order.MerchantName
	dto.MOpenId 			= order.MOpenId
	dto.OpenId 				= order.OpenId
	dto.OrderNo 			= order.No
	dto.Title 				= order.Title
	dto.PayStatus 			= order.PayStatus
	dto.OrderStatus 		= order.OrderStatus
	dto.RealPrice 			= order.RealPrice
	dto.PayPrice 			= order.PayPrice
	dto.CreateTime 			= qtime.ToyyyyMMddHHmm(order.CreateTime)
	dto.GmOrdernum 			= order.GmOrdernum
	dto.GmPassnum 			= order.GmPassnum
	dto.GmPassway 			= order.GmPassway
	dto.WayStatus 			= order.WayStatus
	dto.DetailTitle 		= order.DetailTitle
	
	dto.Address 			= order.Address
	dto.AddressMobile 		= order.AddressMobile
	dto.AddressName 		= order.AddressName
	
	dto.Show 				= order.Show
	dto.Mobile 				= order.Mobile
	dto.YdgyName 			= order.YdgyName
	
	
	

	return 	dto
}

//定时开奖
func PurchaseCodesOpen(prodId int64)  {
	log.Info("定时开奖"+fmt.Sprintf("%d", prodId))
	prod,_		:=dao.ProdPurchaseCodeWithProdId(prodId)
	orderItem,_ :=dao.OrderItemPurchaseCodesWithTime(prod.OpenTime,comm.PRODUCT_YYG_BUY_CODES)
	codeCount,_ :=dao.OrderItemPurchaseCodesWithProdId(prod.ProdId)//商品份数
	var c  int64 = 0
	for _, tSum := range orderItem {
		s1	:=fmt.Sprintf("%s%s",time.Unix(int64(tSum.BuyTime/1e3), 0).Format("150405"),dao.Right(fmt.Sprintf("%d",tSum.BuyTime),3))
		i1,_:=strconv.ParseInt(s1,10,64)			
		c	 =c+i1
	}
	openCode:=fmt.Sprintf("%d",c%codeCount+10000001)	
	user,err:=dao.GetOpenIdbyOpenCode(prod.ProdId,openCode)		
	if err!=nil {
		log.Error(err)
		return
	}
	if user==nil {
		return
	}
	dao.ProductAndPurchaseCodesOpened(prod.ProdId,user.OpenId,user.Mobile,openCode)
}
























