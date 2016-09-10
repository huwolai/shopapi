package service

import (
	"shopapi/dao"
)

//通过类型查询标记
func FlagsWithTypes(stype []string,status []string,appId string) ([]*dao.Flags,error)  {

	return dao.NewFlags().WithTypes(stype,appId,status)
}
