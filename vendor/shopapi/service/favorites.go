package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"errors"
)

type Favorites struct  {
	ObjId int64
	AppId string
	OpenId string
	CoverImg string
	Title string
	Remark string
	Type int
	Flag string
	Json string

}

func FavoritesAdd(favorites *Favorites) error  {

	dfavorites :=dao.NewFavorites()

	exist,err :=dfavorites.WithTypeAndObjId(favorites.ObjId,favorites.Type,favorites.OpenId,favorites.AppId)
	if err!=nil{
		return errors.New("查询收藏信息失败!")
	}
	if exist {
		return errors.New("已收藏,不能再收藏!")
	}

	dfavorites.AppId = favorites.AppId
	dfavorites.ObjId = favorites.ObjId
	dfavorites.Title = favorites.Title
	dfavorites.Remark = favorites.Remark
	dfavorites.Type = favorites.Type
	dfavorites.OpenId = favorites.OpenId
	dfavorites.Flag = favorites.Flag
	dfavorites.Json = favorites.Json
	dfavorites.CoverImg = favorites.CoverImg
	err =dfavorites.Insert()
	if err!=nil{
		log.Error(err)
		return err
	}

	return err
}

func FavoritesGet(openId,appId string) ([]*dao.Favorites,error)  {

	favorties :=dao.NewFavorites()
	return favorties.WithOpenId(openId,appId)
}

func FavoritesDelete(id int64) error {

	return dao.NewFavorites().DeleteWithId(id)
}

func FavoritesIsExist(objId int64,typ int,appId string,openId string) (*dao.Favorites,error)  {

	return dao.NewFavorites().WithTypeAndObjId(objId,typ,openId,appId)
}