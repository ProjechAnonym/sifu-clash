package controller

import (
	"fmt"
	"sifu-clash/execute"
	"sifu-clash/models"
	"sifu-clash/utils"
	"sync"
)

func UpdateConfig(addr, config string,lock *sync.Mutex) error {
    // 从数据库中查询服务器信息
    var host models.Host
    if err := utils.DiskDb.Model(&host).Select("localhost", "url", "username", "password").Where("url = ?", addr).First(&host).Error; err != nil {
        // 日志记录查询服务器信息失败
        utils.LoggerCaller("获取主机失败", err, 1)
        return fmt.Errorf("获取主机失败")
    }
	var providers []models.Provider
    // 获取代理配置中的代理信息
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
			// 日志记录查询代理信息失败
			utils.LoggerCaller("获取代理信息失败", err, 1)
			return fmt.Errorf("获取代理信息失败")
	}
    if err := execute.ExecUpdate(config,providers,host,true,lock);err != nil{
        // 日志记录更新配置失败
        utils.LoggerCaller("更新singbox配置失败", err, 1)
        return fmt.Errorf("更新singbox配置失败")
    }
    return nil
}