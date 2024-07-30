package execute

import (
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"
)

func BackupFile(origin, backup string, host models.Host) error {
    if host.Localhost {
        if err := utils.FileCopy(origin, backup); err != nil {
			utils.LoggerCaller("复制文件失败",err,1)
            return err
        }
        if err := utils.FileDelete(origin); err != nil {
			utils.LoggerCaller("删除文件失败",err,1)
            return err
        }
    } else {
        content, err := utils.SftpRead(host, origin)
        if err != nil {
			utils.LoggerCaller("读取远程文件失败",err,1)
            return err
        }
        if err := utils.FileWrite(content, backup); err != nil {
			utils.LoggerCaller("写入本地文件失败",err,1)
            return err
        }
        if err := utils.SftpDelete(host, origin); err != nil {
			utils.LoggerCaller("删除远程文件失败",err,1)
            return err
        }
    }
    return nil
}

func UpdateFile(originFile, newFile, backupFile string, host models.Host) error{
    if err := utils.DirCreate(filepath.Dir(backupFile));err != nil{
        utils.LoggerCaller("创建备份目录失败！",err,1)
        return err
    }
    if host.Localhost{
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
        }
        if err := utils.FileCopy(newFile,originFile); err != nil {
            utils.LoggerCaller("设置新配置文件失败",err,1)
            return err
        }
    }else{
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
            return err
        }
        content,err := utils.FileRead(newFile)
        if err != nil {
            utils.LoggerCaller("read new config file failed",err,1)
            return err
        }
        if err := utils.SftpWrite(host,content,originFile);err != nil{
            utils.LoggerCaller("上传新配置文件到远程服务器失败",err,1)
            return err
        }
    }
    return nil
}

func RecoverFile(origin_file,backup_file string, host models.Host) error{
    if host.Localhost{
        if err := utils.FileCopy(backup_file,origin_file); err != nil {
            utils.LoggerCaller("恢复原配置文件失败",err,1)
            return err
        }
    }else{
        content,err := utils.FileRead(backup_file)
        if err != nil {
            utils.LoggerCaller("读取备份文件内容失败",err,1)
            return err
        }
        
        if err := utils.SftpWrite(host,content,origin_file); err != nil {
            
            utils.LoggerCaller("写入远程主机原配置文件内容失败",err,1)
            return err
        }
    }
    
    return nil
}