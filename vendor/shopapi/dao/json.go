package dao

import (
    "encoding/json"
)
func  JsonToMap( s string) (map[string]interface{},error) {
    var result map[string]interface{}
    if err := json.Unmarshal([]byte(s), &result); err != nil {
        return nil, err
    }
    return result, nil
}



