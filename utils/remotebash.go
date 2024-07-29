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
    
	config,addr,err := InitClient(host)
	if err != nil {
		return nil,nil,err
	}
	
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil,nil,err
	}
    
	defer client.Close()
	
	session,err := client.NewSession()
	if err != nil {
		return nil,nil,err
	}
    
	defer session.Close()
    
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil,nil,err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return nil,nil,err
	}

    
	resultsChSsh := make(chan string)
	errorsChSsh := make(chan string)
	procErrsSsh := make(chan error)
    
	if err := session.Run(command + " " + strings.Join(args," ")); err != nil {
		return nil,nil,err
	}
    
	var sshPipe sync.WaitGroup
	sshPipe.Add(2)
    
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
    
	var results,errors []string
	for result := range resultsChSsh {
		results = append(results,result)
	}
	for msg := range errorsChSsh {
		errors = append(errors,msg)
	}
    
	sshPipe.Wait()
    
	close(procErrsSsh)
	
	procErrsTag := false
	for procErr := range procErrsSsh {
        
		LoggerCaller("没有EOF结束标志",procErr,1)
		procErrsTag = true
	}
    
	if procErrsTag {
		return results,errors,fmt.Errorf("获取命令输出失败")
	}

    
	return results,errors,nil
	
}