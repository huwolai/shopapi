package dao

import "github.com/gocraft/dbr"

type ProdCategory struct  {
	Id int64
	CategoryId int64
	ProdId int64
	Json string
	AppId string
	BaseDModel
}

func NewProdCategory() *ProdCategory  {

	return &ProdCategory{}
}

func (self *ProdCategory) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("prod_category").Columns("app_id","category_id","prod_id","json").Record(self).Exec()

	return err
}
