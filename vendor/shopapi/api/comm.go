package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"io"
	"time"
)

const MERCHANT_IMG_PATH  = "./config/upload/merchant"
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

	name := c.Query("name")
	if name=="" {
		name = "file"
	}

	if c.Request.MultipartForm == nil {
		err := c.Request.ParseMultipartForm(32 << 20)
		log.Error(err)
	}

	log.Error(c.Request.MultipartForm)

	file, _ , err := c.Request.FormFile(name)

	if err != nil {
		log.Debug("获取文件错误!")
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	rootDir :="./config/upload"

	typeS := c.Query("type")
	uploadTime := time.Now().Format("200601")
	fileDir :="/images" +"/" +uploadTime
	fileName :=util.GenerUUId()
	if typeS=="avatar" {
		fileDir = "/avatar"
		fileName =  openId
	}else if typeS=="merchant" {
		fileDir = "/merchant"
		fileName =  openId
	}else{

	}
	err =os.MkdirAll(rootDir+"/" +fileDir,0777)
	if err!=nil {
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	filepath :=rootDir+"/" +fileDir + "/" + fileName
	out, err := os.Create(filepath)
	if err!=nil{
		log.Debug("创建文件失败",filepath)
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Debug("复制文件错误!")
		log.Error(err)
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"path": fileDir+"/" +fileName,
	})


}