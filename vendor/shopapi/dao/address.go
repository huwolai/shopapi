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
	Name string
	Mobile string
	Address *string
	Weight int
	Json string
}

func NewAddress() *Address  {

	return &Address{}
}

func (self *Address) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("address").Columns("name","mobile","app_id","open_id","longitude","latitude","address","weight","json").Record(self).Exec()
	if err!=nil{

		return 0,err
	}
	lastId,err := result.LastInsertId()

	return lastId,err

}

func (self *Address) Insert() (int64,error)  {

	result,err :=db.NewSession().InsertInto("address").Columns("name","mobile","app_id","open_id","longitude","latitude","address","weight","json").Record(self).Exec()
	if err!=nil{

		return 0,err
	}
	lastId,err := result.LastInsertId()

	return lastId,err

}

func (self *Address) Update() error  {
	_,err :=db.NewSession().Update("address").Set("longitude",self.Longitude).Set("latitude",self.Latitude).Set("address",self.Address).Set("json",self.Json).Where("id=?",self.Id).Exec()

	return err
}

func (self *Address) Delete() error  {
	_,err :=db.NewSession().DeleteFrom("address").Where("id=?",self.Id).Exec()

	return err
}

func (self *Address) WithId(id int64) (*Address,error)  {
	var address *Address
	_,err :=db.NewSession().Select("*").From("address").Where("id=?",id).LoadStructs(&address)
	return address,err
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
