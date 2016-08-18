package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"net/http"
	"shopapi/dao"
	"strings"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
)

type OrderDto struct  {
	Items []OrderItemDto `json:"items"`
	Json string `json:"json"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	OrderNo string `json:"order_no"`
	Status int `json:"status"`
}

type OrderDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	PayapiNo string `json:"payapi_no"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	ActPrice float64 `json:"act_price"`
	OmitMoney float64 `json:"omit_money"`
	Price float64 `json:"price"`
	Status int `json:"status"`
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
	//商品ID
	ProdId int64 `json:"prod_id"`

	//商品数量
	Num int `json:"num"`

	Json string `json:"json"`

}

type OrderPrePayDto struct  {
	OrderNo string `json:"order_no"`
	PayType int `json:"pay_type"`
	AppId string `json:"app_id"`
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

	if openId == "" {
		util.ResponseError400(c.Writer,"open_id不能为空!")
	}
	orderDto.AppId = appId
	orderDto.OpenId = openId
	orderDto.Status = comm.ORDER_STATUS_PAY_WAIT

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

	var params OrderPrePayDto
	err =c.BindJSON(&params)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	orderNo := c.Param("order_no")
	params.OrderNo = orderNo
	params.AppId = appId

	resultMap,err := service.OrderPrePay(orderPrePayDtoToModel(params))
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
	//openId,err :=CheckUserAuth(c)
	//if err!=nil{
	//	log.Error(err)
	//	util.ResponseError400(c.Writer,err.Error())
	//	return
	//}
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

	status := c.Param("status")
	if status=="" {
		util.ResponseError(c.Writer,http.StatusBadRequest,"请输入订单状态!")
		return
	}

	statusArray := strings.Split(status,",")

	istatusArray :=make([]int,0)
	if len(statusArray)>0 {
		for _,statusStr :=range statusArray {
			stat,err :=strconv.Atoi(statusStr)
			if err!=nil {
				util.ResponseError(c.Writer,http.StatusBadRequest,"状态不是数字!")
				return
			}

			istatusArray = append(istatusArray,stat)
		}
	}

	orderList,err := service.OrderByUser(openId,istatusArray,appId)
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
	appId,err :=CheckAppAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	_,err =CheckUserAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

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

//post订单事件
func OrderEventPost(c *gin.Context)  {
	
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
	dto.Price = model.Price
	dto.Status = model.Status
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

func orderPrePayDtoToModel(dto OrderPrePayDto ) *service.OrderPrePayModel  {

	model :=&service.OrderPrePayModel{}
	model.AppId = dto.AppId
	model.OrderNo = dto.OrderNo
	model.PayType = dto.PayType
	model.NotifyUrl = config.GetValue("notify_url").ToString()
	return model
}


func orderDtoToModel(dto OrderDto) *service.OrderModel  {

	model := &service.OrderModel{}
	model.AppId = dto.AppId
	model.OpenId = dto.OpenId
	model.Json  =dto.Json
	model.Status = dto.Status
	model.Title = dto.Title


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
	model.ProdId = dto.ProdId
	return model
}