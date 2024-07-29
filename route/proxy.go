package route

import (
	"net/http"
	"path/filepath"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/models"
	"sifu-clash/utils"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func SettingProxy(group *gin.RouterGroup,lock *sync.Mutex) {
	route := group.Group("/proxy")
	route.Use(middleware.TokenAuth())
	route.GET("fetch",func(ctx *gin.Context) {
		config, err := controller.FetchItems()
        
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取代理配置失败"})
            return
        }
        
        ctx.JSON(http.StatusOK, config)
	})
	route.POST("add",func(ctx *gin.Context) {
		var proxy models.Proxy
        if err := ctx.ShouldBindJSON(&proxy); err != nil {
            
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "序列化json失败"})
            return
        }
        
        if errs := controller.AddItems(proxy,lock); len(errs) != 0 {
            var errors []string
            for _,addErr := range errs {
                errors = append(errors, addErr.Error())
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errors})
            return
        }
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})
	route.DELETE("delete",func(ctx *gin.Context) {
        
        deleteMap := make(map[string][]int) 
        if err := ctx.ShouldBindJSON(&deleteMap); err != nil {
            
            utils.LoggerCaller("序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "序列化json失败"})
            return
        }
        
        
        if err := controller.DeleteProxy(deleteMap); err != nil {
            
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
	})

    route.POST("files",func(ctx *gin.Context) {
        form, err := ctx.MultipartForm()
        if err != nil {
            
            utils.LoggerCaller("解析表单失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析表单失败"})
            return
        }
        
        files := form.File["files"]
        
        projectDir, err := utils.GetValue("project-dir")
        if err != nil {
            
            utils.LoggerCaller("获取工作目录失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取工作目录失败"})
            return
        }
        
        providers := make([]models.Provider, len(files))
        
        if err := utils.DirCreate(filepath.Join(projectDir.(string),"temp")); err != nil{
            utils.LoggerCaller("创建temp文件夹失败",err,1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建temp文件夹失败"})
            return
        }
        
        for i, file := range files {
            
            nameSlice := strings.Split(file.Filename, ".")
            var label string
            if len(nameSlice) <= 2 {
                label = nameSlice[0]
            } else {
                label = strings.Join(nameSlice[0:len(nameSlice)-2], "")
            }
            
            providers[i] = models.Provider{Path: filepath.Join(projectDir.(string), "temp", file.Filename), Proxy: false, Name: label, Remote: false}
            
            if err := ctx.SaveUploadedFile(file, filepath.Join(projectDir.(string), "temp", file.Filename)); err != nil {
                
                ctx.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
                return
            }
        }
        
        
        if err := controller.AddItems(models.Proxy{Providers: providers,Rulesets: []models.Ruleset{}},lock); err != nil {
            
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "添加代理配置失败"})
            return
        }
        
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}