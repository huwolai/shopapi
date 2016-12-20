package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"os"
	"shopapi/api"
	"shopapi/task"
	
	//"gitlab.qiyunxin.com/tangtao/utils/security"
	//"gitlab.qiyunxin.com/tangtao/utils/app"
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

	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "tests")
		os.Setenv("APP_ID", "shopapi")
	}

	err := config.Init(false)
	util.CheckErr(err)
	err = startup.InitDBData()
	util.CheckErr(err)
	env := os.Getenv("GO_ENV")
	if env == "tests" {
		gin.SetMode(gin.DebugMode)
	} else if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if env == "preproduction" {
		gin.SetMode(gin.TestMode)
	}
	/* //应用安装
	app.Setup()

	//初始化权限资源
	security.InitSources([]security.Source{
		security.Source{Id:"user",Name:"用户",Permissions:"create,update,delete,browse"},
		security.Source{Id:"product",Name:"产品",Permissions:"create,update,delete,browse"},
		security.Source{Id:"merchant",Name:"商户",Permissions:"create,update,delete,browse"},
	})

	//安装权限控制功能
	security.Setup() */

	//安装消息队列
	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	//开启定时器
	task.StartCron ()

	router := gin.Default()

	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")
	{
		comm := v1.Group("/comm")
		{
			comm.POST("/images/upload", api.ImageUpload)
		}
		//用户
		users := v1.Group("/users")
		{
			users.POST("/loginSMS", api.LoginForSMS)
			//配置登入界面
			users.GET("/getonkey", api.GetOnKey)

			users.GET("/",api.AccountsGet)
			//充值记录 后台
			users.GET("/rechargerecord/admin", api.RechargeRecordByAdmins)
			//获取全部用户总额
			users.GET("/get_user_money/admin", api.GetUsersMoney)
		}
		user := v1.Group("/user")
		{
			user.GET("/:open_id/mobile/:mobile/sms", api.PayPwdUpdateSMS)
			user.PUT("/:open_id/mobile/:mobile/paypwd", api.PayPwdUpdate)
			user.POST("/:open_id/merchantupdate", api.MerchantUpdate)
			user.POST("/:open_id/merchant", api.MerchantAdd)
			user.GET("/:open_id/merchant", api.MerchantWithOpenId)
			user.POST("/:open_id/recharge", api.AccountPreRecharge)
			user.POST("/:open_id/rechargebyadmin", api.AccountPreRechargeByAdmin)
			user.POST("/:open_id/rechargebyadminok", api.AccountPreRechargeByAdminOK)
			//用户账户信息
			user.POST("/:open_id/account", api.AccountDetail)
			//用户订单
			user.GET("/:open_id/orders", api.OrderWithUserAndStatus)
			user.GET("/:open_id/orders/status/:status",api.OrderWithUserAndStatusCount)
			//用户在线状态
			user.GET("/:open_id/merchant/online", api.MerchantOnline)
			user.GET("/:open_id/merchant/online/:status", api.MerchantOnlineAndChange)
			//用户收藏
			user.POST("/:open_id/favorites", api.FavoritesAdd)
			user.GET("/:open_id/favorites", api.FavoritesGet)
			user.DELETE("/:open_id/favorites/:id", api.FavoritesDelete)
			user.GET("/:open_id/existfavorites", api.FavoritesIsExist)

			//用户添加分销
			user.POST("/:open_id/distribution/:id", api.DistributionProductAdd)
			//用户取消分销
			user.DELETE("/:open_id/distribution/:id", api.DistributionProductCancel)
			//添加银行信息
			user.POST("/:open_id/banks", api.UserBankAdd)
			//查询用户银行信息
			user.GET("/:open_id/banks", api.UserBankGet)
			//修改用户银行信息
			user.PUT("/:open_id/banks", api.UserBankUpdate)
			//删除银行信息
			user.DELETE("/:open_id/bank/:id", api.UserBankDel)
			//是否是厨师
			user.GET("/:open_id/ismerchant", api.MerchantIs)	
			//充值记录 后台
			user.GET("/:open_id/rechargerecord/admin", api.RechargeRecordByAdmin)

		}
		//分类
		categories := v1.Group("/categories")
		{
			categories.GET("/", api.CategoryWithFlags)
		}
		//分类
		category := v1.Group("/category")
		{
			category.GET("/:category_id/products", api.ProductListWithCategory)
			category.GET("/:category_id/products/islimit", api.ProductListWithCategoryIsLimit)
		}

		skus := v1.Group("/skus")
		{
			skus.POST("/product/:prod_id", api.ProdSkuAdd)
			// 修改sku
			skus.PUT("/product/:prod_id", api.ProdSkuUpdate)
		}
		//商品
		products := v1.Group("/products")
		{
			products.POST("/:merchant_id", api.ProductAdd)
			products.GET("/", api.ProdDetailListWith)				
		}

		//商品
		product := v1.Group("/product")
		{			
			//根据属性path查询商品的SKU
			product.GET("/:prod_id/sku", api.ProductSkuWithProdIdAndSymbolPath)
			product.GET("/:prod_id/imgs", api.ProdImgsWithProdId)
			product.GET("/:prod_id/detail", api.ProdDetailWithProdId)
			product.GET("/:prod_id/detail/yyg",api.ProdDetailYygWithProdId)
			//修改SKU 库存
			product.POST("/:prod_id/updatestock", api.ProductUpdateStockWithProdId)
			//通过属性生成sku (特殊接口)
			product.POST("/:prod_id/sku", api.ProductAndAttrAdd)
			//一元购生成购买码
			product.POST("/:prod_id/purchase_codes", api.ProductAndPurchaseCodesAdd)
			//product.POST("/:prod_id/purchase_codes/minus", api.ProductAndPurchaseCodesMinus)
			product.GET("/:prod_id/buy_codes", api.ProductBuyCodesWithProdId)
			//通过属性key查询商品属性值
			product.GET("/:prod_id/attr/:attr_key", api.ProductAttrValues)
			//修改商品状态
			product.PUT("/:prod_id/status/:status", api.ProductStatusUpdate)

			//设置商品推荐 0.取消推荐 1.推荐
			product.PUT("/:prod_id/recom/:is_recom", api.ProductRecom)
			//修改商品
			product.PUT("/:prod_id",api.ProductUpdate)
			//录入商品链接
			product.POST("/:prod_id/addlink", api.ProductAndAddLink)
			//重置重新上架的商品数据
			product.POST("/:prod_id/initpro", api.ProductInitPro)
		}
		merchants := v1.Group("/merchants")
		{
			//附近商户
			merchants.GET("/nearby", api.MerchatNear)
			//附近商户搜索 可提供服务的厨师
			merchants.GET("/nearbySearch", api.MerchatNearSearch)
			//商户图片
			merchants.GET("/user/:open_id/imgs", api.MerchantImgWithFlag)
			//商户菜品图片批量命名
			merchants.POST("/imgsnamed", api.MerchantImgsWithNamed)
			//厨师面试登记表
			merchants.POST("/resume/add", api.MerchantResumesWithAdd)
		}

		//商户
		merchant := v1.Group("/merchant")
		{
			//添加商户服务时间
			merchant.POST("/:merchant_id/servicetimes", api.MerchantServiceTimeAdd)
			//查询商户服务时间
			merchant.GET("/:merchant_id/servicetimes", api.MerchantServiceTimeGet)
			//商户营业时间信息
			merchant.GET("/:merchant_id/open", api.MerchantOpenWithMerchantId)
			//商户图片
			merchant.GET("/:merchant_id/imgs", api.MerchantImgWithMerchantId)
			//商户资料
			merchant.GET("/:merchant_id", api.MerchantWithId)
			merchant.GET("/:merchant_id/distance", api.MerchantWithIdDistance)
			//商户商品
			merchant.GET("/:merchant_id/prods", api.MerchantProds)
			//商户订单
			merchant.GET("/:merchant_id/orders", api.MerchantOrders)
			//商户审核
			merchant.POST("/:merchant_id/audit", api.MerchantAudit)
			//商户分销商品
			merchant.GET("/:merchant_id/distributions/products", api.DistributionWithMerchant)
		}

		//推荐
		recom := v1.Group("/recom")
		{
			//商品推荐列表
			recom.GET("/products", api.ProductListWithRecomm)
		}

		//订单
		orders := v1.Group("/orders")
		{
			//订单快递查询
			orders.GET("/expressdelivery", api.ExpressDelivery)
			//订单删除
			orders.POST("/to_delete", api.OrderDeleteBatch)
		}
		order := v1.Group("/order")
		{
			//添加订单
			order.POST("/", api.OrderAdd)
			//订单预支付
			order.POST("/:order_no/prepay", api.OrderPrePay)
			//订单支付
			order.POST("/:order_no/pay", api.OrderPayForAccount)
			order.GET("/:order_no/detail", api.OrderDetailWithNo)
			//取消订单
			order.PUT("/:order_no/cancel", api.OrderCancel)
			//商户同意取消订单
			order.PUT("/:order_no/agree_cancel", api.OrderAgreeCancel)
			//商户拒绝取消订单
			order.PUT("/:order_no/refuse_cancel", api.OrderRefuseCancel)
			//订单确认
			order.PUT("/:order_no/sure", api.OrderSure)
			//订单删除
			order.DELETE("/:order_no/to_delete", api.OrderDelete)
		}

		pay := v1.Group("/pay")
		{
			pay.POST("/payapi/callback", api.CallbackForPayapi)
		}

		address := v1.Group("/address")
		{
			address.GET("/:open_id", api.AddressWithOpenId)
			address.POST("/:open_id", api.AddressAdd)
			address.PUT("/:open_id", api.AddressUpdate)
			address.DELETE("/:open_id/id/:id", api.AddressDelete)
			address.GET("/:open_id/id/:id", api.AddressWithId)
			address.GET("/:open_id/recom", api.AddressWithRecom)
		}

		//分销
		distributions := v1.Group("/distributions")
		{
			//分销的商品
			distributions.GET("/", api.DistributionWith)
			//参与分销的商品
			distributions.GET("/products", api.DistributionProducts)

			//添加或修改分销(有ID为修改 无ID为添加)
			distributions.POST("/", api.ProductJoinOrUpdateDistribution)
		}
		distribution := v1.Group("/distribution")
		{
			//分销商品信息
			distribution.GET("/:id", api.DistributionProductWithId)
			distribution.DELETE("/:id", api.DistributionProductDelete)

		}
		//建议
		suggests := v1.Group("/suggests")
		{
			suggests.POST("/", api.SuggestAdd)
		}

		admin := v1.Group("/admin")
		{
			admin.GET("/orders", api.OrdersGet)			
			
			admin.POST("/orders/addordernum", api.OrdersAddNum)
			admin.POST("/orders/addpassnum", api.OrdersAddPassnum)
			
			admin.GET("/merchants", api.MerchantWith)
			admin.GET("/product/:prod_id/detail", api.ProdDetailWithProdId)	

			admin.GET("/yygorders", api.OrdersYygWin) //中奖管理
			admin.GET("/yygdetail/:prod_id", api.OrdersYygRecord)
		}
		flags := v1.Group("/flags")
		{
			flags.GET("/", api.FlagsWithTypes)
			flags.GET("/json", api.FlagsGetJsonWithTypes)
			flags.PUT("/json/set/:type", api.FlagsSetJsonWithTypes)
		}
		
		//数据初始化
		init := v1.Group("/inits")
		{
			//商品初始化售出数量
			//init.GET("/products/init", api.ProductInitNum)
			//商品 售出数量 定时增加
			//init.GET("/products/add", api.ProductAddNum)
			//判断token是否过期
			init.GET("/token/expired", api.TokenWithExpired)
		}

		//权限资源数据 （此接口提供给,权限资源服务调用）
		v1.GET("/sources",api.SourcesAll)
		
		//一点公益
		ydgy := v1.Group("/ydgy")
		{
			ydgy.POST("/:open_id/setid/:id", api.YdgySetId)
			ydgy.GET("/:open_id/getid", api.YdgyGetId)
			ydgy.POST("/:open_id/setidwithstatus/:status", api.YdgySetIdWithStatus)
			ydgy.DELETE("/:open_id/deleteid", api.YdgySetIdWithDelete)
		}
		//应用
		app := v1.Group("/apps")
		{
			app.GET("/update_log", api.AppUpdateLog)
		}
		//应用
		changeshowstate := v1.Group("/changeshowstate")
		{
			changeshowstate.GET("/order/:no/:show", api.OrderChangeShowState)
			changeshowstate.GET("/product/:id/:show", api.ProductChangeShowState)
		}
		//购物车
		cart := v1.Group("/cart")
		{
			cart.GET("/:open_id", api.CartList)
			cart.POST("/:open_id/add", api.CartAddToList)
			cart.POST("/:open_id/minus", api.CartMinusFromList)
			cart.POST("/:open_id/update", api.CartUpdateList)
			cart.POST("/:open_id/delete/:id", api.CartDelFromList)
		}
		test := v1.Group("/test")
		{
			test.GET("/", api.Test)
		}
	}

	//设置上传目录
	uploadRootDir := config.GetValue("upload_root_dir").ToString()
	if uploadRootDir == "" {
		uploadRootDir = "./config/upload"
	}

	router.Static("/upload", uploadRootDir)
	router.Run(":8080")
}
