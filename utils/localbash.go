package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)
func CommandExec(command string,args ...string) ([]string,[]string,error){
    // 创建一个命令实例
	cmd := exec.Command(command,args...)
    
    // 创建一个标准输出的管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil,nil,err
	}
    // 确保在函数返回时关闭标准输出管道
	defer stdoutPipe.Close()

    // 创建一个错误输出的管道
	errorsPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil,nil,err
	}
    // 确保在函数返回时关闭错误输出管道
	defer errorsPipe.Close()

    // 启动命令的执行
	if err := cmd.Start(); err != nil {
		return nil,nil,err
	}

    // 创建通道用于接收命令的标准输出和错误输出
	resultsCh := make(chan string)
	errorsCh := make(chan string)
	procErrs := make(chan error)
	var pipe sync.WaitGroup
	pipe.Add(2)

    // 并发读取命令的标准输出
	go func ()  {
		defer func(){
			pipe.Done()
			close(resultsCh)
		}()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			resultsCh <- line
		}
		if scanner.Err() != nil {
			procErrs <- scanner.Err()
		}
	}()

    // 并发读取命令的错误输出
	go func ()  {
		defer func(){
			pipe.Done()
			close(errorsCh)
		}()
		scanner := bufio.NewScanner(errorsPipe)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			errorsCh <- line
		}
		if scanner.Err() != nil {
			procErrs <- scanner.Err()
		}
	}()

    // 收集命令的标准输出和错误输出
	var results,errors []string
	for result := range resultsCh {
		results = append(results,result)
	}
	for msg := range errorsCh {
		errors = append(errors,msg)
	}

    // 等待读取操作完成
	pipe.Wait()
	close(procErrs)
	procErrsTag := false
	for proc_err := range procErrs {
		LoggerCaller("没有EOF结尾标志",proc_err,1)
		procErrsTag = true
	}
    // 如果存在读取错误,返回错误信息
	if procErrsTag {
		return results,errors,fmt.Errorf("get pipe output failed")
	}

    // 等待命令执行完成,并检查是否有错误发生
	if err = cmd.Wait();err != nil {
		return results,errors,err
	}
    // 命令执行成功,返回输出结果
	return results,errors,nil
}