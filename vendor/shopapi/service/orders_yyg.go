package service

import (
	"shopapi/dao"	
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
	//"fmt"
	//"strings"
)

type ProdPurchaseCode struct  {
	Id 			int64  `json:"id"`
	ProdId		int64  `json:"prod_id"`
	Sku 		string `json:"sku"`
	OpenStatus 	int64  `json:"open_status"`	
	PordTitle 	string `json:"pord_title"`	
	OrderDto	 	
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
	
	//log.Info(prods==nil)
	
	itemsDto :=make([]*ProdPurchaseCode,0)
	orderDao :=dao.NewOrder()
	prodDao	 :=dao.NewProduct()
	var orderItem []*dao.OrderItem
	//account := dao.NewAccount()
	for _,item :=range prods {
		dto :=&ProdPurchaseCode{}		
		dto.Id			=item.Id
		dto.ProdId		=item.ProdId
		dto.Sku			=item.Sku
		dto.OpenStatus	=item.OpenStatus
		//OpenCode
	
		
		prod,_:=prodDao.ProductWithId(item.ProdId,appId)
		if prod!=nil {		
			dto.PordTitle 			= prod.Title
		}		
		
		if item.OpenStatus==2 {
			order,err:=orderDao.OrderWithPordYyg(item.OpenCode)
			if err!=nil {		
				return nil,0,err
			}
			if order!=nil {
				//account,_ =account.AccountWithOpenId(order.OpenId,appId)
				//order.Mobile	=account.Mobile
				//order.YdgyName	=account.YdgyName
				
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

























