package service

import (	
	"shopapi/dao"
	"errors"
	"fmt"
)
func PushSingle(openId string,appId string,title string,content string,types string) error {
	account,err :=dao.NewAccount().AccountWithOpenId(openId,appId)
	if err!=nil {
		return err
	}
	if account==nil || account.Getui =="" {
		return errors.New("Getui为空")
	}
	
	getuiJson,_:=dao.JsonToMap(account.Getui)	
	
	
	if getuiJson["cid"].(string)=="" {
		return errors.New("cid为空")
	}
	
	var msg dao.AccountMsg
	msg.Cid			=getuiJson["cid"].(string)
	msg.Devicetoken	=getuiJson["devicetoken"].(string)
	msg.Title		=title
	msg.Content		=content
	msg.Types		=types
	msg.AppId		=appId
	msg.OpenId		=openId

	getui:= dao.NewGetui()
	getui.Conn()
	status,err:=getui.Status(msg.Cid)
	
	if err!=nil {
		return err
	}
	
	fmt.Println(status)
	
	if status=="online" {		
		err:=getui.PushSingle(msg)
		return err
	}
	//ios	
	if msg.Devicetoken!="" {
		err:=getui.PushApnsSingle(msg)		
		return err
	}
	//====
	err=getui.PushSingleOffline(msg)
	return err
}













