package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Flags struct  {
	AppId string
	Name string
	Flag string
	Type string
	BaseDModel
}

func NewFlags() *Flags  {

	return &Flags{}
}

func (self *Flags) Insert() error {

	_,err :=db.NewSession().InsertInto("flags").Columns("app_id","name","flag","type").Record(self).Exec()

	return err
}

//查询标记通过类型
func (self *Flags) WithTypes(stype []string,appId string) ([]*Flags,error)  {
	log.Error("type:",stype,"app_id:",appId)
	var flags []*Flags
	builder :=db.NewSession().Select("*").From("flags").Where("app_id=?",appId)
	if stype!=nil{
		builder = builder.Where("type in ?",stype)
	}
	_,err :=builder.LoadStructs(&flags)
	return flags,err
}