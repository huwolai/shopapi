package service

import (
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"errors"
)

type ImprestsModel struct  {
	//用户openID
	OpenId string
	//预付款code
	Code string
	//金额 单位(分)
	Amount int64
	//预付款标题
	Title string
	//备注
	Remark string
}


//领取预付款
func FetchImprests(model *ImprestsModel) (string,error)  {

	payparams :=map[string]interface{}{
		"open_id": model.OpenId,
		"code": model.Code,
		"amount": model.Amount,
		"title": model.Title,
		"remark": model.Remark,
	}

	//获取接口签名信息
	noncestr,timestamp,appid,basesign,sign  :=GetPayapiSign(payparams)
	//header参数
	headers := map[string]string{
		"app_id": appid,
		"sign": fmt.Sprintf("%s.%s",basesign,sign),
		"noncestr": noncestr,
		"timestamp": timestamp,
	}
	paramData,_:= json.Marshal(payparams);
	response,err := network.Post(config.GetValue("payapi_url").ToString()+"/imprests/fetch",paramData,headers)
	if err!=nil{
		log.Error(err)
		return "",err
	}

	if response.StatusCode==http.StatusOK {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)
		return resultMap["sub_trade_no"],err
	}else if response.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(response.Body),&resultMap)

		return "",errors.New(resultMap["err_msg"].(string))
	}else{
		return "",errors.New("请求支付中心失败!")
	}

}
