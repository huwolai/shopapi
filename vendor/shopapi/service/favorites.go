package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"shopapi/comm"
	"errors"
)

type Favorites struct  {
	ObjId int64
	AppId string
	OpenId string
	CoverImg string
	Remark string
	Type int
	Flag string
	Json string

}

func FavoritesAdd(favorites *Favorites) error  {

	if favorites.Type!=comm.FAVORITES_TYPE_PRODUCT {
		return errors.New("暂不支持此类收藏!")
	}

	product :=dao.NewProduct()
	product,err :=product.ProductWithId(favorites.ObjId,favorites.AppId)
	if err!=nil{
		log.Error(err)
		return err
	}
	if  product==nil{
		return errors.New("没有找到此商品信息!")
	}

	dfavorites :=dao.NewFavorites()
	dfavorites.AppId = favorites.AppId
	dfavorites.ObjId = favorites.ObjId
	dfavorites.Title = product.Title
	dfavorites.Remark = favorites.Remark
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