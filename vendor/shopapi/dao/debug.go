package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
)


func Debug(bugs ...string) error {
	s:=""
	for _, bug := range bugs {
		s+=bug
	}	
	_,err :=db.NewSession().UpdateBySql(fmt.Sprintf("INSERT INTO `debug` (bug) VALUES ('%s')",s)).Exec()
	return err
}










