package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Category struct {
	Id int64
	AppId string
	Title string
	Description string
	Icon string
	Flag string
	Json string
}

func NewCategory() *Category  {

	return &Category{}
}

func (self *Category) InsertTx(tx *dbr.Tx) error  {

	_,err := tx.InsertInto("category").Columns("app_id","title","description","icon","flag","json").Record(self).Exec()

	return err
}

func (self *Category) WithFlags(flags []string,noflags []string,appId string) ([]*Category,error)  {

	bulider :=db.NewSession().Select("*").From("category").Where("app_id=?",appId)

	if flags!=nil&&len(flags)>0 {
		bulider = bulider.Where("flag in ?",flags)
	}

	if noflags!=nil&&len(noflags) >0 {
		bulider = bulider.Where("flag not in ?",noflags)
	}
	var list []*Category
	_,err :=bulider.LoadStructs(&list)

	return list,err
}