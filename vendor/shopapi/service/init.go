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
//厨师面试登记表
func MerchantResumesWithAdd( resumes interface{} ) error  {

	resume:=resumes.(dao.MerchantResume)
	
	err:=dao.MerchantResumesSearchByTel(resume.Tel)	
	if err!=nil {
		return err
	}

	return dao.MerchantResumesWithAdd(resumes)
}