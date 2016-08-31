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
)


type DistributionProductDetail struct {
	//商品ID
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//商品标题
	Title string `json:"title"`
	//商品价格
	Price float64 `json:"price"`
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
	a.DistributionId = model.DistributionId
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


