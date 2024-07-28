package route

import (
	"net/http"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/models"
	"sifu-clash/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

func SettingProxy(group *gin.RouterGroup,lock *sync.Mutex) {
	route := group.Group("/proxy")
	route.Use(middleware.TokenAuth())
	route.GET("fetch",func(ctx *gin.Context) {
		config, err := controller.FetchItems()
        // 如果获取失败,记录错误日志并返回内部服务器错误的响应
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取代理配置失败"})
            return
        }
        // 如果获取成功,返回物品信息
        ctx.JSON(http.StatusOK, config)
	})
	route.POST("add",func(ctx *gin.Context) {
		var proxy models.Proxy
        if err := ctx.ShouldBindJSON(&proxy); err != nil {
            // 日志记录JSON绑定失败,并返回错误响应
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "序列化json失败"})
            return
        }
        // 调用控制器方法添加项,处理业务逻辑
        if err := controller.AddItems(proxy); err != nil {
            // 日志记录添加项失败,并返回错误响应
            ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
            return
        }
        // 如果添加成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
	route.DELETE("delete",func(ctx *gin.Context) {
        
        // 解析请求中的JSON数据,填充delete_config结构体
        deleteMap := make(map[string][]int) 
        if err := ctx.BindJSON(&deleteMap); err != nil {
            // 如果解析JSON数据失败,记录错误并返回内部服务器错误
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "序列化json失败"})
            return
        }
        // 调用物品控制器的Delete_items方法,尝试删除指定的物品
        // 使用互斥锁来保证并发安全
        if err := controller.DeleteProxy(deleteMap); err != nil {
            // 如果删除操作失败,记录错误并返回内部服务器错误
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        // 如果删除成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    
	})

}