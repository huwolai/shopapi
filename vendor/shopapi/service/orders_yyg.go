package service

import (
	"shopapi/dao"	
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"fmt"
)

type ProdPurchaseCode struct  {
	Id 			int64  `json:"id"`
	ProdId		int64  `json:"prod_id"`
	Sku 		string `json:"sku"`
	OpenStatus 	int64  `json:"open_status"`	
	PordTitle 	string  `json:"prod_title"`
	Order	 	*OrderDto  `json:"order"`
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
	detailTitle :=make([]string,0)
	var orderItem []*dao.OrderItem
	//account := dao.NewAccount()
	for _,item :=range prods {
		dto :=&ProdPurchaseCode{}		
		dto.Id			=item.Id
		dto.ProdId		=item.ProdId
		dto.Sku			=item.Sku
		dto.OpenStatus	=item.OpenStatus
	
	
		order,err:=orderDao.OrderWithPordYyg(item.ProdId)
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
				
				for _,odItem :=range orderItem {
					detailTitle=append(detailTitle,fmt.Sprintf("%s*%d", odItem.Title,odItem.Num))
				}
				order.DetailTitle	=detailTitle
				detailTitle=make([]string,0)
			}
			dto.Order		=orderToA(order)
			dto.PordTitle	=order.DetailTitle[0]
		}
		itemsDto = append(itemsDto,dto)
	}
	
	return itemsDto,count,nil
}
func OrdersYygRecord(prodId string,appId string,pIndex uint64,pSize uint64) ([]*dao.OrderYyg,error) {
	order,err:=dao.NewOrder().OrdersWithPordYyg(prodId)
	if err!=nil {		
		return nil,err
	}
	return order,nil
}
func orderToA(order *dao.Order) *OrderDto {
	a :=&OrderDto{}
	a.AddressId = order.AddressId
	a.AppId = order.AppId
	a.CancelReason = order.CancelReason
	a.Json = order.Json
	a.MerchantId = order.MerchantId
	a.MerchantName = order.MerchantName
	a.MOpenId = order.MOpenId
	a.OpenId = order.OpenId
	a.OrderNo = order.No
	a.Title = order.Title
	a.PayStatus = order.PayStatus
	a.OrderStatus = order.OrderStatus
	a.RealPrice = order.RealPrice
	a.PayPrice = order.PayPrice
	a.CreateTime = qtime.ToyyyyMMddHHmm(order.CreateTime)
	a.GmOrdernum = order.GmOrdernum
	a.GmPassnum = order.GmPassnum
	a.GmPassway = order.GmPassway
	a.WayStatus = order.WayStatus
	a.DetailTitle = order.DetailTitle
	
	a.Address = order.Address
	a.AddressMobile = order.AddressMobile
	a.AddressName = order.AddressName
	
	a.Show = order.Show
	a.Mobile = order.Mobile
	a.YdgyName = order.YdgyName
	
//	a.YygCode = order.YygCode
	

	return a
}

























