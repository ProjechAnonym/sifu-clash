package utils

import (
	"io"
	"os"
	"path/filepath"
	"sifu-clash/models"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)
func SftpRead(host models.Host,src string) ([]byte,error){
    
	config,addr,err := InitClient(host)
	if err != nil {
		return nil,err
	}

    
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil,err
	}
	defer client.Close()

    
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil,err
	}
	defer sftpClient.Close()

    
	srcFile,err := sftpClient.OpenFile(src,os.O_RDONLY)
	if err != nil{
		return nil,err
	}
	defer srcFile.Close()

    
	content,err := io.ReadAll(srcFile)
	if err != nil{
		return nil,err
	}
	
	return content,nil
}

func SftpWrite(host models.Host,content []byte, dst string) error{
    
    config,addr,err := InitClient(host)
    if err != nil {
        return err
    }
    
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
		return err
    }
    defer client.Close()
    
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
		return err
    }
    defer sftpClient.Close()
    
    if _,err := sftpClient.Stat(filepath.Dir(dst)); err != nil{
        if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
            if err := sftpClient.MkdirAll(filepath.Dir(dst));err != nil {
                return err
            }
        }else{
            return err
        }
    }
    
    file, err := sftpClient.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
    defer func() {
        
        if err := file.Close(); err != nil {
            LoggerCaller("文件无法关闭", err,1)
        }
    }()
    if err != nil {
        return err
    }
    
    _, err = file.Write(content)
    if err != nil {
        return err
    }

    
    return nil
}

func SftpDelete(host models.Host,path string) error{
    
	config,addr,err := InitClient(host)
	if err != nil {
		return err
	}

    
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	defer client.Close()

    
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

    
	if _,err := sftpClient.Stat(path);err != nil{
		if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
			return nil
		}else{
			return err
		}
	}

    
	if err := sftpClient.RemoveAll(path);err != nil{
		return err
	}

    
	return nil
}