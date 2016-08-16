package main

import "testing"
import (
	"github.com/gavv/httpexpect"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
)


const URL  = "http://127.0.0.1:8080/v1"

const APP_ID  =  "23232"

const OPEN_ID  = "099d3c0b061c4c038332cb4e28234c9e"

//添加商品
func TestProductAdd(t *testing.T)  {

	e :=httpexpect.New(t,URL)
	param :=map[string]interface{}{
		"title":"麻辣牛肉面",
		"description":"麻辣牛肉面的描述",
		"category_id": 1,
		"price": 20,
		"dis_price": 19.9,
		"imgnos": "12,23",
	}
	e.POST("/product/1").WithHeader("app_id",APP_ID).WithJSON(param).Expect().
	Status(http.StatusOK).
	JSON().Object().ValueEqual("err_code",0)
}

//成为商户
func TestMerchantAdd(t *testing.T)  {
	e :=httpexpect.New(t,URL)
	param :=map[string]interface{}{
		"name":"商户名称",
		"json":"{\"service\":\"ddddd\"}",
	}
	e.POST("/user/234/merchant").WithHeader("app_id",APP_ID).WithJSON(param).Expect().
	Status(http.StatusOK).
	JSON().Object().ContainsKey("id")
}

//根据分类查询商品列表
func TestProductListWithCategory(t *testing.T)  {
	e :=httpexpect.New(t,URL)

	obj := e.GET("/category/1/products").WithHeader("app_id",APP_ID).Expect().
	Status(http.StatusOK).
	JSON().Array()

	fmt.Println(obj)
}

//查询商品图片
func TestProdImgsWithProdId(t *testing.T)  {
	e :=httpexpect.New(t,URL)

	obj := e.GET("/product/1/imgs").WithHeader("app_id",APP_ID).Expect().
	Status(http.StatusOK).
	JSON().Array()

	fmt.Println(obj)
}

func TestOrderAdd(t *testing.T)  {

	param :=map[string]interface{}{
		"title":"订单标题",
		"json":"{\"service\":\"ddddd\"}",
		"items": []map[string]interface{}{
			map[string]interface{}{
				"prod_id":1,
				"num": 1,
			},
		},
	}

	e :=httpexpect.New(t,URL)
	obj :=e.POST("/order/").WithHeader("app_id",APP_ID).WithHeader("open_id",OPEN_ID).WithJSON(param).Expect().Status(http.StatusOK).JSON().Object()
	fmt.Println(obj)
}

func TestOrderPrePay(t *testing.T)  {
	param :=map[string]interface{}{
		"order_no":"20168159542692233466",
		"pay_type":1,
	}

	e :=httpexpect.New(t,URL)
	obj :=e.POST("/order/20168159542692233466/prepay").WithHeader("app_id",APP_ID).WithHeader("open_id",OPEN_ID).WithJSON(param).Expect().Status(http.StatusOK).JSON().Object()
	fmt.Println(obj)
}

func TestSign(t *testing.T)  {

	param := map[string]interface{}{
		"title":"111",
		"client_ip":"127.0.0.1",
		"notify_url":"http://shopapi.qiyunxin.svc.cluster.local:8080 ",
		"open_id":"123456",
		"out_trade_no":"20168159542692233466",
		"amount":1000,
		"trade_type": 1,
		"pay_type": 1,
	}
	sign :=util.SignWithBaseSign(param,"4537C07A563C4899B5A592DA3CC84A10","9142EB5E2F586C649371244336E341D3",nil)
	fmt.Println(sign)
}

func TestAccountRecharge(t *testing.T)  {
	param :=map[string]interface{}{
		"money":100,
		"pay_type":1,
	}

	e :=httpexpect.New(t,URL)
	obj :=e.POST("/user/"+OPEN_ID+"/recharge").WithHeader("app_id",APP_ID).WithHeader("open_id",OPEN_ID).WithJSON(param).Expect().Status(http.StatusOK).JSON().Object()
	fmt.Println(obj)
}