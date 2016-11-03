package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type Ydgy struct {
	YdgyId		string	`json:"id"`
	YdgyName	string	`json:"name"`
	YdgyMine	int64	`json:"mine"`
	YdgyStatus	int64	`json:"stauts"`
}

//一点公益ID号绑定
func YdgySetId(openId string,ydgyId string,ydgyName string,ydgyMine int64) error {
	_,err :=db.NewSession().Update("account").Set("ydgy_status",1).Set("ydgy_id",ydgyId).Set("ydgy_name",ydgyName).Set("ydgy_mine",ydgyMine).Where("open_id=?",openId).Limit(1).Exec()
	return err
}
//一点公益ID号获取
func YdgyGetId(openId string) (*Ydgy,error) {
	bulider :=db.NewSession().Select("ydgy_id,ydgy_name,ydgy_mine,ydgy_status").From("account").Where("open_id=?",openId)

	var ydgy *Ydgy
	
	_,err :=bulider.Limit(1).LoadStructs(&ydgy)
	
	return ydgy,err
}
//一点公益ID号状态审核
func YdgySetIdWithStatus(openId string,YdgyStatus int64,ydgyRes string) error {

	builder :=db.NewSession().Update("account").Set("ydgy_status",YdgyStatus)
	if len(ydgyRes)>0 {
		builder = builder.Set("ydgy_fail_res",ydgyRes)
	}
	_,err :=builder.Where("open_id=?",openId).Limit(1).Exec()
	return err
}
//一点公益ID号删除
func YdgySetIdWithDelete(openId string) error {
	_,err :=db.NewSession().Update("account").Set("ydgy_id","").Set("ydgy_name","").Set("ydgy_status",0).Where("open_id=?",openId).Limit(1).Exec()
	return err
}


























