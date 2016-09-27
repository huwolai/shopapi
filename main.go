package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"os"
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/api"
	"shopapi/task"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
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

//认证中间件
//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token :=security.GetParamInRequest("Authorization",c.Request)
//		if token=="" {
//			log.Error("没有认证信息!")
//			c.AbortWithStatus(401)
//			return
//		}
//		jwttoken,err :=security.InitJWTAuthenticationBackend().FetchToken(token)
//		if err!=nil{
//			log.Error(err)
//			c.AbortWithStatus(401)
//			return
//		}
//		if !jwttoken.Valid {
//			log.Error("认证信息无效!")
//			c.AbortWithStatus(401)
//			return
//		}
//		c.Next()
//	}
//}

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

	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	//开启定时器
	task.StartCron()

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
			user.GET("/:open_id/mobile/:mobile/sms",api.PayPwdUpdateSMS)
			user.PUT("/:open_id/mobile/:mobile/paypwd",api.PayPwdUpdate)
			user.POST("/:open_id/merchantupdate",api.MerchantUpdate)
			user.POST("/:open_id/merchant",api.MerchantAdd)
			user.GET("/:open_id/merchant",api.MerchantWithOpenId)
			user.POST("/:open_id/recharge",api.AccountPreRecharge)
			//用户账户信息
			user.POST("/:open_id/account",api.AccountDetail)
			//用户订单
			user.GET("/:open_id/orders",api.OrderWithUserAndStatus)
			//用户收藏
			user.POST("/:open_id/favorites",api.FavoritesAdd)
			user.GET("/:open_id/favorites",api.FavoritesGet)
			user.DELETE("/:open_id/favorites/:id",api.FavoritesDelete)
			user.GET("/:open_id/existfavorites",api.FavoritesIsExist)

			//用户添加分销
			user.POST("/:open_id/distribution/:id",api.DistributionProductAdd)
			//用户取消分销
			user.DELETE("/:open_id/distribution/:id",api.DistributionProductCancel)
			//添加银行信息
			user.POST("/:open_id/banks",api.UserBankAdd)
			//查询用户银行信息
			user.GET("/:open_id/banks",api.UserBankGet)
			//修改用户银行信息
			user.PUT("/:open_id/banks",api.UserBankUpdate)
			//删除银行信息
			user.DELETE("/:open_id/bank/:id",api.UserBankDel)
			//是否是厨师
			user.GET("/:open_id/ismerchant",api.MerchantIs)

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

		skus :=v1.Group("/skus")
		{
			skus.POST("/product/:prod_id",api.ProdSkuAdd)
		}
		//商品
		products :=v1.Group("/products")
		{
			products.POST("/:merchant_id",api.ProductAdd)
			products.GET("/",api.ProdDetailListWith)
		}

		//商品
		product :=v1.Group("/product")
		{
			//根据属性path查询商品的SKU
			product.GET("/:prod_id/sku",api.ProductSkuWithProdIdAndSymbolPath)
			product.GET("/:prod_id/imgs",api.ProdImgsWithProdId)
			product.GET("/:prod_id/detail",api.ProdDetailWithProdId)
			//通过属性生成sku (特殊接口)
			product.POST("/:prod_id/sku",api.ProductAndAttrAdd)
			//通过属性key查询商品属性值
			product.GET("/:prod_id/attr/:attr_key",api.ProductAttrValues)
			//修改商品状态
			product.PUT("/:prod_id/status/:status",api.ProductStatusUpdate)

			//设置商品推荐 0.取消推荐 1.推荐
			product.PUT("/:prod_id/recom/:is_recom",api.ProductRecom)

		}
		merchants :=v1.Group("/merchants")
		{
			//附近商户
			merchants.GET("/nearby",api.MerchatNear)
			//附近商户搜索 可提供服务的厨师
			merchants.GET("/nearbySearch",api.MerchatNearSearch)
			//商户图片
			merchants.GET("/user/:open_id/imgs",api.MerchantImgWithFlag)
		}

		//商户
		merchant :=v1.Group("/merchant")
		{
			//添加商户服务时间
			merchant.POST("/:merchant_id/servicetimes",api.MerchantServiceTimeAdd)
			//查询商户服务时间
			merchant.GET("/:merchant_id/servicetimes",api.MerchantServiceTimeGet)
			//商户营业时间信息
			merchant.GET("/:merchant_id/open",api.MerchantOpenWithMerchantId)
			//商户图片
			merchant.GET("/:merchant_id/imgs",api.MerchantImgWithMerchantId)
			//商户资料
			merchant.GET("/:merchant_id",api.MerchantWithId)
			//商户商品
			merchant.GET("/:merchant_id/prods",api.MerchantProds)
			//商户订单
			merchant.GET("/:merchant_id/orders",api.MerchantOrders)
			//商户审核
			merchant.POST("/:merchant_id/audit",api.MerchantAudit)
			//商户分销商品
			merchant.GET("/:merchant_id/distributions/products",api.DistributionWithMerchant)
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
			order.GET("/:order_no/detail",api.OrderDetailWithNo)
			//取消订单
			order.PUT("/:order_no/cancel",api.OrderCancel)
			//商户同意取消订单
			order.PUT("/:order_no/agree_cancel",api.OrderAgreeCancel)
			//商户拒绝取消订单
			order.PUT("/:order_no/refuse_cancel",api.OrderRefuseCancel)
			//订单确认
			order.PUT("/:order_no/sure",api.OrderSure)
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

		//分销
		distributions := v1.Group("/distributions")
		{
			//分销的商品
			distributions.GET("/",api.DistributionWith)
			//参与分销的商品
			distributions.GET("/products",api.DistributionProducts)

			//添加或修改分销(有ID为修改 无ID为添加)
			distributions.POST("/",api.ProductJoinOrUpdateDistribution)
		}
		//建议
		suggests := v1.Group("/suggests")
		{
			suggests.POST("/",api.SuggestAdd)
		}

		admin :=v1.Group("/admin")
		{
			admin.GET("/merchants",api.MerchantWith)
			admin.GET("/product/:prod_id/detail",api.ProdDetailWithProdId)
		}
		flags :=v1.Group("/flags")
		{
			flags.GET("/",api.FlagsWithTypes)
		}
	}

	//设置上传目录
	uploadRootDir :=config.GetValue("upload_root_dir").ToString()
	if uploadRootDir=="" {
		uploadRootDir = "./config/upload"
	}

	router.Static("/upload",uploadRootDir)
	router.Run(":8080")
}
