package service

import "shopapi/dao"

func DistributionProducts(openId,appId string)  ([]*dao.DistributionProductDetail,error)  {

	distDetail := dao.NewDistributionProductDetail()

	return distDetail.DetailWithAppId(openId,appId)
}
