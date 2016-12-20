package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	//"gitlab.qiyunxin.com/tangtao/utils/page"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	"strconv"
	"shopapi/service"
	"shopapi/dao"
	"net/http"
)
func CartList(c *gin.Context)  {
	//pindex,psize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	open_id := c.Param("open_id")
	cartList,err := service.CartList(open_id)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	if len(cartList)<1 {
		results :=make([]string,0)
		c.JSON(http.StatusOK,results)
	}else{
		c.JSON(http.StatusOK,cartList)
	}
}

func CartAddToList(c *gin.Context)  {	
	var err error
	var cart dao.Cart
	err =c.BindJSON(&cart)	
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	} 
	
	err = service.CartAddToList(cart)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
func CartMinusFromList(c *gin.Context)  {
	var err error
	_,err = security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	
	var cart dao.Cart
	err =c.BindJSON(&cart)	
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	} 
	
	err = service.CartMinusFromList(cart)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
func CartDelFromList(c *gin.Context)  {	
	var err error
	
	_,err = security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	}
	
	open_id := c.Param("open_id")
	id,err	:= strconv.ParseUint(c.Param("id"),10,64)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	err = service.CartDelFromList(open_id,id)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}
func CartUpdateList(c *gin.Context)  {
	var err error
	/* _,err = security.GetAuthUser(c.Request)
	if err!=nil{
		util.ResponseError(c.Writer,http.StatusUnauthorized,"认证失败!")
		return
	} */
	
	var cart dao.Cart
	err =c.BindJSON(&cart)	
	if err!=nil {
		util.ResponseError(c.Writer,http.StatusBadRequest,err.Error())
		return
	} 
	
	err = service.CartUpdateList(cart)
	if err!=nil{
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	
	util.ResponseSuccess(c.Writer)
}

















