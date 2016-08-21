package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type MerchantImgs struct  {
	Id int64
	//商户ID
	MerchantId int64
	OpenId string
	AppId string
	Json string
	Url string
	Flag string
	BaseDModel

}

func NewMerchantImgs() *MerchantImgs  {

	return &MerchantImgs{}
}

func (self *MerchantImgs) InsertTx(tx *dbr.Tx) error {

	_,err :=tx.InsertInto("merchant_imgs").Columns("merchant_id","app_id","json","url","flag","open_id").Record(self).Exec()

	return err
}

func (self *MerchantImgs) MerchantImgsUpdateTx( merchantImg *MerchantImgs,tx *dbr.Tx)  error {

	_,err :=tx.Update("merchant_imgs").Set("url",merchantImg.Url).Set("flag",merchantImg.Flag).Set("json",merchantImg.Json).Exec()

	return err
}

func (self *MerchantImgs) MerchantImgsWithId(id int64) (*MerchantImgs,error)  {
	var merchantImgs *MerchantImgs
	_,err :=db.NewSession().Select("*").From("merchant_imgs").Where("id=?",id).LoadStructs(&merchantImgs)

	return merchantImgs,err
}

func (self *MerchantImgs) MerchantImgsWithMerchantId(merchantId int64,flags []string) ([]*MerchantImgs,error)  {
	var merchantImgs []*MerchantImgs
	builder :=db.NewSession().Select("*").From("merchant_imgs").Where("merchant_id=?",merchantId)
	if flags!=nil&&len(flags)>0 {
		builder = builder.Where("flag in ?",flags)
	}

	_,err :=builder.LoadStructs(&merchantImgs)

	return merchantImgs,err
}

func (self *MerchantImgs) MerchantImgsWithFlag(flags []string,openId string,appId string) ([]*MerchantImgs,error)  {
	var merchantImgs []*MerchantImgs
	_,err := db.NewSession().Select("*").From("merchant_imgs").Where("flag in ?",flags).Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&merchantImgs)

	return merchantImgs,err
}

func (self *MerchantImgs) MerchantImgs(openId string,appId string) ([]*MerchantImgs,error)  {
	var merchantImgs []*MerchantImgs
	_,err := db.NewSession().Select("*").From("merchant_imgs").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&merchantImgs)

	return merchantImgs,err
}

