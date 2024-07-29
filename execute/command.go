package execute

import (
	"fmt"
	"sifu-clash/models"
	"sifu-clash/utils"
	"strings"
)
func ReloadConfig(service string,host models.Host) (bool,error){
	finalStatus := true
	var results,errors []string
	currentStatus,err := CheckService(service,host)
	if err != nil {
		utils.LoggerCaller(fmt.Sprintf("%s未运行",service),err,1)
		return false,err
	}
	if host.Localhost{
		if currentStatus {
			_,_,err = utils.CommandExec("systemctl", "reload", service)
			if err != nil {
				utils.LoggerCaller("重载配置失败",err,1)
				return false,err
			}
		} else {
			err = BootService(service,host)
			if err != nil {
				utils.LoggerCaller("启动服务失败",err,1)
				return false,err
			}
		}
		results,errors,err = utils.CommandExec("journalctl", "-u", service,"-n","1")
		if err != nil {
			utils.LoggerCaller("获取日志文件失败",err,1)
			return false,err
		}
	}else{
		if currentStatus {
			_,_,err = utils.CommandSsh(host,"systemctl","reload",service)
			if err != nil {
				utils.LoggerCaller("重载配置失败",err,1)
				return false,err
			}
		} else {
			err = BootService(service,host)
			if err != nil {
				utils.LoggerCaller("启动服务失败",err,1)
				return false,err
			}
		}
		results,errors,err = utils.CommandSsh(host,"journalctl","-u",service,"-n","1")
		if err != nil {
			utils.LoggerCaller("获取日志文件失败",err,1)
			return false,err
		}
	}
	for _,result := range(results){
		if strings.Contains(result,"ERROR"){
			utils.LoggerCaller("重载配置失败",fmt.Errorf(result),1)
			finalStatus = false
			break
		}
	}

	if len(errors) != 0{
		utils.LoggerCaller("错误",fmt.Errorf("命令出现错误返回"),1)
		return false,fmt.Errorf("命令出现错误返回")
	}

	if !finalStatus{
		return false,fmt.Errorf("重载新配置失败")
	}
	
    return finalStatus,nil
}

func BootService(service string,host models.Host) error{
	var status bool
	var err error

	if host.Localhost{
		_,_,err = utils.CommandExec("systemctl", "start", service)
	}else{
		_,_,err = utils.CommandSsh(host,"systemctl","start",service)
	}

	if err != nil {
		utils.LoggerCaller("启动服务失败",err,1)
		return err
	}

	status,err = CheckService(service,host)
	if err != nil {
		utils.LoggerCaller(fmt.Sprintf("%s未运行",service),err,1)
		return err
	}

	if !status {
		return fmt.Errorf("%s状态为关闭",service)
	}

	return nil
}

func CheckService(service string,host models.Host) (bool,error){
	status := false

	var results,errors []string
	var err error
	if host.Localhost{
		results,errors,err = utils.CommandExec("systemctl", "status", service)
	}else{
		results,errors,err = utils.CommandSsh(host,"systemctl", "status", service)
	}
	if err != nil {
		utils.LoggerCaller("获取服务运行状态失败",err,1)
		return false,err
	}
	if len(errors) != 0{
		utils.LoggerCaller("错误",fmt.Errorf("命令出现错误返回"),1)
		return false,fmt.Errorf("命令出现错误返回")
	}
	for _,result := range(results){
		if strings.Contains(result,"active (running)"){
			status = true
			break
		}
	}
	return status,nil
}