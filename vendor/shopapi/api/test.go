package api

import (
	"github.com/gin-gonic/gin"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"time"
	"fmt"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"errors"
	"io/ioutil"
)

const (
	APPID		="TdCk0qqZoa7LE0bw42pyu1"
	APPSECRET	="EOy3GR1AUp58xCwyQeQbz5"
	APPKEY		="QM43LR7hUz8AFx05GnKjS9"
	MASTERSECRET="MDtPCBiDEiAqGSBi8Ukba3"
)

func Test(c *gin.Context) {
	var err error
	getui := &Getui{}
	err=getui.Conn()	
	if err!=nil {
		c.JSON(http.StatusOK,err.Error())
		return
	}	
	c.JSON(http.StatusOK,getui.AuthToken)
	//getui.AuthToken="05a941277358ba42c2e27459034aa09b4fe33747fd33defdf2052f37d5741068"
	//getui.Test()
	//=======================================================
	status,err:=getui.Status("ed105d433efafb78c110729c821748c1")
	if err!=nil {
		c.JSON(http.StatusOK,err.Error())
		return
	}
	c.JSON(http.StatusOK,status)
	
	err=getui.Close()
	if err!=nil {
		c.JSON(http.StatusOK,err.Error())
		return
	}
}

type Getui struct {
	Timestamp string
	AuthToken string
}
func (self *Getui)Conn() error {
	self.Timestamp= fmt.Sprintf("%d",time.Now().UnixNano()/1e6)	
	hash 		 := sha256.New()
	io.WriteString(hash,fmt.Sprintf("%s%s%s",APPKEY,self.Timestamp,MASTERSECRET));
	sign:=fmt.Sprintf("%x",hash.Sum(nil));	
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_sign",APPID)
	payparams :=map[string]interface{}{
		"appkey"	: APPKEY,
		"sign"		: sign,
		"timestamp"	: self.Timestamp,
	}
	paramData,_:= json.Marshal(payparams);
	response,_ := network.Post(url,paramData,map[string]string{
		"Content-Type"	: "application/json",
	})
	/* if response.StatusCode==http.StatusOK {		
	}else if response.StatusCode==http.StatusBadRequest {		
	}else{		
	} */
	if response.StatusCode!=http.StatusOK {
		return errors.New("鉴权错误")
	}
	resultMap 	:= map[string]interface{}{}
	util.ReadJsonByByte([]byte(response.Body),&resultMap)
	if "ok"!=resultMap["result"].(string) {
		return errors.New("鉴权错误,"+resultMap["result"].(string))
	}
	
	self.AuthToken= resultMap["auth_token"].(string)
	return nil
}
func (self *Getui)Status(cid string) (string,error) {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/user_status/%s",APPID,cid)
	response,_ := network.Get(url,map[string]string{},map[string]string{
		"Content-Type"	: "application/json",
		"authtoken"		: self.AuthToken,
	})
	if response.StatusCode!=http.StatusOK {
		return "",errors.New("状态检测错误")
	}
	resultMap 	:= map[string]interface{}{}
	util.ReadJsonByByte([]byte(response.Body),&resultMap)
	if "ok"!=resultMap["result"].(string) {
		return "",errors.New("状态检测错误,"+resultMap["result"].(string))
	}
	return resultMap["status"].(string),nil
}
func (self *Getui)Close() error {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_close",APPID)
	payparams	:=map[string]interface{}{}
	paramData,_ := json.Marshal(payparams);
	response,_  := network.Post(url,paramData,map[string]string{
		"authtoken": self.AuthToken,
	})
	if response.StatusCode!=http.StatusOK {
		return errors.New("关闭错误")
	}
	resultMap 	:= map[string]interface{}{}
	util.ReadJsonByByte([]byte(response.Body),&resultMap)
	if "ok"!=resultMap["result"].(string) {
		return errors.New("关闭错误,"+resultMap["result"].(string))
	}
	return nil
}
func (self *Getui)Test() {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/user_status/%s",APPID,"ed105d433efafb78c110729c821748c1")
	url="http://127.0.0.1/test.php"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authtoken", self.AuthToken)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}


































