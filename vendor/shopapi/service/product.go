package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
	"github.com/gocraft/dbr"
	"errors"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type ProdAndAttrDto struct  {
	AttrKey string `json:"attr_key"`
	AttrValue string `json:"attr_value"`
	AppId string `json:"app_id"`
	ProdId int64 `json:"prod_id"`
	SkuNo string `json:"sku_no"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

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
func ProdAdd(prodbll *ProdBLL) (*ProdBLL,error)  {
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
		return nil,err
	}

	prodbll.Id = prodId

	//保存商品图片信息
	if len(prodbll.Imgs)>0 {
		err := productImgSave(prodbll,prodId,tx)
		if err!=nil{
			tx.Rollback()
			return nil,err
		}
	}

	//保存商品所属信息
	err = merchantProdAdd(prodbll,prodId,tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}
	//商品分类添加
	if err:=prodCategoryAdd(prodbll,prodId,tx);err!=nil{
		tx.Rollback()
		return nil,err
	}
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	return prodbll,nil;
}

type ProdAttrKeyDto struct  {
	Id int64 `json:"id"`
	//商品ID
	ProdId int64 `json:"prod_id"`
	//属性key
	AttrKey string `json:"attr_key"`
	//属性名
	AttrName string `json:"attr_name"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

type ProdAttrValueDto struct  {
	Id int64 `json:"id"`
	ProdId int64 `json:"prod_id"`
	AttrKey string `json:"attr_key"`
	AttrValue string `json:"attr_value"`
	Flag string `json:"flag"`
	Json string  `json:"json"`
}

//添加商品属性
func ProdAttrKeyAdd(dto *ProdAttrKeyDto) (*ProdAttrKeyDto,error) {
	prodAttrKey := ProdAttrKeyDtoToModel(dto)
	lastId,err :=prodAttrKey.Insert()

	dto.Id = lastId

	return dto,err
}

//添加商品属性值
func ProdAttrValueAdd(dto *ProdAttrValueDto) (*ProdAttrValueDto,error)  {
	prodAttrValue :=ProdAttrValueDtoToModel(dto)
	lastId,err :=prodAttrValue.Insert()
	dto.Id = lastId
	return dto,err
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

//查询商品指定key的属性值
func ProductAttrValues(vsearch string,attrKey string,prodId int64) ([]*dao.ProdAttrVal,error)  {
	prodAttrVal := dao.NewProdAttrVal()
	prodAttrVals,err :=prodAttrVal.WithAttrKeyStock(vsearch,attrKey,prodId)

	return prodAttrVals,err
}

func ProductAndAttrAdd(dto *ProdAndAttrDto)  (*ProdAndAttrDto,error) {

	product :=dao.NewProduct()
	product,err :=product.ProductWithId(dto.ProdId,dto.AppId)
	if err!=nil{

		return nil,err
	}
	if product==nil  {
		return nil,errors.New("没有找到此商户发布的产品!")
	}

	prodAttrVal := dao.NewProdAttrVal()
	prodAttrVal,err =prodAttrVal.WithProdValue(dto.AttrValue,product.Id)
	if err!=nil{
		return nil,err
	}

	if prodAttrVal!=nil{
		prodSku :=dao.NewProdSku()
		prodSku,err =prodSku.WithProdIdAndSymbolPath(strconv.Itoa(int(prodAttrVal.Id)),product.Id)
		if err!=nil{
			return nil,err
		}
		if prodSku!=nil{
			if prodSku.Stock==0{
				return nil,errors.New("此商品已无库存!")
			}else{
				dto.SkuNo = prodSku.SkuNo
				return dto,nil
			}
		}
	}

	tx,_ := db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
		}
	}()
	//生成商品属性
	prodAttrVal = dao.NewProdAttrVal()
	prodAttrVal.AttrValue = dto.AttrValue
	prodAttrVal.AttrKey = dto.AttrKey
	prodAttrVal.Flag = dto.Flag
	prodAttrVal.ProdId = product.Id
	prodAttrVal.Json = dto.Json

	lastId,err :=prodAttrVal.InsertTx(tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}
	prodAttrVal.Id=lastId

	//添加商品sku
	prodSku :=dao.NewProdSku()
	prodSku.ProdId = product.Id
	prodSku.AppId = dto.AppId
	prodSku.AttrSymbolPath = strconv.Itoa(int(product.Id))
	prodSku.DisPrice = product.DisPrice
	prodSku.Price = product.Price
	prodSku.Stock=1
	prodSku.SkuNo = util.GenerUUId()
	err =prodSku.InsertTx(tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	dto.SkuNo = prodSku.SkuNo

	return dto,nil
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

func ProdAttrKeyToDto(model *dao.ProdAttrKey) *ProdAttrKeyDto  {

	dto :=&ProdAttrKeyDto{}
	dto.AttrKey = model.AttrKey
	dto.AttrName = model.AttrName
	dto.Id = model.Id
	dto.ProdId = model.ProdId
	dto.Flag = model.Flag
	dto.Json = model.Json

	return dto
}

func ProdAttrKeyDtoToModel(dto *ProdAttrKeyDto) *dao.ProdAttrKey  {
	model := &dao.ProdAttrKey{}
	model.ProdId = dto.ProdId
	model.AttrName = dto.AttrName
	model.AttrKey = dto.AttrKey
	model.Flag = dto.Flag
	model.Json = dto.Json
	model.Id = dto.Id

	return model
}

func ProdAttrValueDtoToModel(dto *ProdAttrValueDto) *dao.ProdAttrVal  {
	model := dao.NewProdAttrVal()
	model.Flag = dto.Flag
	model.ProdId = dto.ProdId
	model.AttrKey = dto.AttrKey
	model.AttrValue = dto.AttrValue
	model.Json = dto.Json
	model.Id = dto.Id


	return model
}
