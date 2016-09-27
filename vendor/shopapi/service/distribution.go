package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

func DistributionProducts(added,openId,appId string)  ([]*dao.DistributionProductDetail,error)  {

	distDetail := dao.NewDistributionProductDetail()

	return distDetail.DetailWithAppId(added,openId,appId)
}

func DistributionWithMerchant(merchantId int64,appId string) ([]*dao.DistributionProductDetail,error)  {
	distDetail := dao.NewDistributionProductDetail()
	return distDetail.DistributionWithMerchant(merchantId,appId)
}

func DistributionProductAdd(distributionId int64,openId string,appId string) (*dao.UserDistribution,error)  {

	distributionProduct := dao.NewDistributionProduct()
	distributionProduct,err :=distributionProduct.WithId(distributionId)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	if distributionProduct==nil{
		return nil,errors.New("分销商品不存在!")
	}

	userDistribution := dao.NewUserDistribution()
	userDistribution.Code = util.GenerUUId()
	userDistribution.DistributionId = distributionId
	userDistribution.OpenId = openId
	userDistribution.ProdId = distributionProduct.ProdId
	userDistribution.AppId = appId
	err =userDistribution.Insert()
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	return userDistribution,nil
}

func ProductUpdateDistribution(distributionId int64,csnRate float64,appId string) error  {

	return dao.NewDistributionProduct().UpdateWithId(distributionId,csnRate,appId)
}

func ProductJoinDistribution(prodId int64,csnRate float64,appId string) error  {

	merchantProd,err := dao.NewMerchantProd().WithProdId(prodId,appId)
	if err!=nil{
		return err
	}
	if merchantProd==nil{
		return errors.New("商品没有关联商户!")
	}

	distrib,err := dao.NewDistributionProduct().WithProdId(prodId)
	if err!=nil{
		return err
	}
	if distrib!=nil{

		return errors.New("商品已经加入分销!")
	}

	distributionProduct := dao.NewDistributionProduct()
	distributionProduct.ProdId = prodId
	distributionProduct.MerchantId = merchantProd.MerchantId
	distributionProduct.CsnRate = csnRate
	distributionProduct.AppId = appId

	return  distributionProduct.Insert()
}

func DistributionProductCancel(distributionId int64,openId,appId string) error  {

	return dao.NewUserDistribution().DeleteWithDistributionId(distributionId,openId,appId)
}

func DistributionWith(keyword string,pageIndex,pageSize uint64,noflags []string,flags []string) ([]*dao.DistributionProductDetail2,error)  {

	return dao.NewDistributionProductDetail2().With(keyword,pageIndex,pageSize,noflags,flags)
}

func DistributionWithCount(keyword string,noflags []string,flags []string) (int64,error)  {

	return dao.NewDistributionProductDetail2().WithCount(keyword,noflags,flags)
}