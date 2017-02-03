package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	//"gitlab.qiyunxin.com/tangtao/utils/network"
	"shopapi/dao"
	"net/http"
	"fmt"
	"bytes"
	//"io/ioutil"
)
//判断token是否过期
func TokenWithExpired(c *gin.Context)  {
	_,err :=security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError400(c.Writer,"重新登入!")
		return
	}
	util.ResponseSuccess(c.Writer)
}
//厨师面试登记表
func MerchantResumesWithAdd(c *gin.Context)  {
	var resumes dao.MerchantResume
	c.BindJSON(&resumes)
	//resumes.Name= c.PostForm("name")
	
	err:=service.MerchantResumesWithAdd(resumes)	
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}
//menu
func Menu(c *gin.Context)  {
	url	:="http://config.huwolai.com/shopapi/prod/user/%s/acl"
	memu:=[]string{
		"75b3be97d0b749768bcb272e4684786a",
		"46ac15de93d74d828b87ef861300fd0f",
		"e53224bebe264dad9d8c2bc4ad42adac",
		"2bd209e36084479cbbb7258f12fce02f",
		"eee81c2e755a4447a6e58c2447f261c6",
		"75e8f083ce674590bc8380ec0679e742",
		"ebe6a2245e6e4ece888688f81beaa4df",
		"d2a37b2a259d486b825fef7df40b8660",
		"aa711d464b6f4f50a2e9ab5639a4a7b2",
	}
	
	params:=[]map[string]interface{}{
		map[string]interface{}{
			"id"		:"1",
			"name"		:"商品管理",
			"payload"	:"{\"icon\":\"file\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"1001",
			"name"		:"商品管理",			
			"payload"	:"{\"path\":\"/prodmanager\"}",
			"pid"		:"1",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"2",
			"name"		:"用户管理",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"2001",
			"name"		:"用户管理",
			"payload"	:"{\"path\":\"/usermanager\"}",
			"pid"		:"2",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"3",
			"name"		:"商户管理",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"3001",
			"name"		:"商户管理",
			"payload"	:"{\"path\":\"/merchantmanager\"}",
			"pid"		:"3",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"4",
			"name"		:"订单管理",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"4001",
			"name"		:"订单管理",
			"payload"	:"{\"path\":\"/ordermanager\"}",
			"pid"		:"4",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"4002",
			"name"		:"中奖管理",
			"payload"	:"{\"path\":\"/yygmanager\"}",
			"pid"		:"4",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"5",
			"name"		:"分销管理",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"5001",
			"name"		:"分销管理",
			"payload"	:"{\"path\":\"/distribmanager\"}",
			"pid"		:"5",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"6",
			"name"		:"充值管理",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"6001",
			"name"		:"客服充值",
			"payload"	:"{\"path\":\"/kefu\"}",
			"pid"		:"6",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"6002",
			"name"		:"三方充值(APP商城)",
			"payload"	:"{\"path\":\"/sanfang\"}",
			"pid"		:"6",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"6003",
			"name"		:"提现审核",
			"payload"	:"{\"path\":\"/cashout\"}",
			"pid"		:"6",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"7",
			"name"		:"系统设置",
			"payload"	:"{\"icon\":\"team\"}",
			"pid"		:"0",
			"type"		:"1",
		},
		map[string]interface{}{
			"id"		:"7001",
			"name"		:"系统设置",
			"payload"	:"{\"path\":\"/setting\"}",
			"pid"		:"7",
			"type"		:"1",
		},
	}
	
	client := &http.Client{}	
	b,_		:= json.Marshal(params)
	for _,item := range memu {
		req, _ := http.NewRequest("POST", fmt.Sprintf(url,item),bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		
		/* body, _ := ioutil.ReadAll(resp.Body)
		resultMap 	:= map[string]interface{}{}
		util.ReadJsonByByte([]byte(body),&resultMap) */
	}
	util.ResponseSuccess(c.Writer)		
}













































