package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)
func CommandExec(command string,args ...string) ([]string,[]string,error){
    
	cmd := exec.Command(command,args...)
    
    
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil,nil,err
	}
    
	defer stdoutPipe.Close()

    
	errorsPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil,nil,err
	}
    
	defer errorsPipe.Close()

    
	if err := cmd.Start(); err != nil {
		return nil,nil,err
	}

    
	resultsCh := make(chan string)
	errorsCh := make(chan string)
	procErrs := make(chan error)
	var pipe sync.WaitGroup
	pipe.Add(2)

    
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

    
	var results,errors []string
	for result := range resultsCh {
		results = append(results,result)
	}
	for msg := range errorsCh {
		errors = append(errors,msg)
	}

    
	pipe.Wait()
	close(procErrs)
	procErrsTag := false
	for proc_err := range procErrs {
		LoggerCaller("没有EOF结尾标志",proc_err,1)
		procErrsTag = true
	}
    
	if procErrsTag {
		return results,errors,fmt.Errorf("get pipe output failed")
	}

    
	if err = cmd.Wait();err != nil {
		return results,errors,err
	}
    
	return results,errors,nil
}