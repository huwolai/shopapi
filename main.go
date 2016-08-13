package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	err :=config.Init(false)
	util.CheckErr(err)


	err = startup.InitDBData()
	fmt.Println("-----===")
	fmt.Println(err)

	env := os.Getenv("GO_ENV")
	if env=="tests" {
		gin.SetMode(gin.TestMode)
	}else if env== "production" {
		gin.SetMode(gin.ReleaseMode)
	}else if env == "preproduction" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(CORSMiddleware())
}
