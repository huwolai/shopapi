package comm

const (
	//产品正常状态
	PRODUCT_STATUS_NORMAL = 1
)

const (
	//商户正常状态
	MERCHANT_STATUS_NORMAL = 1
)

const (
	//订单等待预支付
	ORDER_STATUS_PAY_WAIT = 1

	//订单付款成功
	ORDER_STATUS_PAY_SUCCESS=2
)

const (
	//账户状态正常
	ACCOUNT_STATUS_NORMAL =1

	//账户待绑定支付
	ACCOUNT_STATUS_WAIT_BINDPAY =2

	//账户已锁
	ACCOUNT_STATUS_LOCK =0
)

//交易类型
const  (

	//充值订单
	Trade_Type_Recharge = 1
	//购买
	Trade_Type_Buy =2
	//预付款
	Trade_Type_Imprest =3
)

const (
	//支付宝支付
	Pay_Type_Alipay = 1
	//微信支付
	Pay_Type_WX = 2
	//现金支付
	Pay_Type_Cash = 3
	//账户余额支付
	Pay_Type_Account =4
)

//商户默认商品标记
const MERCHANT_DEFAULT_PRODUCT_FLAG  = "merchant_default"