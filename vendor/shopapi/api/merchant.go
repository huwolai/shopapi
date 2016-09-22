package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"net/http"
	"shopapi/dao"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"strings"
	"gitlab.qiyunxin.com/tangtao/utils/page"
	"os"
	"shopapi/comm"
)


type MerchantDto struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Status int `json:"status"`
	Json string `json:"json"`
	HasAvatar int `json:"has_avatar"`
	Address string `json:"address"`
	//权重
	Weight int `json:"weight"`
	//手机号
	Mobile string `json:"mobile"`
	CoverDistance float64 `json:"cover_distance"`
	//经度
	Longitude float64 `json:"longitude"`
	//维度
	Latitude float64 `json:"latitude"`
}
type MerchantDetailParam struct  {
	Name string `json:"name"`
	OpenId string `json:"open_id"`
	Json string `json:"json"`
	Address string `json:"address"`
	Id int64 `json:"id"`
	//经度
	Longitude float64 `json:"longitude"`
	//纬度
	Latitude float64 `json:"latitude"`
	//覆盖距离 (单位米)
	CoverDistance float64 `json:"cover_distance"`
	AppId string `json:"app_id"`
	//商户图片
	Imgs []MerchantImgDto
}

//商户主要图片DTO
type MerchantImgDto struct  {
	Id int64 `json:"id"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Url string `json:"url"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}


