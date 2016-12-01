package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"shopapi/comm"
	"net/http"
	"shopapi/dao"
	"strings"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"gitlab.qiyunxin.com/tangtao/utils/page"
	"fmt"
	"time"
)

type OrderDto struct  {
	Items []OrderItemDto `json:"items"`
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

type OrderDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	PayapiNo string `json:"payapi_no"`
	OpenId string `json:"open_id"`
	Name string `json:"name"`
	Mobile string `json:"mobile"`
	AddressId int64 `json:"address_id"`
	Address string `json:"address"`
	AddressName string `json:"address_name"`
	AddressMobile string `json:"address_mobile"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	Price float64 `json:"price"`
	RealPrice float64 `json:"real_price"`
	PayPrice float64 `json:"pay_price"`
	OmitMoney float64 `json:"omit_money"`
	RejectCancelReason string `json:"reject_cancel_reason"`
	CancelReason string `json:"cancel_reason"`
	OrderStatus int `json:"order_status"`
	PayStatus int `json:"pay_status"`
	Items []*OrderItemDetailDto `json:"items"`
	Json string `json:"json"`
	CreateTime string `json:"create_time"`

	GmOrdernum string  `json:"ordernum"`
	GmPassnum string   `json:"passnum"`
	GmPassway string `json:"passway"`
	WayStatus int64 `json:"way_status"`
	
	Yyg	*OrderItemYygDto `json:"yyg"`
}
type OrderItemYygDto struct  {
	OpenCode   string `json:"open_code"`
	OpenStatus int64 `json:"open_status"`
	OpenMobile string `json:"open_mobile"`
	BuyCodes   []string `json:"buy_codes"`
}
type OrderItemDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	//商户名称
	MerchantName string `json:"merchant_name"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//商品cover 封面图 url
	ProdCoverImg string `json:"prod_coverimg"`
	ProdTitle string `json:"prod_title"`
	//商品标识
	ProdFlag string `json:"prod_flag"`
	ProdId int64 `json:"prod_id"`
	Num int `json:"num"`
	OfferUnitPrice float64 `json:"offer_unit_price"`
	OfferTotalPrice float64 `json:"offer_total_price"`
	BuyUnitPrice float64 `json:"buy_unit_price"`
	BuyTotalPrice float64 `json:"buy_total_price"`
	Json string `json:"json"`
	
	GmOrdernum string  `json:"ordernum"`
	GmPassnum string   `json:"passnum"`
	GmPassway string `json:"passway"`
	WayStatus int64 `json:"way_status"`

}

type OrderItemDto struct  {
	//分销编号 (推荐人的分销编号)
	DbnNo string `json:"dbn_no"`
	//sku编号
	SkuNo string `json:"sku_no"`
	//商品数量
	Num int `json:"num"`

	Json string `json:"json"`

}

type OrderCountWithUserAndStatus struct {
	Count int64 `json:"count"`
	Status string `json:"status"`
}

//添加订单
func OrderAdd(c *gin.Context)  {
	appId,err :=CheckAppAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	openId,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	var orderDto OrderDto
	err =c.BindJSON(&orderDto)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	if orderDto.MerchantId== 0 {
		log.Error(err)
		util.ResponseError400(c.Writer,"商户ID不能为空!")
		return
	}

	if openId == "" {
		util.ResponseError400(c.Writer,"open_id不能为空!")
		return
	}
	orderDto.AppId = appId
	orderDto.OpenId = openId
	orderDto.OrderStatus = comm.ORDER_STATUS_WAIT_SURE
	
	//限购商品
	items := orderDto.Items
	if items!=nil {
		prodSkuDetail := dao.NewProdSkuDetail()
		var prodSkuDetailItem *dao.ProdSkuDetail
		var buyCount int64
		for _,item :=range items  {		
			prodSkuDetailItem,_ =prodSkuDetail.WithSkuNo(item.SkuNo)
			if prodSkuDetailItem.LimitNum>0 {
				buyCount,_=service.ProdOrderCountWithId(prodSkuDetailItem.ProdId,orderDto.OpenId,fmt.Sprintf(time.Now().Format("2006-01-02")))
				if (buyCount+int64(item.Num))>prodSkuDetailItem.LimitNum {
					util.ResponseError400(c.Writer,"您购买的商品数量已超过限购数!")
					return
				}
			}
		}
	}	
	
	order,err := service.OrderAdd(orderDtoToModel(orderDto))
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	orderDto.OrderNo = order.No

	c.JSON(http.StatusOK,orderDto)
}


