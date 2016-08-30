package service

import "shopapi/dao"

func DistributionProducts(appId string)  ([]*dao.DistributionProductDetail,error)  {

	distDetail := dao.NewDistributionProductDetail()

	return distDetail.DetailWithAppId(appId)
}
