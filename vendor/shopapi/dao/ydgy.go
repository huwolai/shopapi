package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

//一点公益ID号绑定
func YdgySetId(openId string,ydgyId int64) error {
	_,err :=db.NewSession().Update("account").Set("ydgy_id",ydgyId).Where("open_id=?",openId).Limit(1).Exec()
	return err
}
//一点公益ID号获取
func YdgyGetId(openId string) (int64,error) {
	bulider :=db.NewSession().Select("ydgy_id").From("account").Where("open_id=?",openId)

	var ydgyId int64
	
	_,err :=bulider.Limit(1).LoadStructs(&ydgyId)
	
	return ydgyId,err
}