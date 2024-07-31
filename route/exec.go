package route

import (
	"net/http"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/utils"
	"strconv"
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
        
        if config == "" {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新配置文件为空"})
            return
        }
        
        if err := controller.UpdateConfig(addr, config,lock); err != nil {
            utils.LoggerCaller("更新配置文件失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
    route.GET("refresh",func(ctx *gin.Context) {
        if errs := controller.RefreshItems(lock); errs != nil {
            errors := make([]string, len(errs))
            for i, err := range errs {
                errors[i] = err.Error()
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errors})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
    route.POST("check",func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        status,err := controller.CheckStatus(url,service)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        if status{
            ctx.JSON(http.StatusOK, gin.H{"message": true})
        }else{
            ctx.JSON(http.StatusOK, gin.H{"message": false})
        }
    })
    route.POST("boot",func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        if err := controller.BootService(url,service,lock); err!=nil{
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
    route.POST("interval",func(ctx *gin.Context) {
        span := ctx.PostFormArray("span")
        timeSpan := make([]int,len(span))
        var err error
        for i,num := range(span){
            timeSpan[i],err = strconv.Atoi(num)
            if err != nil{
                ctx.JSON(http.StatusBadRequest,gin.H{
                    "message": "间隔必须是整数",
                })
                return
            }
        }
        if err := controller.SetInterval(timeSpan,cronTask,id,lock); err != nil{
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": false})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}