//预支付订单
func OrderPrePay(c *gin.Context)  {

	appId,err :=CheckAppAuth(c)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}

	var params service.OrderPrePayDto
	err =c.BindJSON(&params)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	if params.AddressId==0 {
		util.ResponseError400(c.Writer,"地址ID不能为空!")
		return
	}

	orderNo := c.Param("order_no")
	params.OrderNo = orderNo
	params.AppId = appId
	
	resultMap,err := service.OrderPrePay(service.OrderPrePayDtoToModel(params))
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,resultMap)
}

//账户余额付款
func OrderPayForAccount(c *gin.Context)  {
	//获取用户openid

	user,err := security.GetAuthUser(c.Request)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败！")
		return
	}
	orderNo := c.Param("order_no")
	if orderNo =="" {
		util.ResponseError400(c.Writer,"订单号不能为空!")
		return
	}
	var paramMap map[string]interface{}
	err = c.BindJSON(&paramMap)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	payToken :=paramMap["pay_token"]
	if payToken==nil{
		util.ResponseError400(c.Writer,"支付token不能为空!")
		return
	}
	appId :=security.GetAppId2(c.Request)
	err =service.OrderPayForAccount(user.OpenId,orderNo,payToken.(string),appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	util.ResponseSuccess(c.Writer)
}
//取消订单
func OrderCancel(c *gin.Context)  {
	_,err := security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	appId := security.GetAppId2(c.Request)
	var paramMap map[string]interface{}
	err = c.BindJSON(&paramMap)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	err =service.OrderCancel(orderNo,paramMap["reason"].(string),appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	util.ResponseSuccess(c.Writer)
}

//订单拒绝取消
func OrderRefuseCancel(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	var params map[string]interface{}
	err =c.BindJSON(&params)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	appId := security.GetAppId2(c.Request)
	err =service.OrderRefuseCancel(orderNo,params["reason"].(string),appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

//确认订单
func OrderSure(c *gin.Context)  {

	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	appId :=security.GetAppId2(c.Request)
	err =service.OrderSure(orderNo,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"订单确认失败!")
		return
	}

	util.ResponseSuccess(c.Writer)
}

//同意取消订单
func OrderAgreeCancel(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	appId := security.GetAppId2(c.Request)
	err = service.OrderAgreeCancel(orderNo,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

//根据编号查询订单信息
func OrderByNo(c *gin.Context)  {

}

//订单详情
func OrderDetailByNo(c *gin.Context)  {
	//appId,err :=CheckAuth(c)
	//if err!=nil{
	//	util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
	//	return
	//}
	////获取用户openid
	//openId,err :=GetOpenId(c)
	//if err!=nil{
	//	log.Error(err)
	//	util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
	//	return
	//}
}


//获取用户订单
func OrderWithUserAndStatus(c *gin.Context)  {
	//分页
	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))

	appId,err :=CheckAppAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	openId,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	orderStatus := c.Query("order_status")
	iorderStatusArray :=make([]int,0)
	if orderStatus!="" {
		orderStatusArray := strings.Split(orderStatus,",")
		if len(orderStatusArray)>0 {
			for _,statusStr :=range orderStatusArray {
				stat,err :=strconv.Atoi(statusStr)
				if err!=nil {
					log.Error(err)
					util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
					return
				}
				iorderStatusArray = append(iorderStatusArray,stat)
			}
		}
	}


	payStatus := c.Query("pay_status")
	ipayStatusArray :=make([]int,0)
	if payStatus!="" {
		payStatusArray := strings.Split(payStatus,",")
		if len(payStatusArray)>0 {
			for _,statusStr :=range payStatusArray {
				stat,err :=strconv.Atoi(statusStr)
				if err!=nil {
					log.Error(err)
					util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
					return
				}
				ipayStatusArray = append(ipayStatusArray,stat)
			}
		}
	}

	orderList,err := service.OrderByUser(openId,iorderStatusArray,ipayStatusArray,appId,pIndex,pSize)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}

	orderDetailDtos :=make([]*OrderDetailDto,0)
	if orderList!=nil{
		for _,orderDetail :=range orderList {
			orderDetailDtos = append(orderDetailDtos,orderDetailToDto(orderDetail))
		}
	}

	c.JSON(http.StatusOK,orderDetailDtos)
}

//获取用户指定状态订单数量
func OrderWithUserAndStatusCount(c *gin.Context)  {
	appId,err :=CheckAppAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	
	//获取用户openid
	openId,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	
	stat,err :=strconv.Atoi(c.Param("status"))
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
		return
	}
	
	//订单状态
	iorderStatusArray :=make([]int,0)
	iorderStatusArray = append(iorderStatusArray,stat)
	
	//订单支付状态
	ipayStatusArray	  :=make([]int,0)
	
	
	orderCount,err := service.OrderWithUserAndStatusCount(openId,iorderStatusArray,ipayStatusArray,appId)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}
	
	count :=&OrderCountWithUserAndStatus{}
	count.Count=orderCount.Count
	if orderCount.Count>0 {
		count.Status="true"
	}else{
		count.Status="false"
	}

	c.JSON(http.StatusOK,count)
}