type MerchantDetailDto struct  {
	Id string `json:"id"`
	Name string `json:"name"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	Status int `json:"status"`
	Json string `json:"json"`
	//商户地址
	Address string `json:"address"`
	//权重
	Weight int `json:"weight"`
	CoverDistance float64 `json:"cover_distance"`
	//距离(单位米)
	Distance float64 `json:"distance"`

}

type MerchantOpenDto struct  {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	MerchantId int64 `json:"merchant_id"`
	IsOpen int
	OpenTimeStart string `json:"open_time_start"`
	OpenTimeEnd string `json:"open_time_end"`
}

type MerchantServiceTimeDto struct  {
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//服务时间
	Stime []string `json:"stimes"`

}

func MerchantUpdate(c *gin.Context)  {
	openId,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	appId,err :=security.GetAppId(c.Request)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}


	var param MerchantDetailParam
	err =c.BindJSON(&param)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	if param.Id==0 {
		util.ResponseError400(c.Writer,"商户ID不能为空")
		return
	}
	param.OpenId = openId
	param.AppId =appId

	err =service.MerchantUpdate(merchantDetailParamToDll(param))
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

func MerchantAdd(c *gin.Context)  {

	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	appId,err :=security.GetAppId(c.Request)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	var param MerchantDetailParam
	err =c.BindJSON(&param)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	openId := c.Param("open_id")
	if openId=="" {
		util.ResponseError400(c.Writer,"用户open_id不能为空!")
		return
	}
	param.OpenId = openId
	if param.Name=="" {
		util.ResponseError400(c.Writer,"名称不能为空!")
		return
	}

	param.AppId =appId

	mdll,err :=service.MerchantAdd(merchantDetailParamToDll(param))
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	param.Id =mdll.Id
	c.JSON(http.StatusOK,param)

}

func MerchantWith(c *gin.Context)  {

	flags := c.Query("flags")
	noflags := c.Query("noflags")
	status :=c.Query("status")
	orderBy := c.Query("order_by")

	appId :=security.GetAppId2(c.Request)
	flagsArray,noflagArray := GetFlagsAndNoFlags(flags,noflags)

	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))

	merchants,err := service.MerchantWith(flagsArray,noflagArray,status,orderBy,pIndex,pSize,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败!")
		return
	}

	count,err := service.MerchantCountWith(flagsArray,noflagArray,status,orderBy,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询数量失败!")
		return
	}
	merchantsDto :=make([]*MerchantDto,0)
	if merchants!=nil{
		for _,merchant :=range merchants {
			merchantsDto = append(merchantsDto,merchantToDto(merchant))
		}
	}
	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(count),merchantsDto))
}

func MerchantOpenWithMerchantId(c *gin.Context)  {
	merchantId := c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}

	merchantOpen,err :=service.MerchantOpenWithMerchantId(imerchantId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if merchantOpen==nil{
		util.ResponseError400(c.Writer,"商户没有设置营业时间!")
		return
	}

	c.JSON(http.StatusOK,merchantOpenToDto(merchantOpen))
}

func MerchantWithId(c *gin.Context)  {
	appId :=security.GetAppId2(c.Request)
	id := c.Param("merchant_id")
	iid,err := strconv.ParseInt(id,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"ID有误!")
		return
	}

	merchant,err := service.MerchantWithId(iid,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	if merchant!=nil{
		c.JSON(http.StatusOK,merchantToDto(merchant))
		return
	}

	util.ResponseError400(c.Writer,"没有找到信息!")
}

//是否是商户
func MerchantIs(c *gin.Context)  {
	appId :=security.GetAppId2(c.Request)
	openId := c.Param("open_id")
	merchant,err := service.MerchantWithOpenId(openId,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if merchant!=nil&&merchant.Status==comm.MERCHANT_STATUS_NORMAL{

		c.JSON(http.StatusOK,gin.H{
			"is_merchant": 1,
			"merchant_id":merchant.Id,
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"is_merchant": 0,
		})
	}
}


func MerchantWithOpenId(c *gin.Context)  {

	appId :=security.GetAppId2(c.Request)
	openId := c.Param("open_id")

	merchant,err := service.MerchantWithOpenId(openId,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if merchant!=nil{
		_,err :=os.Open(MERCHANT_IMG_PATH+"/"+merchant.OpenId)
		dto :=merchantToDto(merchant)
		if (!os.IsNotExist(err)) {
			dto.HasAvatar = 1
		}
		c.JSON(http.StatusOK,dto)
		return
	}

	util.ResponseError400(c.Writer,"没有找到信息!")


}

func MerchantServiceTimeGet(c *gin.Context)  {
	merchantId :=c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}

	merchantServiceTimes,err := service.MerchantServiceTimeGet(imerchantId)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if merchantServiceTimes==nil{
		dto :=&MerchantServiceTimeDto{}
		dto.MerchantId = imerchantId
		c.JSON(http.StatusOK,dto)
		return
	}

	c.JSON(http.StatusOK,merchantServiceTimesToDto(merchantServiceTimes,imerchantId))
}

func MerchantServiceTimeAdd(c *gin.Context)  {
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	var params *MerchantServiceTimeDto
	err =c.BindJSON(&params)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	merchantId :=c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}
	err =service.MerchantServiceTimeAdd(imerchantId,params.Stime)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

//附近商户
func MerchatNear(c *gin.Context)  {
	appId := security.GetAppId2(c.Request)

	longitude :=c.Query("longitude")
	latitude :=c.Query("latitude")
	openId := security.GetOpenId(c.Request)
	
	
	
	pIndex,pSize := page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	

	flongitude,_ :=strconv.ParseFloat(longitude,20)
	flatitude,_  :=strconv.ParseFloat(latitude,20)
	
	mDetailList,err := service.MerchantNear(flongitude,flatitude,openId,appId,pIndex,pSize)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	mDetailDtos :=make([]*MerchantDetailDto,0)
	if mDetailList!=nil {

		for _,mDetail :=range mDetailList {
			mDetailDtos = append(mDetailDtos,merchantDetailToDto(mDetail))
		}
	}

	c.JSON(http.StatusOK,mDetailDtos)
}

//附近商户搜索 可提供服务的厨师
func MerchatNearSearch(c *gin.Context)  {
	appId := security.GetAppId2(c.Request)

	longitude :=c.Query("longitude")
	latitude :=c.Query("latitude")
	openId := security.GetOpenId(c.Request)
	
	pIndex,pSize := page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))	

	flongitude,_ :=strconv.ParseFloat(longitude,20)
	flatitude,_  :=strconv.ParseFloat(latitude,20)
	
	
	serviceTime :=c.Query("service_time")
	//搜索不能为空
	if serviceTime=="" {
		util.ResponseError400(c.Writer,"服务时间不能为空!")
		return
	}
	
	mDetailList,err := service.MerchantNearSearch(flongitude,flatitude,openId,appId,pIndex,pSize,serviceTime)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	mDetailDtos :=make([]*MerchantDetailDto,0)
	if mDetailList!=nil {

		for _,mDetail :=range mDetailList {
			mDetailDtos = append(mDetailDtos,merchantDetailToDto(mDetail))
		}
	}

	c.JSON(http.StatusOK,mDetailDtos)
}

func MerchantProds(c *gin.Context)  {

	merchantId := c.Param("merchant_id")
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID格式有误!")
		return
	}
	appId := security.GetAppId2(c.Request)
	flags := c.Query("flags")
	noflags := c.Query("noflags")
	var flagsArray []string
	if flags != "" {
		flagsArray = strings.Split(flags,",")
	}
	var noflagsArray []string
	if noflags!="" {
		noflagsArray = strings.Split(noflags,",")
	}

	prodList,err:=service.MerchantProds(imerchantId,appId,flagsArray,noflagsArray)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
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

func MerchantImgWithMerchantId(c *gin.Context)  {
	flags :=c.Query("flags")
	merchantId := c.Param("merchant_id")
	imerchantId,err :=strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}

	var flagsArray []string;
	if flags!="" {
		flagsArray = strings.Split(flags,",")
	}
	appId := security.GetAppId2(c.Request)

	merchantImgs,err :=service.MerchantImgWithMerchantId(imerchantId,flagsArray,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	imgddtos := make([]*MerchantImgDto,0)
	if merchantImgs!=nil {
		for _,mImgsModel :=range merchantImgs{
			imgddtos = append(imgddtos,merchantImgToDto(mImgsModel))
		}
	}
	c.JSON(http.StatusOK,imgddtos)
}

//根据图片标记查询商户图片
func MerchantImgWithFlag(c *gin.Context)  {

	appId :=security.GetAppId2(c.Request)
	flags :=c.Query("flags")
	mopenId := c.Param("open_id")

	merchantImgs,err := service.MerchantImgWithFlag(flags,mopenId,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	imgddtos := make([]*MerchantImgDto,0)
	if merchantImgs!=nil {
		for _,mImgsModel :=range merchantImgs{
			imgddtos = append(imgddtos,merchantImgToDto(mImgsModel))
		}
	}
	c.JSON(http.StatusOK,imgddtos)

}

func MerchantAudit(c *gin.Context){
	_,err := security.CheckUserAuth(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,err.Error())
		return
	}
	merchantId := c.Param("merchant_id")
	appId := security.GetAppId2(c.Request)
	imerchantId,err := strconv.ParseInt(merchantId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"商户ID有误!")
		return
	}
	err =service.MerchantAudit(imerchantId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.ResponseSuccess(c.Writer)
}

func merchantImgToDto(model *dao.MerchantImgs) *MerchantImgDto   {
	dto := &MerchantImgDto{}
	dto.MerchantId = model.MerchantId
	dto.Json = model.Json
	dto.AppId = model.AppId
	dto.Flag = model.Flag
	dto.Url = model.Url
	dto.Id = model.Id

	return dto
}

func merchantDetailToDto(model *dao.MerchantDetail) *MerchantDetailDto  {
	dto :=&MerchantDetailDto{}
	dto.AppId=model.AppId
	dto.Distance = model.Distance
	dto.Json = model.Json
	dto.Name = model.Name
	dto.OpenId = model.OpenId
	dto.Status = model.Status
	dto.Address = model.Address
	dto.Weight = model.Weight
	dto.Id = model.Id
	dto.CoverDistance =  model.CoverDistance * 1000

	return dto

}

func merchantDetailParamToDll(param MerchantDetailParam)  *service.MerchantDetailDLL {

	dll :=&service.MerchantDetailDLL{}
	dll.OpenId = param.OpenId
	dll.Name = param.Name
	dll.Json = param.Json
	dll.AppId = param.AppId
	dll.Address = param.Address
	dll.CoverDistance = param.CoverDistance * 1000
	dll.Latitude = param.Latitude
	dll.Longitude = param.Longitude
	dll.Id = param.Id

	if param.Imgs!=nil {
		imgdlls := make([]service.MerchantImgDLL,0)
		for _,imgDto :=range param.Imgs  {
			imgDto.AppId = dll.AppId
			imgdlls = append(imgdlls,merchantImgToDLL(imgDto))
		}

		dll.Imgs=imgdlls
	}

	return dll
}

func merchantImgToDLL(model MerchantImgDto) service.MerchantImgDLL  {
	dll :=service.MerchantImgDLL{}
	dll.Url = model.Url
	dll.Flag = model.Flag
	dll.AppId = model.AppId
	dll.Json = model.Json
	dll.MerchantId = model.MerchantId

	return dll
}

func merchantToDto(model *dao.Merchant)  *MerchantDto {
	dto:=&MerchantDto{}
	dto.Json = model.Json
	dto.Address = model.Address
	dto.AppId = model.AppId
	dto.CoverDistance = model.CoverDistance * 1000
	dto.Latitude=model.Latitude
	dto.Longitude = model.Longitude
	dto.Name = model.Name
	dto.Id = model.Id
	dto.OpenId = model.OpenId
	dto.Weight = model.Weight
	dto.Status = model.Status

	if len(model.Mobile)==11 {
		dto.Mobile =  strings.Replace(model.Mobile,model.Mobile[3:8],"*****",1)
	}

	return dto
}

func merchantOpenToDto(model *dao.MerchantOpen) *MerchantOpenDto  {
	dto :=&MerchantOpenDto{}
	dto.Id = model.Id
	dto.AppId = model.AppId
	dto.IsOpen = model.IsOpen
	dto.MerchantId = model.MerchantId
	dto.OpenTimeEnd = model.OpenTimeEnd
	dto.OpenTimeStart = model.OpenTimeStart

	return dto
}

func merchantServiceTimesToDto(models []*dao.MerchantServiceTime,merchantId int64) *MerchantServiceTimeDto  {

	dto :=&MerchantServiceTimeDto{}
	stimes :=make([]string,0)
	for _,model :=range models  {
		stimes = append(stimes,model.Stime)
	}
	dto.MerchantId = merchantId
	dto.Stime = stimes

	return dto
}