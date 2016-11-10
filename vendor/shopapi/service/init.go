package service

import (
	"shopapi/dao"
)
//厨师面试登记表
func MerchantResumesWithAdd( resumes interface{} ) error  {

	resume:=resumes.(dao.MerchantResume)
	
	err:=dao.MerchantResumesSearchByTel(resume.Tel)	
	if err!=nil {
		return err
	}

	return dao.MerchantResumesWithAdd(resumes)
}