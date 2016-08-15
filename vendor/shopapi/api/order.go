package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"shopapi/comm"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"net/http"
)

type OrderDto struct  {
	Items []OrderItemDto `json:"items"`
	Json string `json:"json"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	OrderNo string `json:"order_no"`
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

	appId,err :=CheckAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	openId,err :=GetOpenId(c)
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

	appId,err :=CheckAuth(c)
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


//根据编号查询订单信息
func OrderByNo(c *gin.Context)  {
	
}

//订单详情
func OrderDetailByNo(c *gin.Context)  {
	
}


//获取用户订单
func OrderByUser(c *gin.Context)  {
	appId,err :=CheckAuth(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	//获取用户openid
	openId,err :=GetOpenId(c)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	service.Or
}

//post订单事件
func OrderEventPost(c *gin.Context)  {
	
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
	model.Status = comm.ORDER_STATUS_PREPAY_WAIT
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