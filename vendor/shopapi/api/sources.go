package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Source struct {
	//资源ID
	Id string `json:"id"`
	//资源名称
	Name string `json:"name"`
	Description string `json:"description"`
	//资源
	Resource string `json:"resource"`
	//
	Permissions string `json:"permissions"`
}

func SourcesAll(c *gin.Context)  {


	c.JSON(http.StatusOK,[]Source{
		Source{Id:"user",Name:"用户",Resource:"users",Permissions:"create,update,delete,browse"},
		Source{Id:"product",Name:"产品",Resource:"products",Permissions:"create,update,delete,browse"},
		Source{Id:"merchant",Name:"商户",Resource:"merchants",Permissions:"create,update,delete,browse"},
	})
}