package service

import "shopapi/dao"

func DistributionProducts(added,openId,appId string)  ([]*dao.DistributionProductDetail,error)  {

	distDetail := dao.NewDistributionProductDetail()

	return distDetail.DetailWithAppId(added,openId,appId)
}
