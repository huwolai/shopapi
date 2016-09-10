package service

import (
	"shopapi/dao"
	"strings"
)

//通过类型查询标记
func FlagsWithTypes(stype string,appId string) ([]*dao.Flags,error)  {

	return dao.NewFlags().WithTypes(strings.Split(stype,","),appId)
}
