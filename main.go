package main

import (
	"os"
	"path/filepath"
	"sifu-clash/middleware"
	"sifu-clash/models"
	"sifu-clash/route"
	"sifu-clash/singbox"
	"sifu-clash/utils"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	utils.SetValue(utils.GetProjectDir(),"project-dir")
	utils.GetCore()
	utils.GetDatabase()
	if err := utils.LoadConfig(filepath.Join("config","mode.config.yaml"),"mode"); err != nil {
		utils.LoggerCaller("加载服务模式配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载服务模式配置完成",nil,1)
	if err := utils.LoadConfig(filepath.Join("config","proxy.config.yaml"),"proxy"); err != nil {
		utils.LoggerCaller("加载代理集合配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载代理集合配置完成",nil,1)
	if err := utils.LoadTemplate(); err != nil {
		utils.LoggerCaller("加载模板配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载模板配置完成",nil,1)
	utils.LoggerCaller("服务启动成功",nil,1)
}
func main() {
	serverMode,err := utils.GetValue("mode")
	if err != nil {
		utils.LoggerCaller("获取服务模式失败",err,1)
		os.Exit(2)
	}
	if serverMode.(models.Server).Mode {
		var lock sync.Mutex
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(),middleware.Recovery(true),cors.New(middleware.Cors()))
		apiGroup := server.Group("/api")
		apiGroup.GET("verify",middleware.TokenAuth())
		route.SettingHost(apiGroup)
		route.SettingFiles(apiGroup)
		route.SettingProxy(apiGroup,&lock)
		server.Run(serverMode.(models.Server).Listen)
	}else{
		singbox.Workflow()
	}
}