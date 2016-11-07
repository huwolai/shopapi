package service

import (
	"shopapi/dao"
)
//应用更新日志
func AppUpdateLog()(*dao.AppLog,error) {
	appLog :=dao.NewAppLog()
	appLog,err := appLog.AppUpdateLog()
	return appLog,err
}