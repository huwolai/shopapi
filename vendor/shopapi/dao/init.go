package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
	"math/rand"
	"reflect"
	"errors"
	//"gitlab.qiyunxin.com/tangtao/utils/log" //log.Debug(x)
)

type ProductIds struct  {
	Id int
}

type MerchantResume struct {
	Name 				string `json:"name"`
	Nation				string `json:"nation"`
	Gender 				int64  `json:"gender"`
	Birth 				string `json:"birth"`
	Goodat 				string `json:"goodat"`
	Location 			string `json:"location"`
	Household			string `json:"household"`
	Workyear 			int64  `json:"workyear"`
	Professionaltitle 	string `json:"professionaltitle"`
	Honor 				string `json:"honor"`
	Worktype 			int64  `json:"worktype"`
	Servicearea 		string `json:"servicearea"`
	Workexperience 		string `json:"workexperience"`	
	Address	 			string `json:"address"`
	Qq	 				string `json:"qq"`
	Tel	 				string `json:"tel"`
	Wx	 				string `json:"wx"`
	Email	 			string `json:"email"`
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
//厨师面试登记表
func MerchantResumesWithAdd( resumes interface{} ) error  {
	
	resume:=resumes.(MerchantResume)	
	
	builder :=db.NewSession().InsertInto("merchant_resume")
	
	t := reflect.TypeOf(resume)
	v := reflect.ValueOf(resume)
	for i := 0; i < t.NumField(); i++ {
		//fmt.Printf("\n\n%s=%v\n===========================\n", t.Field(i).Name, v.Field(i))		
		builder = builder.Pair(fmt.Sprintf("%v", t.Field(i).Name),fmt.Sprintf("%v", v.Field(i)))	
	}
	_,err :=builder.Exec() 
		
	return err
}
func MerchantResumesSearchByTel( Tel string ) error  {
	count, _ := db.NewSession().Select("count(id)").From("merchant_resume").Where("tel = ?", Tel).ReturnInt64()
	if count >0 {
		return errors.New("该手机号码已注册！")
	}
	return 	nil
}