package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
	"math/rand"  
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

//随机数
func RandInt(min, max int) int {
	if min >= max || min==0 || max==0{
		return max
	}
	return rand.Intn(max-min)+min
}

//商品初始化售出数量
func ProductInitNum() error  {
	x := 0	
	for {		
        count, _ := db.NewSession().Select("count(id)").From("product").Where("sold_num = ?", 0).ReturnInt64()
        if count <1 {
            break
        }else{
			x=RandInt(10,50)
			log.Debug(x)
			db.NewSession().UpdateBySql(fmt.Sprintf("update product set sold_num_init=%d,sold_num=%d where sold_num=0 limit 1",x,x)).Exec()
        }
    }	
	return nil
}