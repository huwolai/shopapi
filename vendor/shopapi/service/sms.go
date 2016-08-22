package service

import (
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"crypto/md5"
	"time"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"errors"
	"encoding/hex"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

const SMS_BASE_URL  = "https://app.cloopen.com:8883"

//发送验证码短信
func SendCodeSMS(mobile string,code string) (error) {


	accountSid := "aaf98f8950f4a62c0151061744fd43b8"
	accountToken :="834b696ff4954450b751a8a4590f1fe0"
	appId :="aaf98f89512446e2015128f0802e0d4e"
	date  :=time.Now()
	dateStr :=date.Format("20060102150405")
	log.Info(dateStr)

	auth :=base64.StdEncoding.EncodeToString([]byte(accountSid+":"+dateStr))
	h:= md5.New()
	h.Write([]byte(accountSid+accountToken+dateStr))
	sign :=h.Sum(nil)
	signStr :=hex.EncodeToString(sign)
	header :=map[string]string{
		"Authorization": auth,
		"Content-Type": "application/json;charset=utf-8",
		"Accept": "application/json",
	}

	param :=map[string]interface{}{
		"to":mobile,
		"appId":appId,
		"templateId":"51954",
		"datas":[]string{code,"qiyunxin"},
	}



	jsonData,_ := json.Marshal(param)
	 resopnse,err :=network.Post(SMS_BASE_URL+"/2013-12-26/Accounts/"+accountSid+"/SMS/TemplateSMS?sig="+signStr,jsonData,header)
	if err!=nil {

		return err
	}

	log.Debug(resopnse.Body)

	if resopnse.StatusCode==http.StatusOK {
		var resutlMap map[string]interface{}
		err :=util.ReadJsonByByte([]byte(resopnse.Body),&resutlMap)
		if err!=nil {
			return err
		}
		if resutlMap["statusCode"].(string)== "000000" {
			return nil
		}

		return errors.New("短信发送错误["+resutlMap["statusCode"].(string)+"]")
	}else{
		return errors.New("请求短信接口失败!")
	}
}


