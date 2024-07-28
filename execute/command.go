package execute

import (
	"fmt"
	"sifu-clash/models"
	"sifu-clash/utils"
	"strings"
)
func ReloadConfig(service string,host models.Host) (bool,error){
	status := true
	var results,errors []string
	if host.Localhost{
		// 使用systemctl命令重新加载sing-box服务
		_,_,err := utils.CommandExec("systemctl", "reload", service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.LoggerCaller("重载配置失败",err,1)
			return false,err
		}

		// 使用journalctl命令获取sing-box服务最近10条日志
		results,errors,err = utils.CommandExec("journalctl", "-u", service,"-n","1")
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.LoggerCaller("获取日志文件失败",err,1)
			return false,err
		}
	}else{
		_,_,err := utils.CommandSsh(host,"systemctl","reload",service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.LoggerCaller("重载配置失败",err,1)
			return false,err
		}
		results,errors,err = utils.CommandSsh(host,"journalctl","-u",service,"-n","1")
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.LoggerCaller("获取日志文件失败",err,1)
			return false,err
		}
	}
	// 检查日志中是否包含错误信息
	for _,result := range(results){
		if strings.Contains(result,"ERROR"){
			// 如果日志包含错误信息,则记录错误并返回
			utils.LoggerCaller("重载配置失败",fmt.Errorf(result),1)
			status = false
			break
		}
	}

	// 如果命令的标准错误输出不为空,则记录错误并返回
	if len(errors) != 0{
		utils.LoggerCaller("错误",fmt.Errorf("命令出现错误返回"),1)
		return false,fmt.Errorf("命令出现错误返回")
	}

	// 如果没有错误但服务未成功重新加载,则返回相应错误
	if !status{
		return false,fmt.Errorf("重载新配置失败")
	}
	
    // 如果一切正常,则返回nil
    return status,nil
}

func BootService(service string,host models.Host) error{
	// 初始化状态变量和错误变量
	var status bool
	var err error

	// 根据服务器是否为本地主机,选择合适的启动服务方法
	if host.Localhost{
		// 使用systemctl命令在本地主机启动服务
		_,_,err = utils.CommandExec("systemctl", "start", service)
	}else{
		// 通过SSH方式在远程主机启动服务
		_,_,err = utils.CommandSsh(host,"systemctl","start",service)
	}

	// 检查启动服务过程中是否发生错误
	if err != nil {
		// 如果发生错误,则记录错误信息并返回错误
		utils.LoggerCaller("启动服务失败",err,1)
		return err
	}

	// 检查服务是否成功启动
	status,err = CheckService(service,host)
	if err != nil {
		// 如果检查服务状态发生错误,则记录错误信息并返回错误
		utils.LoggerCaller(fmt.Sprintf("%s未运行",service),err,1)
		return err
	}

	// 如果服务状态检查失败,则返回错误信息
	if !status {
		return fmt.Errorf("%s状态为关闭",service)
	}

	// 如果一切顺利,返回nil表示服务启动成功
	return nil
}

func CheckService(service string,host models.Host) (bool,error){
	// 初始化服务运行状态为false
	status := false
	// 初始化用于存储命令执行结果和错误信息的切片
	var results,errors []string
	var err error
	// 根据服务器是否为本地服务器,选择不同的方式检查服务状态
	if host.Localhost{
		// 对于本地服务器,直接使用systemctl命令检查服务状态
		results,errors,err = utils.CommandExec("systemctl", "status", service)
	}else{
		// 对于远程服务器,通过SSH连接执行systemctl命令检查服务状态
		results,errors,err = utils.CommandSsh(host,"systemctl", "status", service)
	}
	// 检查命令执行过程中是否发生错误
	if err != nil {
		// 记录错误并返回
		utils.LoggerCaller("获取服务运行状态失败",err,1)
		return false,err
	}
	// 检查命令执行的错误输出是否为空
	if len(errors) != 0{
		// 记录错误并返回
		utils.LoggerCaller("错误",fmt.Errorf("命令出现错误返回"),1)
		return false,fmt.Errorf("命令出现错误返回")
	}
	// 遍历命令执行的结果,检查服务是否处于运行状态
	for _,result := range(results){
		// 如果结果中包含"active (running)",表示服务正在运行
		if strings.Contains(result,"active (running)"){
			status = true
			// 设置状态为运行并终止循环
			break
		}
	}
	// 返回服务运行状态和错误对象
	return status,nil
}