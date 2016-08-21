package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"strconv"
	"net/http"
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"strings"
)

type ProductParam struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品标题
	Title string `json:"title"`
	//描述
	Description string `json:"description"`
	//类别ID
	CategoryId int64 `json:"category_id"`
	//商品价格
	Price float64 `json:"price"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`
	//图片集合
	Imgs []ProductImgParam `json:"imgs"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//附加数据
	Json  string `json:"json"`
}

type ProductImgParam struct {
	Flag string
	Url string
	Json string
	ProdId int64
}


type ProductListDto struct  {
	Id int64 `json:"id"`
	//商品标题
	Title string `json:"title"`
	//描述
	Description string `json:"description"`
	//商品价格
	Price float64 `json:"price"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`

	Json string `json:"json"`
}

type ProductBaseDto struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品标题
	Title string `json:"title"`
	//商品描述
	Description string `json:"description"`
	//商品价格
	Price float64 `json:"price"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`
	//是否推荐
	IsRecom int `json:"is_recom"`
	//商品状态
	Status int `json:"status"`
	//附加数据
	Json string `json:"json"`

}

type ProductDetailDto struct {
	//商品ID
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品标题
	Title string `json:"title"`
	//商品价格
	Price float64 `json:"price"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`
	//商品状态
	Status int `json:"status"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//商户名称
	MerchantName string `json:"merchant_name"`
	Json string `json:"json"`
	//商品图片集合
	ProdImgs []*ProdImgsDetailDto `json:"prod_imgs"`
}

type ProdImgsDetailDto struct  {
	//产品ID
	ProdId int64 `json:"prod_id"`
	AppId string `json:"app_id"`
	Url string `json:"url"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

type ProdAttrValDto struct  {
	Id int64 `json:"id"`
	ProdId int64 `json:"prod_id"`
	AttrKey string `json:"attr_key"`
	AttrValue string `json:"attr_value"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

type CategoryDto struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Icon string `json:"icon"`
	Flag string `json:"flag"`
	Json string `json:"json"`

}

/**
添加商品
 */
func ProductAdd(c *gin.Context)  {

	_,err := CheckUserAuth(c)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	var param *ProductParam
	err = c.BindJSON(&param)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}

	log.Debug(*param)

	midstr := c.Param("merchant_id")

	if midstr==""{
		util.ResponseError400(c.Writer,"商户ID不能为空!")
		return
	}
	if param.Title=="" {
		util.ResponseError400(c.Writer,"商品标题不能为空")
		return
	}
	if param.CategoryId==0 {
		util.ResponseError400(c.Writer,"分类ID不能为空!")
		return
	}
	if param.Price<=0 {
		util.ResponseError400(c.Writer,"请输入商品价格!")
		return
	}
	if param.Imgs==nil {
		util.ResponseError400(c.Writer,"请上传商品图片")
		return
	}

	if param.Description==""{
		util.ResponseError400(c.Writer,"请输入商品描述")
		return
	}
	mid,err := strconv.Atoi(midstr)
	param.MerchantId = int64(mid)
	param.AppId = security.GetAppId2(c.Request)

	prodBll := productParamToBLL(param)
	prodBll,err =service.ProdAdd(prodBll)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	param.Id = prodBll.Id
	c.JSON(http.StatusOK,param)
}

//商品推荐列表
func ProductListWithRecomm(c *gin.Context)  {
	appId,err :=CheckAppAuth(c)
	if err!=nil {
		util.ResponseError400(c.Writer,"校验失败!")
		return
	}
	prodList,err := service.ProductListWithRecomm(appId)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}
	prodListDtos :=make([]*ProductDetailDto,0)
	if prodList!=nil {

		for _,prodDetail :=range prodList {
			prodListDtos = append(prodListDtos,productDetailToDto(prodDetail))
		}
	}

	c.JSON(http.StatusOK,prodListDtos)
}

/**
商品列表(根据分类查询)
 */
