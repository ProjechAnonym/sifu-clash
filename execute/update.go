package execute

import (
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"
)

func BackupFile(origin, backup string, host models.Host) error {
    if host.Localhost {
        // 本地文件复制
        if err := utils.FileCopy(origin, backup); err != nil {
			utils.LoggerCaller("复制文件失败",err,1)
            return err
        }
        // 本地文件删除
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
        // 将远程文件内容写入本地备份文件
        if err := utils.FileWrite(content, backup); err != nil {
			utils.LoggerCaller("写入本地文件失败",err,1)
            return err
        }
        // 通过SFTP协议在远程主机上删除原始文件
        if err := utils.SftpDelete(host, origin); err != nil {
			utils.LoggerCaller("删除远程文件失败",err,1)
            return err
        }
    }
    // 操作成功,返回nil
    return nil
}

func UpdateFile(originFile, newFile, backupFile string, host models.Host) error{
    
    // 创建备份文件所在目录,确保目录存在并设置权限为0755
    if err := utils.DirCreate(filepath.Dir(backupFile));err != nil{
        utils.LoggerCaller("创建备份目录失败！",err,1)
        return err
    }
  
    // 若在本地服务器上操作,执行文件备份及替换
    if host.Localhost{
        
        // 备份原始配置文件
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
            return err
        }
        
        // 将新配置文件复制到原始文件位置,并设置文件权限
        if err := utils.FileCopy(newFile,originFile); err != nil {
            utils.LoggerCaller("设置新配置文件失败",err,1)
            return err
        }
    }else{
        // 若在远程服务器上操作,进行文件备份及上传
        // 备份原始配置文件
        if err := BackupFile(originFile,backupFile,host); err != nil {
            utils.LoggerCaller("备份原文件失败",err,1)
            return err
        }
        // 读取新配置文件内容,准备上传至远程服务器
        content,err := utils.FileRead(newFile)
        if err != nil {
            utils.LoggerCaller("read new config file failed",err,1)
            return err
        }
        // 使用SFTP协议,将新配置文件内容上传至远程服务器的原始文件位置
        if err := utils.SftpWrite(host,content,originFile);err != nil{
            utils.LoggerCaller("上传新配置文件到远程服务器失败",err,1)
            return err
        }
    }
    
    // 文件更新完毕,返回nil表示成功
    return nil
}

func RecoverFile(origin_file,backup_file string, host models.Host) error{
    // 检查是否为本地主机
    if host.Localhost{
        // 对于本地主机,直接使用文件复制函数恢复原始文件
        if err := utils.FileCopy(backup_file,origin_file); err != nil {
            // 记录复制失败的日志并返回错误
            utils.LoggerCaller("恢复原配置文件失败",err,1)
            return err
        }
    }else{
        // 对于远程主机,先从备份文件读取内容
        content,err := utils.FileRead(backup_file)
        if err != nil {
            // 记录读取失败的日志并返回错误
            utils.LoggerCaller("读取备份文件内容失败",err,1)
            return err
        }
        // 使用SFTP协议将内容写入远程主机的原始文件位置
        if err := utils.SftpWrite(host,content,origin_file); err != nil {
            // 记录写入失败的日志并返回错误
            utils.LoggerCaller("写入远程主机原配置文件内容失败",err,1)
            return err
        }
    }
    // 操作成功,返回nil
    return nil
}