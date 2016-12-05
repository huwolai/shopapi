package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type ProdPurchaseCode struct  {
	Id 			int64  `json:"id"`
	ProdId		int64  `json:"prod_id"`
	Sku 		string `json:"sku"`
	OpenStatus 	int64  `json:"open_status"`
}

func OrdersYygWin(search dao.OrdersYygSearch,pIndex uint64,pSize uint64) ([]*ProdPurchaseCode,int64,error) {
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
	for _,item :=range prods {
		dto :=&ProdPurchaseCode{}
		dto.Id			=item.Id
		dto.ProdId		=item.ProdId
		dto.Sku			=item.Sku
		dto.OpenStatus	=item.OpenStatus 
		
		itemsDto = append(itemsDto,dto)
	}
	
	return itemsDto,count,nil
}