func ProductListWithCategory(c *gin.Context)  {

	appId,err :=CheckAppAuth(c)
	if err!=nil {
		util.ResponseError400(c.Writer,"校验失败!")
		return
	}
	categoryId :=c.Param("category_id")

	icategoryId,_ := strconv.Atoi(categoryId)

	prodList,err := service.ProductListWithCategory(appId,int64(icategoryId))
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}
	prodListDtos :=make([]*ProductDetailDto,0)
	if prodList!=nil {

		for _,prodDetail :=range prodList {
			prodListDtos = append(prodListDtos,productDetailToDto(prodDetail))
		}
	}

	c.JSON(http.StatusOK,prodListDtos)
}

//商品详情
func ProdDetailWithProdId(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"校验失败!")
		return
	}

	prodId := c.Param("prod_id")
	iprodId,_ := strconv.ParseInt(prodId,10,64)
	appId := security.GetAppId2(c.Request)
	product,err := service.ProdDetailWithProdId(iprodId,appId)

	if product==nil {
		util.ResponseError400(c.Writer,"商品没找到!")
		return
	}
	c.JSON(http.StatusOK,productToDto(product))
}

//商品图片
func ProdImgsWithProdId(c *gin.Context)  {
	appId,err :=CheckAppAuth(c)
	if err!=nil {
		util.ResponseError400(c.Writer,"校验失败!")
		return
	}
	prodId := c.Param("prod_id")
	iprodId,_ := strconv.Atoi(prodId)

	dlls,err := service.ProdImgsWithProdId(int64(iprodId),appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	detailDtos := make([]*ProdImgsDetailDto,0)
	if dlls!=nil{
		for _,dll :=range dlls  {

			detailDtos = append(detailDtos,prodImgsDetailDLLToDto(dll))
		}
	}
	c.JSON(http.StatusOK,detailDtos)
}



func ProductAndAttrAdd(c *gin.Context) {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"校验失败!")
		return
	}
	param := &service.ProdAndAttrDto{}
	err =c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	prodId :=c.Param("prod_id")

	if prodId=="" {
		util.ResponseError400(c.Writer,"商品ID不能为空!")
		return
	}
	if param.AttrValue=="" {
		util.ResponseError400(c.Writer,"属性值不能为空!")
		return
	}

	iprodId,err := strconv.ParseInt(prodId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商品ID格式有误!")
		return
	}

	param.ProdId=iprodId

	param.AppId = security.GetAppId2(c.Request)

	dto,err :=service.ProductAndAttrAdd(param)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	c.JSON(http.StatusOK,dto)
}



func ProductAttrValues(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusUnauthorized,"校验失败!")
		return
	}

	attrKey :=c.Param("attr_key")
	vsearch  :=c.Query("vsearch")
	prodId :=c.Param("prod_id")

	if attrKey==""{
		util.ResponseError400(c.Writer,"属性key不能为空")
		return
	}

	iprodId,err := strconv.ParseInt(prodId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商品ID错误!")
		return
	}

	prodAttrVals,err  :=service.ProductAttrValues(vsearch,attrKey,iprodId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询商品属性失败!")
		return
	}

	prodAttrDtos :=make([]*ProdAttrValDto,0)
	if prodAttrVals!=nil{

		for _,prodAttrVal :=range prodAttrVals {
			prodAttrDtos = append(prodAttrDtos,prodAttrValToDto(prodAttrVal))
		}
	}
	c.JSON(http.StatusOK,prodAttrDtos)
}

