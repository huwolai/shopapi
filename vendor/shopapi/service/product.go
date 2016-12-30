package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"shopapi/comm"
	"github.com/gocraft/dbr"
	"errors"
	"strconv"	
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	
	"strings"
	"fmt"
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

type ProdSku struct  {
	Id int64
	SkuNo string
	ProdId int64
	AppId string
	SoldNum int
	Price float64
	DisPrice float64
	AttrSymbolPath string
	Stock int
	Json string
}

func ProdSkuAdd(prodSku *ProdSku) (*ProdSku,error)  {

	pSku :=dao.NewProdSku()

	pSku,err :=pSku.WithProdIdAndSymbolPath(prodSku.AttrSymbolPath,prodSku.ProdId)
	if err!=nil{
		return nil,err
	}
	if pSku!=nil{
		return nil,errors.New("已存在sku!")
	}

	pSku =dao.NewProdSku()
	pSku.AppId = prodSku.AppId
	pSku.AttrSymbolPath = prodSku.AttrSymbolPath
	pSku.Price = prodSku.Price
	pSku.DisPrice = prodSku.DisPrice
	pSku.Json = prodSku.Json
	pSku.ProdId = prodSku.ProdId
	pSku.SkuNo = util.GenerUUId()
	pSku.Stock = prodSku.Stock
	id,err :=pSku.Insert()

	prodSku.Id=id

	return prodSku,err
}

func ProdSkuUpdate(prodSku *ProdSku) (*ProdSku,error) {
	pSku :=dao.NewProdSku()
	err :=pSku.UpdatePriceWithProdIdAndSymbolPath(prodSku.Price,prodSku.DisPrice,prodSku.AttrSymbolPath,prodSku.ProdId,prodSku.Stock)

	return prodSku,err
}

func ProdDetailListWith(keywords interface{},merchantId int64,flags []string,noflags []string,orderBy string,pageIndex uint64,pageSize uint64,appId string) ([]*dao.ProductDetail,error)  {

	return dao.NewProductDetail().ProdDetailListWith(keywords,merchantId,flags,noflags,orderBy,pageIndex,pageSize,appId)
}

func ProdDetailListCountWith(keywords interface{},merchantId int64,flags []string,noflags []string) (int64,error)  {

	return dao.NewProductDetail().ProdDetailListCountWith(keywords,merchantId,flags,noflags)
}

//商品详情
func ProdDetailWithProdId(prodId int64,appId string) (*dao.ProductDetail,error)  {
	product := dao.NewProduct()	
	prod,err:=product.ProductDetailWithId(prodId)
	if err!=nil{		
		return nil,err
	}
	
	prodids := make([]int64,0)
	prodids = append(prodids,prod.Id)
	
	prodSkuDao := dao.NewProdSku()
	prodSkus,err:=prodSkuDao.ProdSkuWithProdIds(prodids)
	if err!=nil{		
		return nil,err
	}
	
	prod.ProdSkus=prodSkus
	
	return prod,nil
}
func ProdDetailWithProdParentId(prodParentId int64,appId string) ([]*dao.ProductDetail,error)  {
	product := dao.NewProduct()
	return product.ProdDetailWithProdParentId(prodParentId)
}
//商品已购买数量
func ProdOrderCountWithId(prodId int64,OpenId string,Date string) (int64,error) {
	return dao.ProdOrderCountWithId(prodId,OpenId,Date)
}
//修改SKU 库存
func ProductUpdateStockWithProdId(prodId int64,stock int,soldNum int) error {
	prodids := make([]int64,0)
	prodids = append(prodids,prodId)
	
	prodSkuDao := dao.NewProdSku()

	prodSkus,err:=prodSkuDao.ProdSkuWithProdIds(prodids)
	if err!=nil{		
		return err
	}
	
	for _,prodSkud :=range prodSkus {
		err:=prodSkuDao.UpdateStockWithSkuNo(stock,prodSkud.SkuNo)
		if err!=nil{		
			return err 
		}
		err=prodSkuDao.UpdateSoldNumWithSkuNo(soldNum,prodSkud.SkuNo)
		if err!=nil{		
			return err 
		}
	}
	return dao.NewProduct().UpdateSoldNumWithSkuNo(soldNum,prodId)
}

