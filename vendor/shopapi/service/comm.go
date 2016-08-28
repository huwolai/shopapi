package service

import (
	"fmt"
	"time"
	"crypto/md5"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"net/http"
	"errors"
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

func GetCouponSign(params map[string]interface{}) (noncestr,timestamp,appid,basesign,sign string)  {
	appid= config.GetValue("coupon_appid").ToString()
	apikey :=config.GetValue("coupon_appkey").ToString()
	noncestr ="5671223"
	timestamp =fmt.Sprintf("%d",time.Now().Unix())
	signStr := apikey+noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign =fmt.Sprintf("%X",bytes)

	log.Debug("apikey=",apikey)

	sign = util.SignWithBaseSign(params,apikey,basesign,nil)
	return
}

func GetPayAuthHeader(openId string) ( map[string]string) {

	appId := config.GetValue("payapi_appid").ToString()

	header :=map[string]string{
		"open_id":openId,
		"app_id":appId,
	}
	return header
}

type CouponResultDto struct  {
	//交易号
	TradeNo string `json:"trade_no"`
	//优惠金额
	Amount float64 `json:"amount"`
	//优惠编号
	CouponNo string `json:"coupon_no"`
	//标题
	Title string `json:"title"`
	//备注
	Remark string `json:"remark"`
	//标记
	Flag  string `json:"flag"`
	//附加数据
	Json string `json:"json"`

}
//请求优惠API
func RequestCouponApi(url string,params map[string]interface{}) ([]*CouponResultDto,error) {
	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetCouponSign(params)
	log.Info(fmt.Sprintf("%s.%s",basesign,sign))
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(params);

	response,err := network.Post(url,paramData,headers)
	if err!=nil{
		return nil,err
	}
	if response.StatusCode==http.StatusOK {
		var dto []*CouponResultDto
		err =util.ReadJsonByByte([]byte(response.Body),&dto)
		return dto,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}else{
		return nil,errors.New("请求API失败!")
	}
}

func RequestPayApi(path string,params map[string]interface{}) (map[string]interface{},error) {
	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetPayapiSign(params)
	log.Info(fmt.Sprintf("%s.%s",basesign,sign))
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(params);

	response,err := network.Post(config.GetValue("payapi_url").ToString()+path,paramData,headers)
	if err!=nil{
		return nil,err
	}

	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)


		return resultMap,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}else{
		return nil,errors.New("请求API失败!")
	}
}