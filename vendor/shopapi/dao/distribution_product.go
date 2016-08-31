package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type DistributionProduct struct {
	Id int64
	AppId string
	ProdId int64
	MerchantId int64
	CsnRate float64
	BaseDModel
}

func NewDistributionProduct() *DistributionProduct  {

	return &DistributionProduct{}
}


type DistributionProductDetail struct {
	//商品ID
	Id int64
	AppId string
	//商品标题
	Title string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//分销编号
	DbnNo string
	//商品状态
	Status int
	//佣金比例
	CsnRate float64
	//是否已添加分销
	Added int
	DistributionId int64
	//商户ID
	MerchantId int64
	//商户名称
	MerchantName string
	Flag string
	//附加数据
	Json string
	//商品图片集合
	ProdImgs []*ProdImgsDetail
}

func NewDistributionProductDetail() *DistributionProductDetail  {

	return &DistributionProductDetail{}
}

func (self *DistributionProduct) WithId(id int64) (*DistributionProduct,error) {
	var model *DistributionProduct
	_,err :=db.NewSession().Select("*").From("distribution_product").Where("id=?",id).LoadStructs(&model)

	return model,err
}

func (self *DistributionProductDetail) DistributionWithMerchant(merchantId int64,appId string) ([]*DistributionProductDetail,error)  {
	var prodList []*DistributionProductDetail
	_,err :=db.NewSession().SelectBySql("select pt.id,pt.app_id,pt.title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json,dp.csn_rate,dp.id distribution_id,ud.`code` dbn_no from merchant mt,product pt,distribution_product dp, user_distribution ud  where   dp.prod_id = ud.prod_id and dp.prod_id = pt.id and ud.open_id = mt.open_id and mt.id = ? and pt.status=1 and dp.app_id=?",merchantId,appId).LoadStructs(&prodList)

	err = FillDistributionProdImgs(appId,prodList)

	return prodList,err
}

func (self *DistributionProductDetail) DetailWithAppId(added string,openId string,appId string) ([]*DistributionProductDetail,error)  {
	session := db.NewSession()
	var prodList []*DistributionProductDetail


	if added=="" {
		_,err :=session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json,dp.csn_rate,if(ud.open_id is null,0,1) added,dp.id distribution_id from merchant_prod md,merchant mt,product pt,distribution_product dp left JOIN user_distribution ud on  dp.prod_id = ud.prod_id and ud.open_id=? where pt.status=1 and pt.app_id=dp.app_id and pt.id = dp.prod_id and dp.merchant_id = md.id  and md.merchant_id=mt.id and pt.app_id=? ",openId,appId).LoadStructs(&prodList)
		if err!=nil{
			return nil,err
		}
	}else{
		_,err :=session.SelectBySql("select * from (select pt.id,pt.app_id,pt.title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json,dp.csn_rate,if(ud.open_id is null,0,1) added,dp.id distribution_id from merchant_prod md,merchant mt,product pt,distribution_product dp left JOIN user_distribution ud on  dp.prod_id = ud.prod_id and ud.open_id=? where pt.status=1 and pt.app_id=dp.app_id and pt.id = dp.prod_id and dp.merchant_id = md.id  and md.merchant_id=mt.id and pt.app_id=?) tt where tt.added = ?",openId,appId,added).LoadStructs(&prodList)
		if err!=nil{
			return nil,err
		}
	}

	err := FillDistributionProdImgs(appId,prodList)

	return prodList,err
}

//填充商品图片数据
func FillDistributionProdImgs(appId string,prodList []*DistributionProductDetail) error {
	prodids := make([]int64,0)
	if prodList!=nil{
		for _,prod :=range prodList {
			prodids = append(prodids,prod.Id)
		}
	}

	if len(prodids)<=0 {
		return nil
	}

	prodImgDetail := NewProdImgsDetail()
	prodImgDetails,err := prodImgDetail.ProdImgsWithProdIds(prodids)
	if err!=nil{
		return err
	}
	prodimgsMap := make(map[int64][]*ProdImgsDetail)
	if prodImgDetails!=nil{
		for _,prodimgd :=range prodImgDetails {
			key := prodimgd.ProdId
			pdimgdetails :=prodimgsMap[key]
			if pdimgdetails==nil{
				pdimgdetails = make([]*ProdImgsDetail,0)
			}


			pdimgdetails= append(pdimgdetails,prodimgd)

			prodimgsMap[key] = pdimgdetails
			log.Debug(prodimgsMap)
		}
	}
	for _,prod :=range prodList {
		key := prod.Id
		prodimgs := prodimgsMap[key]
		prod.ProdImgs = prodimgs
	}

	return nil
}