package dao

import (
    "encoding/json"
    "errors"
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



