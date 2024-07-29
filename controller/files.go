package controller

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"
	"strings"
)

func FetchLinks() (map[string][]map[string]string, error) {
	
	project_dir, err := utils.GetValue("project-dir")
	if err != nil {
		
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return nil, fmt.Errorf("获取工作目录失败")
	}

	
	
	staticDir, err := os.Open(filepath.Join(project_dir.(string), "static"))
	if err != nil {
		
		utils.LoggerCaller("打开静态文件目录失败", err, 1)
		return nil, fmt.Errorf("打开静态文件目录失败")
	}
	defer staticDir.Close()

	
	
	dirs, err := staticDir.ReadDir(-1) 
	if err != nil {
		
		utils.LoggerCaller("无法读取分类目录", err, 1)
		return nil,fmt.Errorf("无法读取分类目录")
	}

	
	var providers []models.Provider
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		
		utils.LoggerCaller("无法获得代理配置", err, 1)
		return nil,fmt.Errorf("无法获得代理配置")
	}

	
	serverConfig, err := utils.GetValue("mode")
	if err != nil {
		
		utils.LoggerCaller("无法读取服务配置", err, 1)
		return nil, fmt.Errorf("无法读取服务配置")
	}

	
	md5Token, err := utils.EncryptionMd5(serverConfig.(models.Server).Key)
	if err != nil {
		
		utils.LoggerCaller("md5加密失败", err, 1)
		return nil,fmt.Errorf("文件托管密钥加密失败")
	}

	
	templateLinks := make(map[string][]map[string]string)
	for _, dir := range dirs {
		
		templateFileDir, err := os.Open(filepath.Join(project_dir.(string), "static", dir.Name()))
		if err != nil {
			
			utils.LoggerCaller(fmt.Sprintf("无法打开'%s'目录",dir.Name()), err, 1)
		}
		defer templateFileDir.Close()

		
		templateFileList, err := templateFileDir.ReadDir(-1)
		if err != nil {
			
			utils.LoggerCaller(fmt.Sprintf("无法读取'%s'目录文件",dir.Name()), err, 1)
		}

		
		var links []map[string]string
		for _, file := range templateFileList {
			
			if file.IsDir() {
				
				utils.LoggerCaller("分类文件夹下存在子文件夹", fmt.Errorf("'%s'是个子文件夹", file.Name()), 1)
				continue
			}
			
			for _, provider := range providers {
				md5Link, err := utils.EncryptionMd5(provider.Name)
				if err != nil {
					
					utils.LoggerCaller("无法加密", err, 1)
				}
				
				if md5Link == strings.Split(file.Name(), ".")[0] {
					
					path, _ := url.JoinPath("api", "files", file.Name())
					params := url.Values{}
					params.Add("token", md5Token)
					params.Add("template", dir.Name())
					params.Add("label", provider.Name)
					path += "?" + params.Encode()
					
					links = append(links, map[string]string{"label": provider.Name, "path": path})
					break
				}
			}
			
			templateLinks[dir.Name()] = links
		}
	}

	return templateLinks, nil
}

func VerifyLink(token string) error {
	
	serverConfig, err := utils.GetValue("mode")
	if err != nil {
		utils.LoggerCaller("获取运行配置失败", err, 1)
		return fmt.Errorf("获取运行配置失败")
	}
	md5Token, err := utils.EncryptionMd5(serverConfig.(models.Server).Key)
	if err != nil {
		utils.LoggerCaller("加密预置密钥失败", err, 1)
		return fmt.Errorf("加密预置密钥失败")
	}
	if token == md5Token {
		return nil
	} else {
		return errors.New("密钥错误")
	}

}