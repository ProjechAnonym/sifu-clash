package controller

import (
	"fmt"
	"sifu-clash/execute"
	"sifu-clash/models"
	"sifu-clash/singbox"
	"sifu-clash/utils"
	"sync"

	"github.com/robfig/cron/v3"
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

func RefreshItems(lock *sync.Mutex) []error {
    for {
        if lock.TryLock(){
            break
        }
    }
    defer lock.Unlock()
    // 配置工作流
    if errs := singbox.Workflow(); errs != nil {
        // 记录配置工作流失败的日志并返回错误信息
        return errs
    }

    // 从数据库中获取服务器列表
    var hosts []models.Host
    if err := utils.DiskDb.Find(&hosts).Error; err != nil {
        // 记录获取服务器失败的日志并返回错误信息
        utils.LoggerCaller("获取服务器失败", err, 1)
        return []error{fmt.Errorf("获取服务器失败")}
    }
    var providers []models.Provider
    // 获取代理配置中的URL信息
    if err := utils.MemoryDb.Find(&providers).Error; err != nil {
        // 记录获取代理信息失败的日志并返回错误信息
        utils.LoggerCaller("获取代理信息失败", err, 1)
        return []error{fmt.Errorf("获取代理信息失败")}
    }

    // 检查代理配置中的URL列表是否为空
    if len(providers) == 0 {
        // 如果为空,记录错误日志并返回错误信息
        err := fmt.Errorf("配置中没有机场信息")
        utils.LoggerCaller("配置中没有机场信息", err, 1)
        return []error{err}
    }

    // 标记是否需要更新服务器配置
    serverUpdate := false

    // 遍历服务器列表,检查是否需要更新配置
    for _, host := range hosts {
        // 遍历代理配置中的URL,查找匹配的服务器配置
        for _, provider := range providers {
            // 如果找到匹配的配置,则标记需要更新并跳出循环
            if host.Config == provider.Name {
                serverUpdate = true
                break
            }
        }
        // 如果没有找到匹配的配置,将服务器配置更新为第一个URL的标签
        if !serverUpdate {
            host.Config = providers[0].Name
            if err := utils.DiskDb.Model(&models.Host{}).Where("url = ?",host.Url).Update("config",providers[0].Name).Error;err != nil{
                return []error{err}
            }
        }
    }

    // 执行服务器配置更新
    if len(hosts) != 0{
        if errs := execute.GroupUpdate(hosts, providers, lock,false);errs != nil{
            return errs
        }
    }
    // 更新完成,返回nil表示无错误
    return nil
}
func CheckStatus(addr, service string) (bool, error) {
    var host models.Host
    if err := utils.DiskDb.Model(&host).Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("获取主机失败", err, 1)
        return false, fmt.Errorf("获取主机失败")
    }
    status, err := execute.CheckService(service, host)
    if err != nil {
        utils.LoggerCaller("检查服务进程失败", err, 1)
        return false, fmt.Errorf("检查服务进程失败")
    }
    return status, nil
}
func BootService(addr, service string, lock *sync.Mutex) error {
    var host models.Host
    if err := utils.DiskDb.Model(&host).Where("url = ?", addr).First(&host).Error; err != nil {
        utils.LoggerCaller("获取主机失败", err, 1)
        return fmt.Errorf("获取主机失败")
    }

    for {
        if lock.TryLock() {
            break
        }
    }
    defer lock.Unlock()

    if err := execute.BootService(service, host); err != nil {
        utils.LoggerCaller("启动服务失败", err, 1)
        return fmt.Errorf("启动服务失败")
    }

    return nil
}

func SetInterval(span []int, cronTask *cron.Cron, id *cron.EntryID, lock *sync.Mutex) error {
    var newTime string
    switch len(span) {
    case 0:
        newTime = ""
    case 1:
        newTime = fmt.Sprintf("*/%d * * * *",span[0])
    case 2:
        newTime = fmt.Sprintf("%d %d * * *",span[0],span[1])
    case 3:
        newTime = fmt.Sprintf("%d %d * * %d",span[0],span[1],span[2])
    }
    cronTask.Remove(*id)
    var err error
    if newTime != "" {
        *id,err = cronTask.AddFunc(newTime, func() {
            for {
                if lock.TryLock() {
                    break
                }
            }
            defer lock.Unlock()
            singbox.Workflow()
            var hosts []models.Host
            var providers []models.Provider
            if err := utils.DiskDb.Find(&hosts).Error; err != nil {
                utils.LoggerCaller("获取主机列表失败", err, 1)
                return
            }
            if err := utils.MemoryDb.Find(&providers).Error; err != nil {
                utils.LoggerCaller("获取代理信息失败", err, 1)
                return
            }
            execute.GroupUpdate(hosts, providers, lock, false)
        })
        if err != nil{
            utils.LoggerCaller("修改定时任务失败", err, 1)
            return fmt.Errorf("修改定时任务失败")
        }
    }
    return nil
}