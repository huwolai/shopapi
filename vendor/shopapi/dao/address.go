package dao

import "github.com/gocraft/dbr"

type Address struct  {
	Id int64
	AppId string
	OpenId string
	Longitude float64
	Latitude float64
	Address string
	//'0不是默认,1是默认地址'
	IsDefault int
	Json string
}

func (self *Address) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("address").Columns("app_id","open_id","longitude","latitude","address","is_default","json").Record(self).Exec()
	if err!=nil{

		return 0,err
	}
	lastId,err := result.LastInsertId()

	return lastId,err

}
