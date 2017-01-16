package dao

import (
    "encoding/json"
    "errors"
	"strings"
	"strconv"
	"fmt"
)
func  JsonToMap( s string) (map[string]interface{},error) {
	if len(s)<1 {
		return nil,errors.New("字符串不能为空")
	}

    var result map[string]interface{}
    if err := json.Unmarshal([]byte(s), &result); err != nil {
        return nil, err
    }
    return result, nil
}
func Right(str string, length int) string {
    rs := []rune(str)
    rl := len(rs)
   
	start:=rl-length
   
    return string(rs[start:])
}
func Left(str string, length int) string {
    rs := []rune(str)      
    return string(rs[0:length])
}
func Float2int(f float64) int64 {
    i64,_ := strconv.ParseInt(strings.Replace(fmt.Sprintf("%.2f",f), ".", "", -1), 10, 64)
	return i64
}



