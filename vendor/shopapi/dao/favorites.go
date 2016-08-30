package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

type Favorites struct  {
	OpenId string
	AppId string
	Title string
	CoverImg string
	Remark string
	Type int
	ObjId int64
	Flag string
	Json string
	BaseDModel
}

func NewFavorites() *Favorites  {

	return &Favorites{}
}

func (self *Favorites) Insert() error  {
	_,err :=db.NewSession().InsertInto("favorites").Columns("open_id","app_id","title","cover_img","remark","type","obj_id","flag","json").Record(self).Exec()
	return err
}

func (self *Favorites) WithOpenId(openId,appId string) ([]*Favorites,error)   {
	var list []*Favorites
	_,err :=db.NewSession().Select("*").From("favorites").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&list)
	return list,err
}

func (self *Favorites) WithTypeAndObjId(objId int64,typ int,openId string,appId string) (bool,error)   {
	var count int64
	_,err :=db.NewSession().Select("count(*)").From("favorites").Where("open_id=?").Where("obj_id=?",objId).Where("type=?",typ).Where("app_id=?",appId).LoadStructs(&count)
	if err!=nil{
		return false,err
	}
	if count<=0 {
		return false,nil
	}
	return true,nil
}

func (self *Favorites) DeleteWithId(id int64) error  {

	_,err :=db.NewSession().DeleteFrom("favorites").Where("id=?",id).Exec()

	return err
}