func OrderDetailWithNo(c *gin.Context)  {

	//获取用户openid
	_,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}

	appId := security.GetAppId2(c.Request)
	orderNo := c.Param("order_no")
	log.Debug(orderNo,appId)
	orderDetail,err := service.OrderDetailWithNo(orderNo,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if orderDetail!=nil {
		c.JSON(http.StatusOK,orderDetailToDto(orderDetail))
	}else{
		util.ResponseError400(c.Writer,"没有找到订单详情!")
	}
}

func OrdersGet(c *gin.Context)  {

	//==search
	var search dao.OrderSearch
	search.MerchantName 	=c.Query("shanghu")
	search.Title 			=c.Query("titile")
	search.OrderNo	 		=c.Query("ordernum")
	search.PayStatus,_ 		=strconv.ParseUint(c.Query("paystate"),10,64)
	search.OrderStatus,_ 	=strconv.ParseUint(c.Query("orderstate"),10,64)
	search.AddressMobile	=c.Query("address_mobile")
	search.Show,_ 			=strconv.ParseUint(c.Query("show"),10,64)
	
	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	appId :=security.GetAppId2(c.Request)
	orders,err :=service.OrdersGet(search,pIndex,pSize,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败！");
		return
	}

	total,err :=service.OrdersGetCount(search,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询订单总数量失败！")
		return
	}

	results 	:=make([]*OrderDto,0)
	detailTitle :=make([]string,0)
	
	if orders!=nil&&len(orders) > 0 {
		var orderItem []*dao.OrderItem
		account := dao.NewAccount()
		
		for _,od :=range orders {
			account,_ =account.AccountWithOpenId(od.OpenId,appId)
			od.Mobile	=account.Mobile
			od.YdgyName	=account.YdgyName
		
			orderItem,_=service.OrderItems(od.No);			
			if len(orderItem)>0 {
				od.GmOrdernum	=orderItem[0].GmOrdernum
				od.GmPassnum	=orderItem[0].GmPassnum
				od.GmPassway	=orderItem[0].GmPassway
				od.WayStatus	=orderItem[0].WayStatus
				for _,odItem :=range orderItem {
					detailTitle=append(detailTitle,fmt.Sprintf("%s*%d", odItem.Title,odItem.Num))
				}
				od.DetailTitle	=detailTitle
				detailTitle=make([]string,0)
			}else{
				log.Info("========")				
			}	
			results = append(results,orderToA(od))
		}
	}
	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(total),results))
}

