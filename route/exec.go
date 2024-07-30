package route

import (
	"net/http"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func SettingExec(group *gin.RouterGroup,lock *sync.Mutex,cronTask *cron.Cron,id *cron.EntryID){
 	route := group.Group("/exec")
	route.Use(middleware.TokenAuth())
	route.POST("/update",func(ctx *gin.Context) {
		addr := ctx.PostForm("addr")
        config := ctx.PostForm("config")
        
        // 检查config参数是否为空,如果为空,则返回内部服务器错误和错误信息
        if config == "" {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新配置文件为空"})
            return
        }
        
        // 调用controller层的Update_config方法来尝试更新配置
        // 如果更新失败,则记录错误日志并返回内部服务器错误和错误信息
        if err := controller.UpdateConfig(addr, config,lock); err != nil {
            utils.LoggerCaller("更新配置文件失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        
        // 如果更新成功,则返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
}