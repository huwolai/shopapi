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
	Json string
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
func (self *Flags) WithTypes(stype []string,appId string,status []string) ([]*Flags,error)  {
	log.Error("type:",stype,"app_id:",appId)
	var flags []*Flags
	builder :=db.NewSession().Select("*").From("flags").Where("app_id=?",appId)
	if stype!=nil{
		builder = builder.Where("type in ?",stype)
	}
	if status!=nil {
		builder = builder.Where("status in ?",status)
	}
	_,err :=builder.LoadStructs(&flags)
	return flags,err
}

func (self *Flags) FlagsSetJsonWithTypes(types string,json string,appId string) error  {

	_,err :=db.NewSession().Update("flags").Set("json",json).Where("type=?",types).Where("app_id=?",appId).Exec()

	return err
}

















