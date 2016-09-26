package service


type ProdBLL struct {
	Id int64
	AppId string
	//商品标题
	Title string
	//子标题
	SubTitle string
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
	//图片集合
	Imgs []ProdImgBLL
	//标记
	Flag string
	Json string
}

type ProdImgBLL struct {
	Flag string
	Url string
	Json string
	ProdId int64
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
	//产品ID
	ProdId int64
	AppId string
	Url string
	Flag string
	Json string
}
