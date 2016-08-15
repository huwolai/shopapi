package dao

import "github.com/gocraft/dbr"

type Merchant struct  {
	Name string `json:"name"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Status int `json:"status"`
	Json string `json:"json"`
	BaseDModel
}

func NewMerchant() *Merchant  {

	return &Merchant{}
}

func (self *Merchant) InsertTx(tx *dbr.Tx) (int64,error) {

	result,err :=tx.InsertInto("merchant").Columns("name","app_id","open_id","status","json").Record(self).Exec()
	if err!=nil{
		return 0,err
	}
	lastId,err := result.LastInsertId()
	return lastId,err
}