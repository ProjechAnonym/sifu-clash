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
        
        // 尝试将请求体中的JSON数据绑定到content变量上
        if err := ctx.BindJSON(&content); err != nil {
            // 如果绑定失败,记录错误日志并返回内部服务器错误
            utils.LoggerCaller("反序列化json失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "反序列化失败"})
            return
        }
		// 判断IP是否指向本机
        isLocalhost,err := controller.IsLocalhost(content.Url) 
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		content.Localhost = isLocalhost
        // 尝试将content变量中的服务器信息插入到数据库中
        if err := utils.DiskDb.Create(&content).Error; err != nil {
            // 如果插入失败,记录错误日志并返回内部服务器错误
            utils.LoggerCaller("写入数据库失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "写入数据库失败"})
            return
        }
        
        // 如果插入成功,返回状态码200和成功信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
	route.DELETE("delete",func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        
        // 使用GORM从数据库中删除URL对应的服务器记录
        // 如果删除操作失败,记录错误并返回内部服务器错误响应
        if err := utils.DiskDb.Where("url = ?", url).Delete(&models.Host{}).Error; err != nil {
            utils.LoggerCaller("从数据库删除数据失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "无法从数据库删除数据"})
            return
        }
        
        // 如果删除成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
       })
}