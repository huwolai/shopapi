package service

import (
	"shopapi/dao"
)

//一点公益ID号绑定
func YdgySetId(openId string,ydgyId string,ydgyName string,ydgyMine int64) error {
	return dao.YdgySetId(openId,ydgyId,ydgyName,ydgyMine)
}
//一点公益ID号获取
func YdgyGetId(openId string) (*dao.Ydgy,error) {
	return dao.YdgyGetId(openId)
}
//一点公益ID号状态审核
func YdgySetIdWithStatus(openId string,YdgyStatus int64) error {
	return dao.YdgySetIdWithStatus(openId,YdgyStatus)
}
//一点公益ID号删除
func YdgySetIdWithDelete(openId string) error {
	return dao.YdgySetIdWithDelete(openId)
}