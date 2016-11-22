package service

import (
	"shopapi/dao"
)

func CartList(openId string) ([]dao.Cart,error) {
	return dao.NewCart().CartList(openId)
}
func CartAddToList(cart dao.Cart) error {
	return dao.NewCart().CartAddToList(cart)
}
func CartDelFromList(openId string,id uint64) error {
	return dao.NewCart().CartDelFromList(openId,id)
}