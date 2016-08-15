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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	os.Setenv("GO_ENV","tests")
	os.Setenv("APP_ID","shopapi")

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
		//用户
		user :=v1.Group("/user")
		{
			user.POST("/:open_id/merchant",api.MerchantAdd)
		}
		//分类
		category := v1.Group("/category")
		{
			category.GET("/:category_id/products",api.ProductListWithCategory)
		}
		//商品
		product :=v1.Group("/product")
		{
			product.POST("/:merchant_id",api.ProductAdd)
			product.GET("/:prod_id/imgs",api.ProdImgsWithProdId)
		}
		//订单
		order := v1.Group("/order")
		{
			order.POST("/",api.OrderAdd)
			order.POST("/:order_no/prepay",api.OrderPrePay)
			order.GET("/:order_no",api.OrderByNo)
			//order.GET("/user/:open_id",api.OrderByUser)
			order.POST("/:order_no/event",api.OrderEventPost)
		}
	}
	router.Run(":8080")
}
