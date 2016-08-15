package service

import (
	"fmt"
	"time"
	"crypto/md5"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

func NewInOrderNo() string  {

	now := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d%d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second(),now.Nanosecond())
}

func GetPayapiSign(params map[string]interface{}) (noncestr,timestamp,appid,basesign,sign string)  {

	appid= config.GetValue("payapi_appid").ToString()
	apikey :=config.GetValue("payapi_appkey").ToString()
	noncestr ="23435"
	timestamp =fmt.Sprintf("%d",time.Now().Unix())
	signStr := apikey+noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign =fmt.Sprintf("%X",bytes)

	log.Debug("apikey=",apikey)

	sign = util.SignWithBaseSign(params,apikey,basesign,nil)
	return
	
}