package service

import (
	"shopapi/dao"
)

//一点公益ID号绑定
func YdgySetId(openId string,ydgyId int64) error {
	return dao.YdgySetId(openId,ydgyId)
}
//一点公益ID号获取
func YdgyGetId(openId string) (int64,error) {
	return dao.YdgyGetId(openId)
}