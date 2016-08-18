package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type Address struct  {
	Id int64
	AppId string
	OpenId string
	Longitude float64
	Latitude float64
	Address string
	Weight int
	Json string
}

func NewAddress() *Address  {

	return &Address{}
}

func (self *Address) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("address").Columns("app_id","open_id","longitude","latitude","address","is_default","json").Record(self).Exec()
	if err!=nil{

		return 0,err
	}
	lastId,err := result.LastInsertId()

	return lastId,err

}

func (self *Address) Insert() (int64,error)  {

	result,err :=db.NewSession().InsertInto("address").Columns("app_id","open_id","longitude","latitude","address","is_default","json").Record(self).Exec()
	if err!=nil{

		return 0,err
	}
	lastId,err := result.LastInsertId()

	return lastId,err

}

//查询推荐用户地址
func (self *Address) AddressWithRecom(openId string,appId string) (*Address,error) {
	var address *Address
	_,err :=db.NewSession().Select("*").From("address").Where("open_id=?",openId).Where("app_id=?",appId).OrderDir("weight",false).LoadStructs(&address)
	return address,err
}

func (self *Address) AddressWithOpenId(openId string,appId string) ([]*Address,error)  {
	var address []*Address
	_,err :=db.NewSession().Select("*").From("address").Where("open_id=?",openId).Where("app_id=?",appId).OrderDir("weight",false).LoadStructs(&address)
	return address,err
}
