package service

import (
	"shopapi/dao"
)
//商品初始化售出数量
func ProductInitNum() error  {
	err :=dao.ProductInitNum()
	return err
}
//商品 售出数量 定时增加
func ProductAddNum() error  {
	err :=dao.ProductAddNum()
	return err
}