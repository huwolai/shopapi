package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/dao"
	"shopapi/comm"
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/page"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
)


type DistributionProductDetail struct {
	//商品ID
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品标题
	Title string `json:"title"`
	SubTitle string `json:"sub_title"`
	Description string `json:"description"`	
	//分销编号
	DbnNo string `json:"dbn_no"`
	//商品价格
	Price float64 `json:"price"`
	SoldNum float64 `json:"sold_num"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`
	//是否已添加分销
	Added int `json:"added"`
	//佣金比例
	CsnRate float64 `json:"csn_rate"`
	//商品佣金
	CsnAmount float64 `json:"csn_amount"`
	//商品状态
	Status int `json:"status"`
	//分销ID
	DistributionId int64 `json:"distribution_id"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//商户名称
	MerchantName string `json:"merchant_name"`
	Json string `json:"json"`
	//商品图片集合
	ProdImgs []*DisProdImgsDetailDto `json:"prod_imgs"`
}

//分销商品详情
type DistributionProductDetail2 struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品ID
	ProdId int64 `json:"prod_id"`
	//商品标题
	Title string `json:"title"`
	//商品价格
	Price float64 `json:"price"`
	//折扣价格
	DisPrice float64 `json:"dis_price"`
	//佣金比例
	CsnRate float64 `json:"csn_rate"`
	//分销编号
	DbnNo string `json:"dbn_no"`
	CreateTime string `json:"create_time"`
}

type DistributionProduct struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	ProdId int64 `json:"prod_id"`
	MerchantId int64 `json:"merchant_id"`
	//佣金比例
	CsnRate float64 `json:"csn_rate"`
}

type DisProdImgsDetailDto struct  {
	//产品ID
	ProdId int64 `json:"prod_id"`
	AppId string `json:"app_id"`
	Url string `json:"url"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

//获取正在参与分销的商品
func DistributionProducts(c *gin.Context)  {

	appId :=security.GetAppId2(c.Request)

	openId := security.GetOpenId(c.Request)
	added :=c.Query("added")
	list,err := service.DistributionProducts(added,openId,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}

	details :=make([]*DistributionProductDetail,0)
	if list!=nil{
		for _,detail :=range list {
			details = append(details,distributionProductDetailToA(detail))
		}
	}
	c.JSON(http.StatusOK,details)
}

//商户分销的商品
func DistributionWithMerchant(c *gin.Context)  {

	appId :=security.GetAppId2(c.Request)
	merchantId := c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID格式有误!")
		return
	}
	list,err := service.DistributionWithMerchant(imerchantId,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}

	details :=make([]*DistributionProductDetail,0)
	if list!=nil{
		for _,detail :=range list {
			details = append(details,distributionProductDetailToA(detail))
		}
	}
	c.JSON(http.StatusOK,details)
}

//添加分销
func DistributionProductAdd(c *gin.Context)  {
	appId :=security.GetAppId2(c.Request)
	openId := c.Param("open_id")
	distributionId := c.Param("id")
	idistributionId,err :=strconv.ParseInt(distributionId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"id格式有误")
		return
	}

	ud,err :=service.DistributionProductAdd(idistributionId,openId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加分销失败!")
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":ud.Code,
	})
}

func ProductJoinOrUpdateDistribution(c *gin.Context)  {
	var param *DistributionProduct
	err :=c.BindJSON(&param)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	appId :=security.GetAppId2(c.Request)
	if param.Id!=0 {
		err = service.ProductUpdateDistribution(param.Id,param.CsnRate,appId)
	}else{
		err = service.ProductJoinDistribution(param.ProdId,param.CsnRate,appId)

	}
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加或修改失败!")
		return
	}

	util.ResponseSuccess(c.Writer)
}

