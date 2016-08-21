package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"os"
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/api"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, app_id, open_id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	if os.Getenv("GO_ENV")=="" {
		os.Setenv("GO_ENV","tests")
		os.Setenv("APP_ID","shopapi")
	}

	err :=config.Init(false)
	util.CheckErr(err)
	err = startup.InitDBData()
	env := os.Getenv("GO_ENV")
	if env=="tests" {
		gin.SetMode(gin.DebugMode)
	}else if env== "production" {
		gin.SetMode(gin.ReleaseMode)
	}else if env == "preproduction" {
		gin.SetMode(gin.TestMode)
	}

	router := gin.Default()

	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")
	{
		comm :=v1.Group("/comm")
		{
			comm.POST("/images/upload",api.ImageUpload)
		}


		//用户
		users :=v1.Group("/users")
		{
			users.POST("/loginSMS",api.LoginForSMS)
		}
		user :=v1.Group("/user")
		{
			user.POST("/:open_id/merchantupdate",api.MerchantUpdate)
			user.POST("/:open_id/merchant",api.MerchantAdd)
			user.GET("/:open_id/merchant",api.MerchantWithOpenId)
			user.POST("/:open_id/recharge",api.AccountPreRecharge)
			user.POST("/:open_id/account",api.AccountDetail)
		}
		//分类
		categories := v1.Group("/categories")
		{
			categories.GET("/",api.CategoryWithFlags)
		}
		//分类
		category := v1.Group("/category")
		{
			category.GET("/:category_id/products",api.ProductListWithCategory)
		}

		//商品
		products :=v1.Group("/products")
		{
			products.POST("/:merchant_id",api.ProductAdd)
		}
		//商品
		product :=v1.Group("/product")
		{

			product.GET("/:prod_id/imgs",api.ProdImgsWithProdId)
			product.GET("/:prod_id/detail",api.ProdDetailWithProdId)
			//通过属性生成sku (特殊接口)
			product.POST("/:prod_id/sku",api.ProductAndAttrAdd)
			//通过属性key查询商品属性值
			product.GET("/:prod_id/attr/:attr_key",api.ProductAttrValues)

		}
		merchants :=v1.Group("/merchants")
		{
			//附近商户
			merchants.GET("/nearby",api.MerchatNear)
			//商户图片
			merchants.GET("/user/:open_id/imgs",api.MerchantImgWithFlag)
		}

		//商户
		merchant :=v1.Group("/merchant")
		{
			//商户商品
			merchant.GET("/:merchant_id/prods",api.MerchantProds)
			//商户审核
			merchant.POST("/:merchant_id/audit",api.MerchantAudit)
		}

		//推荐
		recom :=v1.Group("/recom")
		{
			//商品推荐列表
			recom.GET("/products",api.ProductListWithRecomm)
		}
		//订单
		order := v1.Group("/order")
		{
			//添加订单
			order.POST("/",api.OrderAdd)
			//订单预支付
			order.POST("/:order_no/prepay",api.OrderPrePay)
			//订单支付
			order.POST("/:order_no/pay",api.OrderPayForAccount)
			order.GET("/detail/:order_no",api.OrderDetailWithNo)
			order.GET("/status/:status",api.OrderWithUserAndStatus)
			order.POST("/:order_no/event",api.OrderEventPost)
		}

		pay := v1.Group("/pay")
		{
			pay.POST("/payapi/callback",api.CallbackForPayapi)
		}

		address := v1.Group("/address")
		{
			address.GET("/:open_id",api.AddressWithOpenId)
			address.POST("/:open_id",api.AddressAdd)
			address.PUT("/:open_id",api.AddressUpdate)
			address.DELETE("/:open_id/id/:id",api.AddressDelete)
			address.GET("/:open_id/id/:id",api.AddressWithId)
			address.GET("/:open_id/recom",api.AddressWithRecom)
		}
	}

	router.Static("/upload","./config/upload")
	router.Run(":8080")
}
