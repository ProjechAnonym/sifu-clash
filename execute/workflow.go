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
	
	if specific{
		for {
			if lock.TryLock(){
				break
			}
		}
		defer lock.Unlock()
	}
	
	projectDir, err := utils.GetValue("project-dir")
	if err != nil {
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return err
	}

	
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

	
	newFile, err := utils.EncryptionMd5(label)
	if err != nil {
		utils.LoggerCaller("MD5加密失败", err, 1)
		return err
	}

	
	link, err := url.Parse(host.Url)
	if err != nil {
		utils.LoggerCaller("主机url解析失败", err, 1)
		return err
	}
	backupFile := link.Hostname()

	
	originalPath := "/opt/singbox/config.json"
	backupPath := filepath.Join(projectDir.(string), "backup", backupFile+".json")
	newPath := filepath.Join(projectDir.(string), "static", "default", newFile+".json")

	
	if err := UpdateFile(originalPath, newPath, backupPath, host); err != nil {
		return err
	}

	
	if result, err := ReloadConfig("sing-box", host); err != nil || !result {

		
		if recoverErr := RecoverFile(originalPath, backupPath, host); recoverErr != nil {
			return recoverErr
		}

		
		if startErr := BootService("sing-box", host); startErr != nil {
			return startErr
		}

		return fmt.Errorf("reload new config failed")
	}

	
	if err := utils.DiskDb.Model(&host).Where("url = ?", host.Url).Update("config", label).Error; err != nil {
		utils.LoggerCaller("更新数据库失败", err, 1)
		return err
	}
	
	utils.LoggerCaller(fmt.Sprintf("更新'%s'成功,当前配置为: %s",host.Url,host.Config), nil, 1)
	return nil
}

func GroupUpdate(hosts []models.Host, providers []models.Provider, lock *sync.Mutex,only bool) []error{
    if only {
		for {
			if lock.TryLock(){
				break
			}
		}
	}
	defer lock.Unlock()
	var hostsWorkflow sync.WaitGroup
	hostsWorkflow.Add(len(hosts) + 1)
	
	errChan := make(chan error,3)
	var errList []error
	
	countChan := make(chan int,3)
    
	for _, host := range hosts {
        
		go func() {
            
			defer func ()  {
				hostsWorkflow.Done()
				countChan <- 1	
			}()
            
			if err := ExecUpdate(host.Config, providers, host, false, lock); err != nil {
				utils.LoggerCaller("update servers config failed", err, 1)
				errChan <- fmt.Errorf("主机'%s'配置'%s'更新失败",host.Url,host.Config)
			}
		}()
	}
	
	go func ()  {
		
		defer func ()  {
			hostsWorkflow.Done()
			close(countChan)
			close(errChan)
		}()
		
		
		sum := 0
		
		for count := range countChan {
			sum += count
			if sum == len(hosts) {
				return
			}
		}
	}()
	
	for err := range errChan {
		errList = append(errList, err)
	}

    
	hostsWorkflow.Wait()
	return errList
}