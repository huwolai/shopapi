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
)

type OrderDto struct  {
	Items []OrderItemDto `json:"items"`
	Json string `json:"json"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	MOpenId string
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	AddressId int64 `json:"address_id"`
	Title string `json:"title"`
	OrderNo string `json:"order_no"`
	OrderStatus int `json:"order_status"`
	PayStatus int `json:"pay_status"`
}

type OrderDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	PayapiNo string `json:"payapi_no"`
	OpenId string `json:"open_id"`
	AddressId int64 `json:"address_id"`
	Address string `json:"address"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	ActPrice float64 `json:"act_price"`
	OmitMoney float64 `json:"omit_money"`
	Price float64 `json:"price"`
	OrderStatus int `json:"order_status"`
	PayStatus int `json:"pay_status"`
	Items []*OrderItemDetailDto `json:"items"`
	Json string `json:"json"`
	CreateTime string `json:"create_time"`

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
	ProdId int64 `json:"prod_id"`
	Num int `json:"num"`
	OfferUnitPrice float64 `json:"offer_unit_price"`
	OfferTotalPrice float64 `json:"offer_total_price"`
	BuyUnitPrice float64 `json:"buy_unit_price"`
	BuyTotalPrice float64 `json:"buy_total_price"`
	Json string `json:"json"`

}

type OrderItemDto struct  {
	//sku编号
	SkuNo string `json:"sku_no"`
	//商品数量
	Num int `json:"num"`

	Json string `json:"json"`

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
	}
	orderDto.AppId = appId
	orderDto.OpenId = openId
	orderDto.OrderStatus = comm.ORDER_STATUS_WAIT_SURE

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
	openId,err :=CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	orderNo := c.Param("order_no")
	if orderNo =="" {
		util.ResponseError400(c.Writer,"订单号不能为空!")
		return
	}

	appId :=security.GetAppId2(c.Request)
	err =service.OrderPayForAccount(openId,orderNo,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

//取消订单
func OrderCancel(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	orderNo :=c.Param("order_no")
	appId := security.GetAppId2(c.Request)

	err =service.OrderCancel(orderNo,appId)
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

	orderList,err := service.OrderByUser(openId,iorderStatusArray,ipayStatusArray,appId)
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
	dto.ActPrice = model.ActPrice
	dto.AppId = model.AppId
	dto.Id = model.Id
	dto.Json = model.Json
	dto.OpenId = model.OpenId
	dto.OmitMoney = model.OmitMoney
	dto.PayapiNo = model.PayapiNo
	dto.AddressId = model.AddressId
	dto.Address = model.Address
	dto.Price = model.Price
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

	return dto
}

func orderItemDetailToDto(model *dao.OrderItemDetail) *OrderItemDetailDto  {

	dto :=&OrderItemDetailDto{}
	dto.AppId = model.AppId
	dto.BuyTotalPrice = model.BuyTotalPrice
	dto.BuyUnitPrice = model.BuyUnitPrice
	dto.Id = model.Id
	dto.Json = model.Json
	dto.No = model.No
	dto.Num = model.Num
	dto.OfferTotalPrice = model.OfferTotalPrice
	dto.OfferUnitPrice = model.OfferUnitPrice
	dto.ProdId = model.ProdId
	dto.ProdTitle = model.ProdTitle
	dto.ProdCoverImg  = model.ProdCoverImg
	dto.MerchantName = model.MerchantName
	dto.MerchantId = model.MerchantId

	return dto
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
	model.Json = dto.Json
	model.Num = dto.Num
	model.SkuNo = dto.SkuNo
	return model
}