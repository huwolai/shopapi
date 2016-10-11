package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
	"math/rand"  
	//"gitlab.qiyunxin.com/tangtao/utils/log" //log.Debug(x)
)

type ProductIds struct  {
	Id int
}

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
			db.NewSession().UpdateBySql(fmt.Sprintf("update product set sold_num_init=%d,sold_num=%d where sold_num=0 limit 1",x,x)).Exec()
        }
    }	
	return nil
}
//商品 售出数量 定时增加
func ProductAddNum() error  {	
	
	var Products []*ProductIds		
	builder :=db.NewSession().Select("id").From("product")	
	//builder = builder.Where("flag = ?","login_type")	
	builder.LoadStructs(&Products)

	x := 0
	for _,Product :=range Products  {
		x=RandInt(1,5)
		db.NewSession().UpdateBySql(fmt.Sprintf("update product set sold_num_init=sold_num_init+%d,sold_num=sold_num+%d where id=%d limit 1",x,x,Product.Id)).Exec()
	}

	return nil
}