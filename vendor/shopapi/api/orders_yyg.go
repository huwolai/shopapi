package api

import (
	"github.com/gin-gonic/gin"
	"net/http"	
	"gitlab.qiyunxin.com/tangtao/utils/page"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"shopapi/service"
	"shopapi/dao"	
)
//中奖管理
func OrdersYygWin(c *gin.Context)  {
	//appId :=security.GetAppId2(c.Request)
	
	var search dao.OrdersYygSearch
	
	pIndex,pSize :=page.ToPageNumOrDefault(c.Query("page_index"),c.Query("page_size"))
	pords,total,err :=service.OrdersYygWin(search,pIndex,pSize)
	if err!=nil {		
		util.ResponseError400(c.Writer,"查询失败！");
		return
	}	
	
	c.JSON(http.StatusOK,page.NewPage(pIndex,pSize,uint64(total),pords))
}










