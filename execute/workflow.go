package execute

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"
	"sync"
)
func ExecUpdate(label string, providers []models.Provider, host models.Host,specific bool,lock *sync.Mutex) error {
	// 特定更新的话则上锁
	if specific{
		for {
			if lock.TryLock(){
				break
			}
		}
		defer lock.Unlock()
	}
	// 获取项目目录以备备份配置文件使用
	projectDir, err := utils.GetValue("project-dir")
	if err != nil {
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return err
	}

	// 确认待更新的标签存在于代理配置中
	labelExist := false
	for _, proxy := range providers {
		if proxy.Name == label {
			labelExist = true
			break
		}
	}
	if !labelExist {
		return fmt.Errorf("标签'%s'不存在目前配置中", label)
	}

	// 对标签进行MD5加密,准备新配置文件名
	newFile, err := utils.EncryptionMd5(label)
	if err != nil {
		utils.LoggerCaller("MD5加密失败", err, 1)
		return err
	}

	// 解析服务器URL,构建备份文件名
	link, err := url.Parse(host.Url)
	if err != nil {
		utils.LoggerCaller("主机url解析失败", err, 1)
		return err
	}
	backupFile := link.Hostname()

	// 准备各文件路径
	originalPath := "/opt/singbox/config.json"
	backupPath := filepath.Join(projectDir.(string), "backup", backupFile+".json")
	newPath := filepath.Join(projectDir.(string), "static", "Default", newFile+".json")

	// 更新配置文件并创建备份
	if err := UpdateFile(originalPath, newPath, backupPath, host); err != nil {
		return err
	}

	// 尝试重载配置并验证
	if result, err := ReloadConfig("sing-box", host); err != nil || !result {

		// 配置重载失败时恢复备份配置
		if recoverErr := RecoverFile(originalPath, backupPath, host); recoverErr != nil {
			return recoverErr
		}

		// 尝试启动服务确保服务状态正常
		if startErr := BootService("sing-box", host); startErr != nil {
			return startErr
		}

		return fmt.Errorf("reload new config failed")
	}

	// 成功后更新数据库中的服务器配置标签
	if err := utils.DiskDb.Model(&host).Where("url = ?", host.Url).Update("config", label).Error; err != nil {
		utils.LoggerCaller("更新数据库失败", err, 1)
		return err
	}
	// 完成后日志记录信息
	utils.LoggerCaller(fmt.Sprintf("更新'%s'成功,当前配置为: %s",host.Url,host.Config), nil, 1)
	return nil
}

func GroupUpdate(hosts []models.Host, providers []models.Provider, lock *sync.Mutex) []error{
    // 持续尝试获取锁,以确保并发安全
	for {
		if lock.TryLock() {
			break
		}
	}
    // 确保在函数退出前释放锁
	defer lock.Unlock()
	
    // 使用 WaitGroup 来等待所有更新操作完成
	var hostsWorkflow sync.WaitGroup
	hostsWorkflow.Add(len(hosts) + 1)
	// 创建错误通道和数组用于接收可能的错误
	errChan := make(chan error,3)
	var errList []error
	// 创建计数通道用于接收更新操作的计数
	countChan := make(chan int,3)
    // 遍历服务器列表,对每台服务器启动一个并发更新任务
	for _, host := range hosts {
        // 使用 goroutine 并发执行更新操作
		go func() {
            // 确保在子goroutine退出前通知 WaitGroup
			defer func ()  {
				hostsWorkflow.Done()
				countChan <- 1	
			}()
            // 尝试执行更新操作,并处理可能的错误
			if err := ExecUpdate(host.Config, providers, host, false, lock); err != nil {
				utils.LoggerCaller("update servers config failed", err, 1)
				errChan <- fmt.Errorf("主机'%s'配置'%s'更新失败",host.Url,host.Config)
			}
		}()
	}
	// 创建协程等待所有子协程完成
	go func ()  {
		
		defer func ()  {
			hostsWorkflow.Done()
			close(countChan)
			close(errChan)
		}()
		// 阻塞等待所有子协程完成
		// sum用于计数完成的协程数量
		sum := 0
		// 完成的协程会往计数通道发1,累加即可判断完成的协程数量
		for count := range countChan {
			sum += count
			if sum == len(hosts) {
				return
			}
		}
	}()
	// 获取错误通道的错误信息
	for err := range errChan {
		errList = append(errList, err)
	}

    // 等待所有更新操作完成
	hostsWorkflow.Wait()
	return errList
}