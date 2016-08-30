package dao

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

