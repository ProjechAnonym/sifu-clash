package main

import (
	"os"
	"path/filepath"
	"sifu-clash/singbox"
	"sifu-clash/utils"
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
	singbox.Workflow()
}