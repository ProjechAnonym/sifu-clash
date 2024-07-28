package utils

import (
	"bufio"
	"fmt"
	"sifu-clash/models"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)
func CommandSsh(host models.Host,command string,args ...string) ([]string,[]string,error){
    // 初始化SSH客户端配置并获取服务器地址
	config,addr,err := InitClient(host)
	if err != nil {
		return nil,nil,err
	}
	// 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil,nil,err
	}
    // 确保SSH连接在函数返回前关闭
	defer client.Close()
	// 创建一个新的SSH会话
	session,err := client.NewSession()
	if err != nil {
		return nil,nil,err
	}
    // 确保SSH会话在函数返回前关闭
	defer session.Close()
    // 设置标准输出和标准错误的管道
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil,nil,err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return nil,nil,err
	}

    // 创建通道以异步收集命令的标准输出和错误输出
	resultsChSsh := make(chan string)
	errorsChSsh := make(chan string)
	procErrsSsh := make(chan error)
    // 执行命令
	if err := session.Run(command + " " + strings.Join(args," ")); err != nil {
		return nil,nil,err
	}
    // 同步读取标准输出和错误输出的等待组
	var sshPipe sync.WaitGroup
	sshPipe.Add(2)
    // 读取标准输出的协程
	go func ()  {
		defer func(){
			sshPipe.Done()
			close(resultsChSsh)
		}()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			resultsChSsh <- line
		}
		if scanner.Err() != nil {
			procErrsSsh <- scanner.Err()
		}
	}()
    // 读取标准错误的协程
	go func ()  {
		defer func(){
			sshPipe.Done()
			close(errorsChSsh)
		}()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			errorsChSsh <- line
		}
		if scanner.Err() != nil {
			procErrsSsh <- scanner.Err()
		}
	}()
    // 收集标准输出和错误输出的结果
	var results,errors []string
	for result := range resultsChSsh {
		results = append(results,result)
	}
	for msg := range errorsChSsh {
		errors = append(errors,msg)
	}
    // 等待读取标准输出和错误输出的协程完成
	sshPipe.Wait()
    // 关闭错误通道
	close(procErrsSsh)
	// 检查读取过程中是否有错误发生
	procErrsTag := false
	for procErr := range procErrsSsh {
        // 如果读取过程中有错误发生,记录错误
		LoggerCaller("没有EOF结束标志",procErr,1)
		procErrsTag = true
	}
    // 如果存在读取错误,返回错误信息
	if procErrsTag {
		return results,errors,fmt.Errorf("获取命令输出失败")
	}

    // 返回命令执行的成功输出和错误输出
	return results,errors,nil
	
}