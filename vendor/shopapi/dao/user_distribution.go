package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type UserDistribution struct  {
	Id int64
	AppId string
	OpenId string
	//分销编号
	Code string
	//分销ID
	DistributionId int64
	//商品ID
	ProdId int64
	BaseDModel
}

type UserDistributionDetail struct  {
	Id int64
	AppId string
	OpenId string
	//分销编号
	Code string
	//分销ID
	DistributionId int64
	//商品ID
	ProdId int64
	MerchantId int64
	//佣金比例
	CsnRate float64

}

func NewUserDistribution() *UserDistribution  {

	return &UserDistribution{}
}

func NewUserDistributionDetail() *UserDistributionDetail  {

	return &UserDistributionDetail{}
}

func (self *UserDistribution) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("user_distribution").Columns("app_id","open_id","code","prod_id","distribution_id").Record(self).Exec()

	return err
}
func (self *UserDistribution) Insert() error  {

	_,err :=db.NewSession().InsertInto("user_distribution").Columns("app_id","open_id","code","prod_id","distribution_id").Record(self).Exec()

	return err
}

func (self *UserDistribution) DeleteWithDistributionId(distributionId int64,openId,appId string) error  {
	_,err :=db.NewSession().DeleteFrom("user_distribution").Where("distribution_id=?",distributionId).Where("open_id=?",openId).Where("app_id=?",appId).Exec()

	return err
}

func (self *UserDistribution) WithCode(code string) (*UserDistribution,error)   {
	var model *UserDistribution
	_,err :=db.NewSession().Select("*").From("user_distribution").Where("code=?",code).LoadStructs(&model)
	return model,err
}

func (self *UserDistributionDetail) WithCode(code string)  (*UserDistributionDetail,error) {
	var model *UserDistributionDetail
	_,err :=db.NewSession().SelectBySql("select ud.*,dt.merchant_id,dt.csn_rate from user_distribution ud,distribution_product dt where ud.distribution_id=dt.id and ud.code=?",code).LoadStructs(&model)
	return model,err
}