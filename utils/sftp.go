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
    // 初始化SSH客户端配置
	config,addr,err := InitClient(host)
	if err != nil {
		return nil,err
	}

    // 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil,err
	}
	defer client.Close()

    // 初始化SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil,err
	}
	defer sftpClient.Close()

    // 打开指定路径的文件
	srcFile,err := sftpClient.OpenFile(src,os.O_RDONLY)
	if err != nil{
		return nil,err
	}
	defer srcFile.Close()

    // 读取文件全部内容
	content,err := io.ReadAll(srcFile)
	if err != nil{
		return nil,err
	}
	
	return content,nil
}

func SftpWrite(host models.Host,content []byte, dst string) error{
    // 初始化SSH客户端配置并获取服务器地址
    config,addr,err := InitClient(host)
    if err != nil {
        return err
    }
    // 建立SSH连接
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
		return err
    }
    defer client.Close()
    // 初始化SFTP客户端
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
		return err
    }
    defer sftpClient.Close()
    // 检查目标文件目录是否存在,不存在则创建
    if _,err := sftpClient.Stat(filepath.Dir(dst)); err != nil{
        if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
            if err := sftpClient.MkdirAll(filepath.Dir(dst));err != nil {
                return err
            }
        }else{
            return err
        }
    }
    // 打开目标文件,如果不存在则创建,并设置为可写模式
    file, err := sftpClient.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
    defer func() {
        // 确保文件关闭
        if err := file.Close(); err != nil {
            LoggerCaller("文件无法关闭", err,1)
        }
    }()
    if err != nil {
        return err
    }
    // 写入内容到文件
    _, err = file.Write(content)
    if err != nil {
        return err
    }

    // 操作成功
    return nil
}

func SftpDelete(host models.Host,path string) error{
    // 初始化SSH客户端配置
	config,addr,err := InitClient(host)
	if err != nil {
		return err
	}

    // 建立SSH连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	defer client.Close()

    // 初始化SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

    // 检查文件或目录是否存在
	if _,err := sftpClient.Stat(path);err != nil{
		if err.(*sftp.StatusError).Code == uint32(sftp.ErrSSHFxNoSuchFile){
			return nil
		}else{
			return err
		}
	}

    // 删除文件或目录
	if err := sftpClient.RemoveAll(path);err != nil{
		return err
	}

    // 删除操作成功,返回nil
	return nil
}