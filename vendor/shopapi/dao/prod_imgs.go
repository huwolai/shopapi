package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type ProdImgs struct  {
	//图片编号
	ImgNo string
	//产品ID
	ProdId int64
	AppId string
	Json string
	BaseDModel

}

type ProdImgsDetail struct  {
	//图片编号
	ImgNo string
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

	_,err :=tx.InsertInto("prod_imgs").Columns("app_id","prod_id","img_no","json").Record(self).Exec()

	return err
}

func (self *ProdImgsDetail) ProdImgsWithProdId(prodId int64,appId string) ([]*ProdImgsDetail,error)  {

	sess := db.NewSession()
	var details []*ProdImgsDetail
	_,err :=sess.SelectBySql("select * from prod_imgs ps,images gs where ps.app_id=gs.app_id and ps.img_no=gs.no and ps.prod_id=? and ps.app_id=?",prodId,appId).LoadStructs(&details)

	return  details,err
}

func (self *ProdImgsDetail) ProdImgsWithProdIds(prodIds []int64,appId string) ([]*ProdImgsDetail,error){
	sess := db.NewSession()
	var details []*ProdImgsDetail
	_,err :=sess.SelectBySql("select * from prod_imgs ps,images gs where ps.app_id=gs.app_id and ps.img_no=gs.no and ps.prod_id in ? and ps.app_id=?",prodIds,appId).LoadStructs(&details)

	return  details,err
}