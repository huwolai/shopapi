package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
	"github.com/gocraft/dbr"
)

func ProdDetailWithProdId(prodId int64,appId string) (*dao.Product,error)  {

	product := dao.NewProduct()
	return product.ProductWithId(prodId,appId)
}

func  ProdImgsWithProdId(prodId int64,appId string) ([]*ProdImgsDetailDLL,error) {
	prodImgDetail := dao.NewProdImgsDetail()
	prodImgDetals,err := prodImgDetail.ProdImgsWithProdId(prodId,appId)
	if err!=nil{

		return nil,err
	}
	detailDLLs := make([]*ProdImgsDetailDLL,0)
	if prodImgDetals!=nil {
		for _,detail :=range prodImgDetals{
			detailDLLs = append(detailDLLs,prodImgsDetailToDLL(detail))
		}
	}

	return detailDLLs,nil
}

//添加商品
func ProdAdd(prodbll *ProdBLL) error  {
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	//保存商品基础信息
	prodId,err := productBaseSave(prodbll,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	//保存商品图片信息
	if len(prodbll.Imgs)>0 {
		err := productImgSave(prodbll,prodId,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
	}

	//保存商品所属信息
	err = merchantProdAdd(prodbll,prodId,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//商品分类添加
	if err:=prodCategoryAdd(prodbll,prodId,tx);err!=nil{
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return err
	}

	return nil;
}

func ProductListWithRecomm(appId string) ([]*dao.ProductDetail,error) {
	productDetail :=dao.NewProductDetail()
	prodList,err := productDetail.ProductListWithRecomm(appId)
	if err!=nil {
		return nil,err
	}

	return prodList,nil
}

func ProductListWithCategory(appId string,categoryId int64) ([]*dao.ProductDetail,error)   {
	productDetail :=dao.NewProductDetail()
	prodList,err := productDetail.ProductListWithCategory(appId,categoryId)
	if err!=nil {
		return nil,err
	}

	return prodList,nil
}

//产品分类添加
func prodCategoryAdd(prodbll *ProdBLL,prodId int64,tx *dbr.Tx) error {
	prodCategory := dao.NewProdCategory()
	prodCategory.CategoryId = prodbll.CategoryId
	prodCategory.ProdId = prodId
	prodCategory.AppId = prodbll.AppId

	if err :=prodCategory.InsertTx(tx);err!=nil{

		return err
	}
	return nil;
}

func prodImgsDetailToDLL(detail *dao.ProdImgsDetail) *ProdImgsDetailDLL  {

	dll := &ProdImgsDetailDLL{}
	dll.AppId = detail.AppId
	dll.Flag = detail.Flag
	dll.Json = detail.Json
	dll.ProdId = detail.ProdId
	dll.Url = detail.Url

	return dll
}

func productToReusltDLL(prod *dao.Product)  *ProductResultDLL  {
	dll := &ProductResultDLL{}
	dll.Description = prod.Description
	dll.DisPrice = prod.DisPrice
	dll.Json = prod.Json
	dll.Price = prod.Price
	dll.Title = prod.Title
	dll.Id = prod.Id
	return dll
}

//商户产品添加
func merchantProdAdd(prodbll *ProdBLL,prodId int64,tx *dbr.Tx) error  {

	mprod := dao.NewMerchantProd()
	mprod.MerchantId = prodbll.MerchantId
	mprod.ProdId = prodId
	mprod.AppId = prodbll.AppId

	return mprod.InsertTx(tx)
}
//保存商品图片
func productImgSave(prodbll *ProdBLL,prodId int64,tx *dbr.Tx) error  {

	imgArray := prodbll.Imgs
	for _,img :=range imgArray  {
		prodImgs :=dao.NewProdImgs()
		prodImgs.AppId = prodbll.AppId
		prodImgs.ProdId = prodId
		prodImgs.Url = img.Url
		prodImgs.Flag = img.Flag
		prodImgs.Json = img.Json
		err :=prodImgs.InsertTx(tx)
		if err!=nil{
			return err
		}
	}

	return nil;
}

//保存商品基础信息
func productBaseSave(prodbll *ProdBLL,tx *dbr.Tx) (int64,error)  {
	prod := dao.NewProduct()
	prod.Title = prodbll.Title
	prod.AppId = prodbll.AppId
	prod.Description = prodbll.Description
	prod.DisPrice = prodbll.DisPrice
	prod.Price = prodbll.Price
	prod.Status = comm.PRODUCT_STATUS_NORMAL
	prod.Json = prodbll.Json
	prodId,err := prod.InsertTx(tx)
	if err!=nil{
		tx.Rollback()
		return 0,err
	}

	return prodId,err
}