func MerchantOrders(c *gin.Context)  {
	_,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}

	merchantId := c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}

	orderStatus := c.Query("order_status")
	iorderStatusArray :=make([]int,0)
	if orderStatus!="" {
		orderStatusArray := strings.Split(orderStatus,",")
		if len(orderStatusArray)>0 {
			for _,statusStr :=range orderStatusArray {
				stat,err :=strconv.Atoi(statusStr)
				if err!=nil {
					log.Error(err)
					util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
					return
				}
				iorderStatusArray = append(iorderStatusArray,stat)
			}
		}
	}
	payStatus := c.Query("pay_status")
	ipayStatusArray :=make([]int,0)
	if payStatus!="" {
		payStatusArray := strings.Split(payStatus,",")
		if len(payStatusArray)>0 {
			for _,statusStr :=range payStatusArray {
				stat,err :=strconv.Atoi(statusStr)
				if err!=nil {
					log.Error(err)
					util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
					return
				}
				ipayStatusArray = append(ipayStatusArray,stat)
			}
		}
	}

	appId := security.GetAppId2(c.Request)
	orderList,err :=service.OrderDetailWithMerchantId(imerchantId,iorderStatusArray,ipayStatusArray,appId)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}

	orderDetailDtos :=make([]*OrderDetailDto,0)
	if orderList!=nil{
		for _,orderDetail :=range orderList {
			orderDetailDtos = append(orderDetailDtos,orderDetailToDto(orderDetail))
		}
	}

	c.JSON(http.StatusOK,orderDetailDtos)
}

func orderDetailToDto(model *dao.OrderDetail) *OrderDetailDto {

	dto :=&OrderDetailDto{}
	dto.No = model.No
	dto.RealPrice = model.RealPrice
	dto.PayPrice = model.PayPrice
	dto.AppId = model.AppId
	dto.Id = model.Id
	dto.Json = model.Json
	dto.OpenId = model.OpenId
	dto.OmitMoney = model.OmitMoney
	dto.PayapiNo = model.PayapiNo
	dto.AddressId = model.AddressId
	dto.Address = model.Address
	dto.Price = model.Price
	dto.RejectCancelReason = model.RejectCancelReason
	dto.CancelReason = model.CancelReason
	dto.Name = model.AddressName
	dto.Mobile = model.AddressMobile
	dto.GmOrdernum = model.GmOrdernum
	dto.GmPassnum = model.GmPassnum	
	dto.GmPassway = model.GmPassway	
	dto.WayStatus = model.WayStatus

	dto.OrderStatus = model.OrderStatus
	dto.PayStatus = model.PayStatus
	dto.Title = model.Title
	dto.CreateTime = qtime.ToyyyyMMddHHmm(model.CreateTime)

	items := model.Items
	if items!=nil {
		itemsDto :=make([]*OrderItemDetailDto,0)
		for _,item :=range items {
			itemsDto = append(itemsDto,orderItemDetailToDto(item))
		}
		dto.Items = itemsDto
	}
	
	if model.Itemsyyg!=nil {
		itemsDto :=&OrderItemYygDto{}
		itemsDto.OpenCode   = model.Itemsyyg.OpenCode
		itemsDto.OpenStatus = model.Itemsyyg.OpenStatus
		itemsDto.OpenMobile = model.Itemsyyg.OpenMobile
		itemsDto.BuyCodes   = model.ItemsyygBuyCodes
		dto.Yyg = itemsDto
	}	

	return dto
}

func orderItemDetailToDto(model *dao.OrderItemDetail) *OrderItemDetailDto  {

	dto :=&OrderItemDetailDto{}
	dto.AppId = model.AppId
	dto.OpenId = model.OpenId
	dto.BuyTotalPrice = model.BuyTotalPrice
	dto.BuyUnitPrice = model.BuyUnitPrice
	dto.Id = model.Id
	dto.Json = model.Json
	dto.No = model.No
	dto.Num = model.Num
	dto.OfferTotalPrice = model.OfferTotalPrice
	dto.OfferUnitPrice = model.OfferUnitPrice
	dto.ProdId = model.ProdId
	dto.ProdFlag = model.ProdFlag
	dto.ProdTitle = model.ProdTitle
	dto.ProdCoverImg  = model.ProdCoverImg
	dto.MerchantName = model.MerchantName
	dto.MerchantId = model.MerchantId
	
	dto.GmOrdernum = model.GmOrdernum
	dto.GmPassnum = model.GmPassnum
	dto.GmPassway = model.GmPassway
	dto.WayStatus = model.WayStatus

	return dto
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
	

	return a
}