//取消分销
func DistributionProductCancel(c *gin.Context) {
	appId :=security.GetAppId2(c.Request)
	openId := c.Param("open_id")
	distributionId := c.Param("id")
	idistributionId,err :=strconv.ParseInt(distributionId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"id格式有误")
		return
	}
	err =service.DistributionProductCancel(idistributionId,openId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"取消失败!")
		return
	}

	util.ResponseSuccess(c.Writer)
}

func DistributionProductDelete(c *gin.Context)  {
	sid :=c.Param("id")
	id,err :=strconv.ParseInt(sid,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"id格式有误")
		return
	}

	err =service.DistributionProductDelete(id)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"删除失败!")
		return
	}

	util.ResponseSuccess(c.Writer)
}

/**
 查询分销商品信息
 */
func DistributionProductWithId(c *gin.Context)  {
	sid :=c.Param("id")
	id ,err :=strconv.ParseInt(sid,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误！")
		return
	}
	result,err := service.DistributionProductWithId(id)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败！")
		return
	}

	c.JSON(http.StatusOK,distributionProductToA(result))
}

func DistributionWith(c *gin.Context)  {

	keyword :=c.Query("keyword")
	flags,noflags := GetFlagsAndNoFlags(c.Query("flags"),c.Query("noflags"))
	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))

	items,err :=service.DistributionWith(keyword,pIndex,pSize,noflags,flags)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}

	results := make([]*DistributionProductDetail2,0)
	var total int64
	if items!=nil&&len(items)>0{
		for _,item :=range items {
			results = append(results,distributionProductDetail2ToA(item))
		}

		total,err = service.DistributionWithCount(keyword,noflags,flags)
		if err!=nil{
			log.Error(err)
			util.ResponseError400(c.Writer,"查询数量失败!")
			return
		}
	}

	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(total),results))

}

func distributionProductToA(model *dao.DistributionProduct) *DistributionProduct  {

	a:=&DistributionProduct{}
	a.AppId = model.AppId
	a.CsnRate = model.CsnRate
	a.Id = model.Id
	a.MerchantId = model.MerchantId
	a.ProdId = model.ProdId

	return a
}

func distributionProductDetail2ToA(model *dao.DistributionProductDetail2) *DistributionProductDetail2  {

	a :=&DistributionProductDetail2{}
	a.AppId = model.AppId
	a.CreateTime = qtime.ToyyyyMMddHHmm(model.CreateTime)
	a.CsnRate = model.CsnRate
	a.DisPrice=model.DisPrice
	a.Price = model.Price
	a.ProdId =model.ProdId
	a.Title = model.Title
	a.Id = model.Id

	return a
}

func distributionProductDetailToA(model *dao.DistributionProductDetail) *DistributionProductDetail  {

	a :=&DistributionProductDetail{}
	a.AppId = model.AppId
	a.CsnRate = model.CsnRate
	a.CsnAmount = comm.Floor(model.DisPrice*model.CsnRate,2)
	a.DisPrice = model.DisPrice
	a.Json = model.Json
	a.MerchantId = model.MerchantId
	a.MerchantName = model.MerchantName
	a.Id = model.Id
	a.Price = model.Price
	a.Title = model.Title
	a.Status = model.Status
	a.Added = model.Added
	a.DbnNo = model.DbnNo
	a.DistributionId = model.DistributionId
		
	a.SubTitle = model.SubTitle
	a.Description = model.Description
	a.SoldNum = model.SoldNum
	
	
	
	if model.ProdImgs!=nil{
		detailDtos :=make([]*DisProdImgsDetailDto,0)

		for _,prodimg :=range model.ProdImgs {
			detailDtos = append(detailDtos,prodImgsDetailToA(prodimg))
		}
		a.ProdImgs=detailDtos
	}

	return a
}

func prodImgsDetailToA(model *dao.ProdImgsDetail) *DisProdImgsDetailDto  {

	dto :=&DisProdImgsDetailDto{}
	dto.AppId = model.AppId
	dto.Flag = model.Flag
	dto.ProdId = model.ProdId
	dto.Url = model.Url

	return dto
}

