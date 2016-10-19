package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type ProdImgs struct  {

	//产品ID
	ProdId int64
	AppId string
	Json string
	Url string
	Flag string
	BaseDModel

}

type ProdImgsDetail struct  {

	//产品ID
	ProdId int64
	AppId string
	Url string
	Flag string
	Json string
	BaseDModel

}

func NewProdImgsDetail() *ProdImgsDetail {

	return &ProdImgsDetail{}
}

func NewProdImgs() *ProdImgs {

	return &ProdImgs{}
}

func (self *ProdImgs) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("prod_imgs").Columns("app_id","prod_id","url","flag","json").Record(self).Exec()

	return err
}

func (self *ProdImgs) DeleteWithIdTx(prodId int64,tx *dbr.Tx) error {

	_,err :=tx.DeleteFrom("prod_imgs").Where("prod_id=?",prodId).Exec()

	return err
}

func (self *ProdImgsDetail) ProdImgsWithProdId(prodId int64,appId string) ([]*ProdImgsDetail,error)  {

	sess := db.NewSession()
	var details []*ProdImgsDetail
	_,err :=sess.SelectBySql("select * from prod_imgs ps where  ps.prod_id=? and ps.app_id=?",prodId,appId).LoadStructs(&details)

	return  details,err
}

func (self *ProdImgsDetail) ProdImgsWithProdIds(prodIds []int64) ([]*ProdImgsDetail,error){
	sess := db.NewSession()
	var details []*ProdImgsDetail
	_,err :=sess.SelectBySql("select * from prod_imgs ps where ps.prod_id in ?",prodIds).LoadStructs(&details)
	log.Debug("----err",err)
	return  details,err
}