func CategoryWithFlags(c *gin.Context)  {

	flags := c.Query("flags")
	noflags := c.Query("noflags")

	var flagsArray []string
	if flags!=""{
		flagsArray =strings.Split(flags,",")
	}
	var  noflagsArray []string
	if noflags!="" {
		noflagsArray =strings.Split(noflags,",")
	}
	appId :=security.GetAppId2(c.Request)
	categories,err :=service.CategoryWithFlags(flagsArray,noflagsArray,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	dtos :=make([]*CategoryDto,0)
	if categories!=nil {

		for _,cateogory :=range categories {
			dtos = append(dtos,categoryToDto(cateogory))
		}
	}
	c.JSON(http.StatusOK,dtos)

}

func categoryToDto(model *dao.Category) *CategoryDto  {

	dto :=&CategoryDto{}
	dto.Id = model.Id
	dto.Flag = model.Flag
	dto.AppId = model.AppId
	dto.Description = model.Description
	dto.Icon = model.Icon
	dto.Title = model.Title
	dto.Json = model.Json

	return dto
}

func prodImgsDetailDLLToDto(dll *service.ProdImgsDetailDLL) *ProdImgsDetailDto{

	dto :=&ProdImgsDetailDto{}
	dto.Url = dll.Url
	dto.ProdId = dll.ProdId
	dto.AppId = dll.AppId
	dto.Flag = dll.Flag
	dto.Json = dll.Json

	return dto
}

func productParamToBLL(param *ProductParam) *service.ProdBLL {

	prodBll := &service.ProdBLL{}
	prodBll.Id = param.Id
	prodBll.AppId = param.AppId
	prodBll.MerchantId = param.MerchantId
	prodBll.CategoryId = param.CategoryId
	prodBll.Description = param.Description
	prodBll.Price = param.Price
	prodBll.DisPrice = param.DisPrice
	prodBll.Title = param.Title
	prodBll.Json  = param.Json

	imgsparams  := param.Imgs
	if imgsparams!=nil {
		imgBllArray :=make([]service.ProdImgBLL,0)
		for _,imgparam :=range imgsparams {
			imgBllArray = append(imgBllArray,productImgParamToBLL(imgparam))
		}
		prodBll.Imgs = imgBllArray
	}

	return prodBll
}

func productImgParamToBLL(param ProductImgParam) service.ProdImgBLL  {
	bll :=service.ProdImgBLL{}
	bll.Json = param.Json
	bll.Flag = param.Flag
	bll.ProdId=param.ProdId
	bll.Url = param.Url

	return bll
}

func productDetailToDto(model *dao.ProductDetail) *ProductDetailDto  {

	dto :=&ProductDetailDto{}
	dto.Id = model.Id
	dto.DisPrice=model.DisPrice
	dto.Json = model.Json
	dto.Title = model.Title
	dto.AppId = model.AppId
	dto.MerchantId = model.MerchantId
	dto.MerchantName = model.MerchantName
	dto.Price = model.Price

	if model.ProdImgs!=nil{
		detailDtos :=make([]*ProdImgsDetailDto,0)

		for _,prodimg :=range model.ProdImgs {
			detailDtos = append(detailDtos,prodImgsDetailToDto(prodimg))
		}
		dto.ProdImgs=detailDtos
	}

	return dto
}

func prodImgsDetailToDto(model *dao.ProdImgsDetail) *ProdImgsDetailDto  {

	dto :=&ProdImgsDetailDto{}
	dto.AppId = model.AppId
	dto.Flag = model.Flag
	dto.ProdId = model.ProdId
	dto.Url = model.Url

	return dto
}

func productToDto(model *dao.Product) *ProductBaseDto {

	dto :=&ProductBaseDto{}
	dto.AppId = model.AppId
	dto.Description =model.Description
	dto.DisPrice = model.DisPrice
	dto.Id = model.Id
	dto.IsRecom = model.IsRecom
	dto.Json = model.Json
	dto.Price = model.Price
	dto.Status = model.Status
	dto.Title = model.Title

	return dto
}

func prodAttrValToDto(model *dao.ProdAttrVal) *ProdAttrValDto {
	dto :=&ProdAttrValDto{}
	dto.AttrKey = model.AttrKey
	dto.AttrValue = model.AttrValue
	dto.Id = model.Id
	dto.Flag  = model.Flag
	dto.ProdId = model.ProdId
	dto.Json = model.Json

	return dto
}