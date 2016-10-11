package service

import (
	"shopapi/dao"
)
//商品初始化售出数量
func ProductInitNum() error  {
	err :=dao.ProductInitNum()
	return err
}