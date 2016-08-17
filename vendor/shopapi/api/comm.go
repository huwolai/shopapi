package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"io"
	"gitlab.qiyunxin.com/tangtao/utils/qtime"
	"time"
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

//图片上传
func ImageUpload(c *gin.Context)  {

	openId,err := CheckUserAuth(c)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	file, header , err := c.Request.FormFile("upload")
	filename := header.Filename
	log.Debug(filename)

	if err != nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	avatar := c.Query("avatar")
	uploadTime := qtime.ToyyyyMMddHHmm(time.Now())
	filepath :="./config/upload/images/" +uploadTime +"/" +util.GenerUUId()
	if avatar=="1" {
		filepath = "./config/upload/avatar/" +openId
	}

	out, err := os.Create(filepath)
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"path": filepath,
	})


}