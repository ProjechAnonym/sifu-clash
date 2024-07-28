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

    route.POST("files",func(ctx *gin.Context) {
        form, err := ctx.MultipartForm()
        if err != nil {
            // 日志记录获取多部分表单失败,并返回错误响应
            utils.LoggerCaller("解析表单失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析表单失败"})
            return
        }
        // 获取上传的文件列表
        files := form.File["files"]
        // 获取项目目录路径
        projectDir, err := utils.GetValue("project-dir")
        if err != nil {
            // 日志记录获取项目目录失败,并返回内部服务器错误响应
            utils.LoggerCaller("获取工作目录失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取工作目录失败"})
            return
        }
        // 初始化配置结构体和URL列表
        providers := make([]models.Provider, len(files))
        // 检查temp目录是否存在,不存在则创建
        if err := utils.DirCreate(filepath.Join(projectDir.(string),"temp")); err != nil{
            utils.LoggerCaller("创建temp文件夹失败",err,1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建temp文件夹失败"})
            return
        }
        // 遍历上传的文件,处理并保存每个文件
        for i, file := range files {
            // 解析文件名,用于生成标签
            nameSlice := strings.Split(file.Filename, ".")
            var label string
            if len(nameSlice) <= 2 {
                label = nameSlice[0]
            } else {
                label = strings.Join(nameSlice[0:len(nameSlice)-2], "")
            }
            // 构建文件保存路径,并初始化URL结构体
            providers[i] = models.Provider{Path: filepath.Join(projectDir.(string), "temp", file.Filename), Proxy: false, Name: label, Remote: false}
            // 保存上传的文件到指定路径
            if err := ctx.SaveUploadedFile(file, filepath.Join(projectDir.(string), "temp", file.Filename)); err != nil {
                // 如果保存文件失败,返回内部服务器错误响应
                ctx.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
                return
            }
        }
        // 将处理后的URL列表赋值给配置结构体
        // 调用控制器方法添加配置,处理业务逻辑
        if err := controller.AddItems(models.Proxy{Providers: providers,Rulesets: []models.Ruleset{}}); err != nil {
            // 日志记录添加失败,并返回错误响应
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "添加代理配置失败"})
            return
        }
        // 如果添加成功,返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}