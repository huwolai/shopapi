package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Product struct  {
	Id int64
	AppId string
	//商品标题
	Title string
	//商品描述
	Description string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//商品状态
	Status int
	//附加数据
	Json string
}

type ProductDetail struct {
	//商品ID
	Id int64
	AppId string
	//商品标题
	Title string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//商品状态
	Status int
	//商户ID
	MerchantId int64
	//商户名称
	MerchantName string
	//附加数据
	Json string
	//商品图片集合
	ProdImgs []*ProdImgsDetail
}

func NewProductDetail() *ProductDetail {

	return &ProductDetail{}
}

func NewProduct() *Product  {

	return &Product{}
}

func (self *Product) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("product").Columns("title","app_id","description","price","dis_price","json","status").Record(self).Exec()
	if err !=nil {

		return 0,err
	}
	pid,err :=  result.LastInsertId()
	return pid,err
}

//商品推荐列表
func (self *ProductDetail) ProductListWithRecomm(appId string) ([]*ProductDetail,error)  {
	session := db.NewSession()
	var prodList []*ProductDetail
	_,err :=session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.price,pt.dis_price,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,prod_category pc,product pt where md.app_id=pc.app_id and md.prod_id=pc.prod_id and md.merchant_id=mt.id and pc.app_id=pt.app_id and pc.prod_id=pt.id and pt.is_recom=1 and pt.app_id=?",appId).LoadStructs(&prodList)
	if err!=nil{
		log.Debug("----err",err)
		return nil,err
	}

	err = fillProdImgs(appId,prodList)

	return prodList,err
}

func (self *ProductDetail) ProductListWithCategory(appId string,categoryId int64) ([]*ProductDetail,error)  {
	session := db.NewSession()
	var prodList []*ProductDetail
	_,err :=session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.price,pt.dis_price,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,prod_category pc,product pt where md.app_id=pc.app_id and md.prod_id=pc.prod_id and md.merchant_id=mt.id and pc.app_id=pt.app_id and pc.prod_id=pt.id and pc.category_id=? and pt.app_id=?",categoryId,appId).LoadStructs(&prodList)
	if err!=nil{
		return nil,err
	}

	err = fillProdImgs(appId,prodList)

	return prodList,err
}

//填充商品图片数据
func fillProdImgs(appId string,prodList []*ProductDetail) error {
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
	prodImgDetails,err := prodImgDetail.ProdImgsWithProdIds(prodids,appId)
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
			for _,test :=range pdimgdetails {
				log.Debug(test.ImgNo)
			}
		}
	}
	for _,prod :=range prodList {
		key := prod.Id
		prodimgs := prodimgsMap[key]
		prod.ProdImgs = prodimgs
	}

	return nil
}


func (self *Product) ProductWithId(appId string,id int64) (*Product,error)  {
	sess :=db.NewSession()
	var prod *Product
	_,err :=sess.Select("*").From("product").Where("app_id=?",appId).Where("id=?",id).LoadStructs(&prod)

	return prod,err
}