package route

import (
	"net/http"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/models"
	"sifu-clash/utils"

	"github.com/gin-gonic/gin"
)

func SettingHost(group *gin.RouterGroup) {
	route := group.Group("/host")
	route.Use(middleware.TokenAuth())
	route.GET("fetch",func(ctx *gin.Context) {
		var hosts []models.Host
        if err := utils.DiskDb.Select("url","config","localhost","secret","port").Find(&hosts).Error; err != nil{
            utils.LoggerCaller("从数据库中获取主机失败",err,1)
            ctx.JSON(http.StatusInternalServerError,gin.H{"message":"连接数据库失败"})
            return
        }
        ctx.JSON(http.StatusOK,hosts)
	})
	route.POST("add",func(ctx *gin.Context) {
		var content models.Host
        
        
        if err := ctx.ShouldBindJSON(&content); err != nil {
            
            utils.LoggerCaller("反序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "反序列化失败"})
            return
        }
		
        isLocalhost,err := controller.IsLocalhost(content.Url) 
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		content.Localhost = isLocalhost
        
        if err := utils.DiskDb.Create(&content).Error; err != nil {
            
            utils.LoggerCaller("写入数据库失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "写入数据库失败"})
            return
        }
        
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
	route.DELETE("delete",func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        
        
        
        if err := utils.DiskDb.Where("url = ?", url).Delete(&models.Host{}).Error; err != nil {
            utils.LoggerCaller("从数据库删除数据失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "无法从数据库删除数据"})
            return
        }
        
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
       })
}