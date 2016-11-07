package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	//"fmt"
	//"math/rand"  
	//"gitlab.qiyunxin.com/tangtao/utils/log" //log.Debug(x)
)

type AppLog struct  {
	Id			 int64  `json:"id"`
	ApkVersion 	 string `json:"apk_version"`
	ApkUrl 		 string `json:"apk_url"`
	IosVersion 	 string `json:"ios_version"`
	IosLog 		 string `json:"ios_log"`
}

func NewAppLog() *AppLog  {
	return &AppLog{}
}

//应用更新日志
func (self *AppLog) AppUpdateLog() (*AppLog,error) {
	var appLog *AppLog
	_,err :=db.NewSession().Select("*").From("app_update").Where("id=?",1).LoadStructs(&appLog)
	return appLog,err
}