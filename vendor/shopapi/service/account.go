package service

import (
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"net/http"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type AccountRechargeModel struct  {
	//充值账户
	OpenId string
	//充值金额
	Money float64
	PayType int  //付款类型(1.支付宝 2.微信)

}

//账户预充值
func AccountPreRecharge(model *AccountRechargeModel) (map[string]interface{},error) {
	//参数
	params := map[string]interface{}{
		"open_id": model.OpenId,
		"amount": int64(model.Money*100),
		"trade_type": 1,
		"pay_type": model.PayType,
		"title": "充值",
		"client_ip": "127.0.0.1",
		"notify_url": config.GetValue("notify_url").ToString(),
		"remark": "充值",
	}
	log.Info(params)

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

	response,err := network.Post(config.GetValue("payapi_url").ToString()+"/pay/makeprepay",paramData,headers)
	if err!=nil{
		return nil,err
	}
	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		if err!=nil{
			return nil,err
		}

		return resultMap,nil
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return nil,errors.New(resultMap["err_msg"].(string))
	}

	return nil,errors.New("充值失败")
}
