package service

import (
	"shopapi/dao"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
)

//通过类型查询标记
func FlagsWithTypes(stype []string,status []string,appId string) ([]*dao.Flags,error)  {

	return dao.NewFlags().WithTypes(stype,appId,status)
}

func FlagsGetJsonWithTypes(stype []string,status []string,appId string) (string,error)  {
	flags,err:=dao.NewFlags().WithTypes(stype,appId,status)
	if err!=nil || len(flags)<1{
		return "",err
	}
	return flags[0].Json,nil
}
func FlagsSetJsonWithTypes(types string,json string,appId string) error  {
	return dao.NewFlags().FlagsSetJsonWithTypes(types,json,appId)
}
