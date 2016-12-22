package api

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"fmt"
	"time"
)

func Test(c *gin.Context) {
	cr := cron.New()
	cr.AddFunc("0/5 * * * * ?", func() { 
		fmt.Println(time.Now().Format("05 04 15 * * ?"))
		//cr.Stop()
	})
	cr.Start()
}
func Test1(c *gin.Context) {
	cr := cron.New()
	cr.AddFunc("0/3 * * * * ?", func() { 
		fmt.Println(time.Now().Format("05 04 15 * * ?"))
		cr.Stop()
	})
	cr.Start()
}

















