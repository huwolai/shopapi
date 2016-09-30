package comm

import "math"

//数组去重
func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

//
func Floor(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.0/pow10_n)*pow10_n) / pow10_n
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}


//数组对象转换（一个数组里的对象都转换为另一个对象的数组）
func ToArray(sourceArray []interface{},convertFunc func(source interface{})(dest interface{})) ([]interface{})  {
	dests :=make([]interface{},0)
	if sourceArray!=nil{
		for _,sr :=range sourceArray  {
			dt := convertFunc(sr)
			dests = append(dests,dt)
		}
	}

	return dests
}