//商品图片
func  ProdImgsWithProdId(prodId int64,appId string) ([]*ProdImgsDetailDLL,error) {
	prodImgDetail := dao.NewProdImgsDetail()
	prodImgDetals,err := prodImgDetail.ProdImgsWithProdId(prodId,appId)
	//log.Error("prodImgDetals=",prodImgDetals)
	//log.Error("size:",len(prodImgDetals),"prod_id:",prodId,"app_id:",appId)
	if err!=nil{

		return nil,err
	}
	detailDLLs := make([]*ProdImgsDetailDLL,0)
	if prodImgDetals!=nil {
		for _,detail :=range prodImgDetals{
			detailDLLs = append(detailDLLs,prodImgsDetailToDLL(detail))
		}
	}

	//log.Error("detailDLLs size:",len(detailDLLs))

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
	if prodbll.Imgs!=nil&&len(prodbll.Imgs)>0 {
		err := productImgSave(prodbll,tx)
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

func ProdUpdate(prodbll *ProdBLL) (*ProdBLL,error)  {
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	//保存商品基础信息
	err := productBaseUpdate(prodbll,tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	//修改商品图片信息
	err = productImgUpdate(prodbll,tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	//修改商品所属信息
	err = merchantProdUpdate(prodbll,tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}
	//商品分类修改
	if err:=produCategoryUpdate(prodbll,tx);err!=nil{
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

//商品分类
func CategoryWithFlags(flags []string,noflags []string,appId string) ([]*dao.Category,error)   {

	category :=dao.NewCategory()
	categories,err := category.WithFlags(flags,noflags,appId)
	return categories,err
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

//获取推荐列表
func ProductListWithRecomm(appId string,pageIndex uint64,pageSize uint64) ([]*dao.ProductDetail,int,error) {
	productDetail :=dao.NewProductDetail()
	prodList,count,err := productDetail.ProductListWithRecomm(appId,pageIndex,pageSize)
	if err!=nil {
		return nil,0,err
	}

	return prodList,count,nil
}

//根据分类获取商品
func ProductListWithCategory(appId string,categoryId int64,flags []string,noflags []string,pageIndex uint64,pageSize uint64) ([]*dao.ProductDetail,int,error)   {
	productDetail :=dao.NewProductDetail()
	prodList,count,err := productDetail.ProductListWithCategory(appId,categoryId,flags,noflags,pageIndex,pageSize)
	if err!=nil {
		return nil,0,err
	}
	return prodList,count,nil
}
func ProductListWithCategoryIsLimit(appId string,categoryId int64,flags []string,noflags []string,pageIndex uint64,pageSize uint64) ([]*dao.ProductDetail,int,error)   {
	productDetail :=dao.NewProductDetail()
	prodList,count,err := productDetail.ProductListWithCategoryIsLimit(appId,categoryId,flags,noflags,pageIndex,pageSize)
	if err!=nil {
		return nil,0,err
	}
	return prodList,count,nil
}

//查询商品指定key的属性值
func ProductAttrValues(vsearch string,attrKey string,prodId int64) ([]*dao.ProdAttrVal,error)  {
	prodAttrVal := dao.NewProdAttrVal()
	prodAttrVals,err :=prodAttrVal.WithAttrKeyStock(vsearch,attrKey,prodId)

	return prodAttrVals,err
}

func ProductSkuWithProdIdAndSymbolPath(prodId int64,symbolPath string) (*dao.ProdSku,error)  {

	prodSku :=dao.NewProdSku()
	prodSku,err :=prodSku.WithProdIdAndSymbolPath(symbolPath,prodId)

	return prodSku,err
}

//商品状态修改
func ProductStatusUpdate(status int,prodId int64) error  {

	return dao.NewProduct().UpdateStatusWithProdId(status,prodId)

}

//商品推荐
func ProductRecom(isRecom int,prodId int64) error  {

	return dao.NewProduct().UpdateRecomWithProdId(isRecom,prodId)
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

	//log.Error("prodAttrVal----------",prodAttrVal)

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
	prodSku.AttrSymbolPath = strconv.Itoa(int(prodAttrVal.Id))
	prodSku.DisPrice = product.DisPrice
	prodSku.Price = product.Price
	prodSku.Stock=1
	prodSku.SkuNo = util.GenerUUId()
	err =prodSku.InsertTx(tx)
	if err!=nil{
		tx.Rollback()
		return nil,err
	}

	err =tx.Commit()
	if err!=nil {
		log.Error(err)
		return nil,errors.New("提交失败!")
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

//商户信息修改
func merchantProdUpdate(prodbll *ProdBLL,tx *dbr.Tx) error  {

	mprod := dao.NewMerchantProd()

	return mprod.UpdateTx(prodbll.Id,prodbll.MerchantId,tx)
}
//保存商品图片
func productImgSave(prodbll *ProdBLL,tx *dbr.Tx) error  {

	imgArray := prodbll.Imgs
	for _,img :=range imgArray  {
		prodImgs :=dao.NewProdImgs()
		prodImgs.AppId = prodbll.AppId
		prodImgs.ProdId = prodbll.Id
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

func productImgUpdate(prodbll *ProdBLL,tx *dbr.Tx) error{
	imgArray := prodbll.Imgs
	err :=dao.NewProdImgs().DeleteWithIdTx(prodbll.Id,tx)
	if err!=nil{
		return err
	}
	if imgArray!=nil&&len(imgArray)>0 {
		for _,img :=range imgArray  {
			prodImgs :=dao.NewProdImgs()
			prodImgs.AppId = prodbll.AppId
			prodImgs.ProdId = prodbll.Id
			prodImgs.Url = img.Url
			prodImgs.Flag = img.Flag
			prodImgs.Json = img.Json
			err :=prodImgs.InsertTx(tx)
			if err!=nil{
				return err
			}
		}
	}


	return nil;
}

//保存商品基础信息
func productBaseSave(prodbll *ProdBLL,tx *dbr.Tx) (int64,error)  {
	//prodbll.Status = comm.PRODUCT_STATUS_NORMAL
	prodbll.Status = comm.PRODUCT_STATUS_ADUIT
	prod := prodToModel(prodbll)
	prodId,err := prod.InsertTx(tx)
	if err!=nil{
		return 0,err
	}

	return prodId,err
}

func productBaseUpdate(prodbll *ProdBLL,tx *dbr.Tx) (error)  {
	prod := prodToModel(prodbll)
	err := prod.UpdateTx(tx)
	if err!=nil{
		return err
	}

	return err
}

func produCategoryUpdate(prodbll *ProdBLL,tx *dbr.Tx) (error)  {
	prodCategory := dao.NewProdCategory()
	prodCategory.ProdId = prodbll.Id
	prodCategory.CategoryId = prodbll.CategoryId
	err := prodCategory.UpdateTx(tx)
	if err!=nil{
		return err
	}

	return err
}


func prodToModel(prodbll *ProdBLL) *dao.Product {
	prod := dao.NewProduct()
	prod.Id = prodbll.Id
	prod.Title = prodbll.Title
	prod.SubTitle = prodbll.SubTitle
	prod.AppId = prodbll.AppId
	prod.Description = prodbll.Description
	prod.DisPrice = prodbll.DisPrice
	prod.Price = prodbll.Price
	prod.Status = prodbll.Status
	prod.Json = prodbll.Json
	prod.Flag = prodbll.Flag
	prod.LimitNum = prodbll.LimitNum
	prod.ParentId = prodbll.ParentId
	prod.Goodsid = prodbll.Goodsid
	prod.IsLimit = prodbll.IsLimit

	return prod
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
//录入商品链接
func ProductAndAddLink(appId string,prodId uint64,shopurl string) error  {
	return dao.ProductAndAddLink(appId,prodId,shopurl)
}
//changeshowstate
func ProductChangeShowState(appId string,id int64,show int64) error  {	
	return dao.NewProduct().ProductChangeShowState(appId,id,show)
}
//一元购生成购买码
func ProductAndPurchaseCodesAdd(ProdPurchaseCode *dao.ProdPurchaseCode) error {
	codeMap:=make(map[int]string)
	for i := 1; i <= ProdPurchaseCode.Num; i++ {
		codeMap[i]=fmt.Sprintf("1%07d",i)
	}
	code:=make([]string,0)
	for _,v := range codeMap {
		code=append(code,v)
	}
	//============
	ProdPurchaseCode.Codes=strings.Join(code,",")
	
	return dao.ProductAndPurchaseCodesAdd(ProdPurchaseCode)
}

//一元购减去购买码
/* func ProductAndPurchaseCodesMinus(ProdPurchaseCode *dao.ProdPurchaseCode) (string,error) {
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	
	productDao:=dao.NewProduct()
	
	codes,err:=productDao.ProductAndPurchaseCodes(ProdPurchaseCode,tx)
	if err!=nil || codes==nil{
		tx.Rollback()
		return "",err
	}
	
	if(codes.Num<ProdPurchaseCode.Num){
		tx.Rollback()
		return "",errors.New("购买数量大于库存数量!")
	}	
	//============================
	s:=strings.Split(codes.Codes, ",")
	ns:=s[ProdPurchaseCode.Num:]	
	ls:=s[0:ProdPurchaseCode.Num]
	
	err=productDao.ProductAndPurchaseCodesMinus(tx,codes.Id,codes.Num,len(ns),strings.Join(ns,","))
	if err!=nil{
		tx.Rollback()
		return "",err
	}
	
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return "",err
	}
	
	return strings.Join(ls,","),nil
} */
//参与计算一元购产品计算中奖号的条数
func ProductBuyCodesWithProdId(prodId int64) ([]*dao.OrderItemPurchaseCodeRrecord,error) {
	prod,_:=dao.ProdPurchaseCodeWithProdId(prodId)
	if prod==nil {
		return nil,errors.New("产品不是一元购产品!")
	}
	return dao.OrderItemPurchaseCodesRrecordWithTime(prod.OpenTime,comm.PRODUCT_YYG_BUY_CODES)
}

func ProdSoldNumRealWithId(prodId int64) (int64,error) {
	return dao.NewProdSku().ProdSoldNumRealWithId(prodId)
}

func ProdDetailListWithIds(appId string,ids []string) ([]*dao.ProductDetail,error)  {
	return dao.NewProductDetail().ProdDetailListWithIds(appId,ids)
}





































