package service


type ProdBLL struct {
	AppId string
	//商品标题
	Title string
	//商户ID
	MerchantId int64
	//描述
	Description string
	//类别ID
	CategoryId int64
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//图片编号集合
	ImgNos string
	Json string
}

type ProductResultDLL struct  {
	Id int64
	//商品标题
	Title string
	//描述
	Description string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	Json string

}

type ProdImgsDetailDLL struct  {
	//图片编号
	ImgNo string
	//产品ID
	ProdId int64
	AppId string
	Url string
	Flag string
	Json string
}
