package comm

const (
	//产品正常状态
	PRODUCT_STATUS_NORMAL = 1
	//参与计算一元购产品计算中奖号的条数
	PRODUCT_YYG_BUY_CODES = 50
)

const (

	//等待审核
	MERCHANT_STATUS_WAIT_AUIT = 0
	//商户正常状态 1/5
	MERCHANT_STATUS_NORMAL = 1 
	
	//记录为null 或者申请中  2
	MERCHANT_STATUS_APPLICATION = 2
	
	//不通过为3
	MERCHANT_STATUS_FAIL = 3
	
	//准厨师(能看到厨师功能界面,但是不显示在首页) 4
	MERCHANT_STATUS_CHEF_ING = 4
	
	//后期用户更新资料中 6
	MERCHANT_STATUS_UPDATE = 6

	
)

const (
	//没有付款
	ORDER_PAY_STATUS_NOPAY = 0
	//付款中
	ORDER_PAY_STATUS_PAYING = 2
	//付款成功
	ORDER_PAY_STATUS_SUCCESS=1

	//等待确认
	ORDER_STATUS_WAIT_SURE = 0
	//已确认
	ORDER_STATUS_SURED = 1
	//已取消
	ORDER_STATUS_CANCELED = 2
	//无效
	ORDER_STATUS_INVALID = 3
	//退货
	ORDER_STATUS_REJECTED = 4
	//订单取消等待确认
	ORDER_STATUS_CANCELED_WAIT_SURE = 5
	//订单取消拒绝确认
	ORDER_STATUS_CANCELED_REJECTED = 6
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

const (
	//修改支付付款密码前缀
	CODE_PAYPWD_PREFIX = "pay_code_"
	//验证码过期时间
	CODE_PAYPWD_EXPIRE  = 60*5
)

const (
	//账户充值等待完成
	ACCOUNT_RECHARGE_STATUS_WAIT = 0
	//账户充值正常状态
	ACCOUNT_RECHARGE_STATUS_NORMAL = 1

)

const (
	//优惠券未使用
	COUNPON_STATUS_NOUSE = 0
	//优惠券已使用
	COUNPON_STATUS_USED = 1
)

const (
	//未激活
	ORDER_COUPON_STATUS_UNACTIVATE = 0

	//已激活
	ORDER_COUPON_STATUS_ACTIVATED = 1
)

//商户默认商品标记
const MERCHANT_DEFAULT_PRODUCT_FLAG  = "merchant_default"

const (
	//收藏厨师
	FAVORITES_TYPE_MERCHANT = 1
	//收藏商品
	FAVORITES_TYPE_PRODUCT =2
)

//管理员
const (
	ADMINOPENID  = "wesdfsfsdf23323"
	ADMINPWD	 = "180181"
	KEFUOPENID	 = "94d4f8e55ca8413ea0d93f6eda926b54"
)

















