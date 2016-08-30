package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"shopapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/dao"
	"net/http"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Favorites struct {
	ObjId int64 `json:"obj_id"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	CoverImg string `json:"cover_img"`
	Remark string `json:"remark"`
	Type int `json:"type"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

func FavoritesAdd(c *gin.Context)  {

	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)
	var param *Favorites
	err :=c.BindJSON(&param)
	if err!=nil{
		util.ResponseError400(c.Writer,"参数有误!")
		return
	}
	param.OpenId = openId
	param.AppId = appId
	err =service.FavoritesAdd(favoritesToS(param))
	if err!=nil{
		util.ResponseError400(c.Writer,"添加收藏失败!")
		return
	}
	 util.ResponseSuccess(c.Writer)
}

func FavoritesGet(c *gin.Context)  {
	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	list,err := service.FavoritesGet(openId,appId)
	if err!=nil{
		util.ResponseError400(c.Writer,"查询收藏失败!")
		return
	}

	 fs :=make([]*Favorites,0)
	if list!=nil{
		for _,f :=range list {
			fs = append(fs,favoritesToA(f))
		}
	}
	c.JSON(http.StatusOK,fs)
}

func FavoritesIsExist(c *gin.Context)  {
	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	objId :=c.Query("obj_id")
	stype := c.Query("type")
	iobjId,err :=strconv.ParseInt(objId,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"obj_id格式有误!")
		return
	}
	itype,err :=strconv.Atoi(stype)
	if err!=nil{
		util.ResponseError400(c.Writer,"type格式有误!")
		return
	}
	favor,err :=service.FavoritesIsExist(iobjId,itype,appId,openId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询收藏失败!")
		return
	}

	isf :=0 //是否收藏
	var fid int64 //收藏ID
	if favor!=nil {
		isf=1
		fid = favor.Id
	}
	c.JSON(http.StatusOK,gin.H{
		"is_favorited":isf,
		"id": fid,
	})
}

func FavoritesDelete(c *gin.Context)  {
	id :=c.Param("id")

	iid,err := strconv.ParseInt(id,10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,"id格式有误!")
		return
	}
	err = service.FavoritesDelete(iid)
	if err!=nil{
		util.ResponseError400(c.Writer,"删除收藏失败!")
		return
	}
	util.ResponseSuccess(c.Writer)
}

func favoritesToA(model *dao.Favorites) *Favorites  {

	a :=&Favorites{}
	a.Type = model.Type
	a.Remark = model.Remark
	a.AppId = model.AppId
	a.CoverImg = model.CoverImg
	a.Flag = model.Flag
	a.Json = model.Json
	a.ObjId = model.ObjId
	a.OpenId = model.OpenId

	return a

}

func favoritesToS(param *Favorites) *service.Favorites  {

	s :=&service.Favorites{}
	s.CoverImg = param.CoverImg
	s.OpenId = param.OpenId
	s.AppId = param.AppId
	s.Flag = param.Flag
	s.Json = param.Json
	s.ObjId = param.ObjId
	s.Remark = param.Remark
	s.Type = param.Type

	return s
}