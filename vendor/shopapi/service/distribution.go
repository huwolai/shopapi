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

func DistributionProductCancel(distributionId int64,openId,appId string) error  {

	return dao.NewUserDistribution().DeleteWithDistributionId(distributionId,openId,appId)
}