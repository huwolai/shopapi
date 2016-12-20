package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	//"strings"
	"bytes"
	"io/ioutil"
	"net/http"
	//"net/url"
	"errors"
	"fmt"
	"encoding/json"
	"time"
	"crypto/sha256"
	"io"
)

const (
	APPID		="TdCk0qqZoa7LE0bw42pyu1"
	APPSECRET	="EOy3GR1AUp58xCwyQeQbz5"
	APPKEY		="QM43LR7hUz8AFx05GnKjS9"
	MASTERSECRET="MDtPCBiDEiAqGSBi8Ukba3"
)

type Getui struct  {
	Authtoken 	[]string
	ContentType []string
	TokenExpire	int64
}

func NewGetui() *Getui  {
	getui:=&Getui{}
	
	s :=make([]string,0)
	s =append(s,"application/json")	
	getui.ContentType=s
	
	return getui
}

func (self *Getui)Conn() error {
	s :=make([]string,0)
	
	type Json struct  {
		 Json string
	}

	var json *Json
	_,err :=db.NewSession().Select("*").From("flags").Where("flag=?","token").Where("type=?","GETUI").LoadStructs(&json)
	
	if len(json.Json)>0 {
		token,_:=JsonToMap(json.Json)			
		if time.Now().Unix()<int64(token["expire"].(float64)) {
			s =append(s,token["auth_token"].(string))
			self.Authtoken=s
			return nil
		}
	}
	//==========================================
	timestamp	 := fmt.Sprintf("%d",time.Now().UnixNano()/1e6)	
	hash 		 := sha256.New()
	io.WriteString(hash,fmt.Sprintf("%s%s%s",APPKEY,timestamp,MASTERSECRET));
	sign		 :=fmt.Sprintf("%x",hash.Sum(nil));	
	url			 :=fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_sign",APPID)
	payparams 	 :=map[string]interface{}{
		"appkey"	: APPKEY,
		"sign"		: sign,
		"timestamp"	: timestamp,
	}
	response,err := self.Post(url,payparams)
	s =append(s,response["auth_token"].(string))		
	self.Authtoken=s
	
	db.NewSession().Update("flags").Set("json",fmt.Sprintf("{\"auth_token\":\"%s\",\"expire\":%d}",response["auth_token"].(string),time.Now().Unix()+43200)).Where("flag=?","token").Where("type=?","GETUI").Exec()
	
	return err
}
func (self *Getui)Status(cid string)(string,error) {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/user_status/%s",APPID,cid)
	response,err := self.Get(url,map[string]string{})
	return response["status"].(string),err
}
func (self *Getui)PushSingle(cid string,title string,body string) error {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/push_single",APPID)
	message:= map[string]interface{}{
		"appkey"		: APPKEY, 
		"is_offline" 	: false,
		"msgtype" 		: "notification",
	}
	notification:= map[string]interface{}{
		"text"					: body, 
		"title" 				: title,		
		"notify_style" 			: 1,		
	}
	postData  := map[string]interface{}{
		"message"		: message, 
		"notification"	: notification, 
		"cid"			: cid, 
		"requestid"		: fmt.Sprintf("%d",time.Now().UnixNano()), 
	}	
	_,err := self.Post(url,postData)
	
	return err
}
func (self *Getui)PushApnsSingle(devictoken string,title string,body string,types string) error {
	url:=fmt.Sprintf("https://restapi.getui.com/v1/%s/push_apns_single",APPID)
	
	alert:= map[string]interface{}{
		"title": title, 
		"body" : body,
	}
	aps := map[string]interface{}{
		"alert"		: alert, 
		"autoBadge"	: "+1", 
	}
	pushInfo := map[string]interface{}{
		"aps"	 	: aps, 
		"payload" 	: "{\"title\":\""+title+"\",\"body\":\""+body+"\",\"type\":\""+types+"\"}",
	}
	postData  := map[string]interface{}{
		"device_token"	: devictoken, 
		"push_info"		: pushInfo, 
	}	
	_,err := self.Post(url,postData)
	
	return err
}
func (self *Getui)Post(url string,postData map[string]interface{}) (map[string]interface{},error) {
	postBody,_:= json.Marshal(postData);
	
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    if err != nil {
        return nil,err
    }
	
	req.Header["Content-Type"]	=self.ContentType
	req.Header["authtoken"]		=self.Authtoken
	
    return self.Json(req)
}
func (self *Getui)Get(url string,data map[string]string) (map[string]interface{},error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil,err
    }
	
	req.Header["Content-Type"]	=self.ContentType
	req.Header["authtoken"]		=self.Authtoken
 
   return self.Json(req)    
}
func (self *Getui)Json(req *http.Request) (map[string]interface{},error) {
	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
 
    body, err := ioutil.ReadAll(resp.Body)
	/* fmt.Println("                            ")
	fmt.Println("                            ")
	fmt.Println("                            ")
	fmt.Println(string(body))
	fmt.Println("                            ")
	fmt.Println("                            ")
	fmt.Println("                            ")	 */
    if err != nil {
        return nil,errors.New("连接失败")
    }
	
	if resp.StatusCode!=http.StatusOK {
		return nil,errors.New("连接失败"+fmt.Sprintf("%d",resp.StatusCode))
	}
	
	resultMap 	:= map[string]interface{}{}
	util.ReadJsonByByte([]byte(body),&resultMap)
	
	if "ok"!=resultMap["result"].(string) {
		return nil,errors.New("参数失败,"+resultMap["result"].(string))
	}
	
	return resultMap,nil
}