func orderDtoToModel(dto OrderDto) *service.OrderModel  {

	model := &service.OrderModel{}
	model.AppId = dto.AppId
	model.OpenId = dto.OpenId
	model.Json  =dto.Json
	model.OrderStatus = dto.OrderStatus
	model.MerchantId = dto.MerchantId
	model.MOpenId = dto.MOpenId
	model.PayStatus = dto.PayStatus
	model.Title = dto.Title
	model.AddressId = dto.AddressId

	items := dto.Items
	if items!=nil {
		itemmodels := make([]service.OrderItemModel,0)
		for _,itemDto :=range items  {
			itemmodels = append(itemmodels,orderItemToModel(itemDto))
		}
		model.Items = itemmodels
	}
	return model
}

func orderItemToModel(dto OrderItemDto) service.OrderItemModel  {
	model :=service.OrderItemModel{}
	model.DbnNo = dto.DbnNo
	model.Json = dto.Json
	model.Num = dto.Num
	model.SkuNo = dto.SkuNo
	return model
}

//订单删除
func OrderDelete(c *gin.Context)  {
	_,err := security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	appId := security.GetAppId2(c.Request)
	var paramMap map[string]interface{}
	err = c.BindJSON(&paramMap)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	err =service.OrderDelete(orderNo,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
//订单删除 批量
func OrderDeleteBatch(c *gin.Context)  {
	var err error

	 _,err = security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	
	var orderNo service.OrderNo
	err =c.BindJSON(&orderNo)	
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	}
	
	appId := security.GetAppId2(c.Request)
	/* var paramMap map[string]interface{}
	err = c.BindJSON(&paramMap)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	} */
	
	err =service.OrderDeleteBatch(orderNo,appId)
	if err!=nil{		
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
//单增加购买订单号
func OrdersAddNum(c *gin.Context)  {
	/* _,err := security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	*/
	appId 		:= security.GetAppId2(c.Request)
	//orderNo		:=c.PostForm("id")
	//ordernum	:=c.PostForm("ordernum")
	
	type Params struct {
		Id			string   `json:"id"`
		Ordernum	string   `json:"ordernum"`
	}
	var param Params
	c.BindJSON(&param)
	
	
	err :=service.OrdersAddNum(param.Id,appId,param.Ordernum)
	
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
//订单增加购买运单号
func OrdersAddPassnum(c *gin.Context)  {
	/* _,err := security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	*/
	appId 		:= security.GetAppId2(c.Request)

	//orderNo		:=c.PostForm("id")
	//passnum		:=c.PostForm("passnum")

	type Params struct {
		Id			string   `json:"id"`
		Passnum		string   `json:"passnum"`
		Passway	string   `json:"passway"`
	}
	var param Params
	c.BindJSON(&param)
	
	err :=service.OrdersAddPassnum(param.Id,appId,param.Passnum,param.Passway)
	
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
//订单快递查询
func ExpressDelivery(c *gin.Context)  {
	logisticCode 	:= c.Query("logisticCode")
	shipperCode 	:= c.Query("ShipperCode")
	
	data,err :=service.ExpressDelivery(logisticCode,shipperCode)
	
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	c.JSON(http.StatusOK,data)
}
//changeshowstate
func OrderChangeShowState(c *gin.Context)  {
	appId,err :=CheckAppAuth(c)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}	
	
	no:=c.Param("no")
	/* id,err :=strconv.ParseInt(c.Param("id"),10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"id格式错误")
		return
	} */
	
	show,err :=strconv.ParseInt(c.Param("show"),10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数格式错误")
		return
	}	
	
	err =service.OrderChangeShowState(appId,no,show)	
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}







