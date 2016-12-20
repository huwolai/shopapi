package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shopapi/dao"
)

func Test(c *gin.Context) {
	getui:= dao.NewGetui()
	getui.Conn()
	status,_:=getui.Status("ed105d433efafb78c110729c821748c1")	
	
	c.JSON(http.StatusOK,status)
	
	/* if status=="online" {		
		err:=getui.PushSingle("ed105d433efafb78c110729c821748c1","go语言1","test语言2")
		c.JSON(http.StatusOK,err)
		return
	}
	
	err:=getui.PushApnsSingle("5d3f792a86cea378bf541454bbb61ce03639f2a4e5b05ed3f2133aab7edfb9b6","golang","tt语言t语言t语言est","chefOrder")
	
	c.JSON(http.StatusOK,err) */
}