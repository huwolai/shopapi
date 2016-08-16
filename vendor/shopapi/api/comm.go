package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

//在请求中获取AppId
func GetQueryParamInRequest(key string,req *http.Request) string  {

	var value string
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		value = values[0]
	}
	if value=="" {
		value = req.Header.Get(key)
	}

	return value

}

//认证校验
func CheckAppAuth(c *gin.Context) (string,error)  {

	appId := GetQueryParamInRequest("app_id",c.Request)

	if appId==""{

		return appId,errors.New("app_id不能为空")
	}

	return appId,nil
}

//用户认证
func CheckUserAuth(c *gin.Context) (string,error)  {
	openId := GetQueryParamInRequest("open_id",c.Request)

	if openId==""{

		return openId,errors.New("open_id不能为空")
	}

	return openId,nil
}