package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/service"
	"strconv"
	"net/http"
)

type ProductParam struct  {
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
	//图片编号集合
	ImgNos string `json:"imgnos"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//附加数据
	Json  string `json:"json"`
}

type ProdImgsDetailDto struct  {
	//图片编号
	ImgNo string `json:"img_no"`
	//产品ID
	ProdId int64 `json:"prod_id"`
	AppId string `json:"app_id"`
	Url string `json:"url"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}
type ProductListDto struct  {
	Id string `json:"id"`
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

/**
添加商品
 */
func ProductAdd(c *gin.Context)  {

	appId,err := CheckAuth(c)
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
	if param.ImgNos=="" {
		util.ResponseError400(c.Writer,"请上传商品图片")
		return
	}

	if param.Description==""{
		util.ResponseError400(c.Writer,"请输入商品描述")
		return
	}
	mid,err := strconv.Atoi(midstr)
	param.MerchantId = int64(mid)
	param.AppId = appId

	prodBll := productParamToDLL(param)
	err =service.ProdAdd(prodBll)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

/**
商品列表
 */
func ProductListWithCategory(c *gin.Context)  {

	appId,err :=CheckAuth(c)
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
	prodListDtos :=make([]*ProductListDto,0)
	if prodList!=nil {

		for _,prodResultBll :=range prodList {
			prodListDtos = append(prodListDtos,productResultDLLToDto(prodResultBll))
		}
	}

	c.JSON(http.StatusOK,prodListDtos)
}

func ProdImgsWithProdId(c *gin.Context)  {
	appId,err :=CheckAuth(c)
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

func prodImgsDetailDLLToDto(dll *service.ProdImgsDetailDLL) *ProdImgsDetailDto{

	dto :=&ProdImgsDetailDto{}
	dto.Url = dll.Url
	dto.ProdId = dll.ProdId
	dto.AppId = dll.AppId
	dto.Flag = dll.Flag
	dto.ImgNo = dll.ImgNo
	dto.Json = dll.Json

	return dto
}

func productParamToDLL(param *ProductParam) *service.ProdBLL {

	prodBll := &service.ProdBLL{}
	prodBll.AppId = param.AppId
	prodBll.MerchantId = param.MerchantId
	prodBll.CategoryId = param.CategoryId
	prodBll.Description = param.Description
	prodBll.Price = param.Price
	prodBll.ImgNos = param.ImgNos
	prodBll.DisPrice = param.DisPrice
	prodBll.Title = param.Title
	prodBll.Json  = param.Json


	return prodBll
}

func productResultDLLToDto(dll *service.ProductResultDLL) *ProductListDto  {

	dto :=&ProductListDto{}
	dto.Description=dll.Description
	dto.DisPrice=dll.DisPrice
	dto.Id = dll.Id
	dto.Json = dll.Json
	dto.Price = dll.Price
	dto.Title = dll.Title

	return dto
}