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
	// 获取项目目录路径
	project_dir, err := utils.GetValue("project-dir")
	if err != nil {
		// 记录获取项目目录失败的日志
		utils.LoggerCaller("获取工作目录失败", err, 1)
		return nil, fmt.Errorf("获取工作目录失败")
	}

	// 打开静态文件目录
	// 打开目录
	staticDir, err := os.Open(filepath.Join(project_dir.(string), "static"))
	if err != nil {
		// 记录打开静态文件目录失败的日志
		utils.LoggerCaller("打开静态文件目录失败", err, 1)
		return nil, fmt.Errorf("打开静态文件目录失败")
	}
	defer staticDir.Close()

	// 读取目录中的所有文件和子目录
	// 读取目录条目
	dirs, err := staticDir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		// 记录读取目录失败的日志
		utils.LoggerCaller("无法读取分类目录", err, 1)
		return nil,fmt.Errorf("无法读取分类目录")
	}

	// 获取代理配置
	var providers []models.Provider
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		// 记录获取代理配置失败的日志
		utils.LoggerCaller("无法获得代理配置", err, 1)
		return nil,fmt.Errorf("无法获得代理配置")
	}

	// 获取服务器配置
	serverConfig, err := utils.GetValue("mode")
	if err != nil {
		// 记录获取服务器配置失败的日志
		utils.LoggerCaller("无法读取服务配置", err, 1)
		return nil, fmt.Errorf("无法读取服务配置")
	}

	// 对服务器配置中的令牌进行MD5加密
	md5Token, err := utils.EncryptionMd5(serverConfig.(models.Server).Key)
	if err != nil {
		// 记录加密令牌失败的日志
		utils.LoggerCaller("md5加密失败", err, 1)
		return nil,fmt.Errorf("文件托管密钥加密失败")
	}

	// 初始化存储链接信息的映射
	templateLinks := make(map[string][]map[string]string)
	for _, dir := range dirs {
		// 打开模板文件目录
		templateFileDir, err := os.Open(filepath.Join(project_dir.(string), "static", dir.Name()))
		if err != nil {
			// 记录打开模板文件目录失败的日志
			utils.LoggerCaller(fmt.Sprintf("无法打开'%s'目录",dir.Name()), err, 1)
		}
		defer templateFileDir.Close()

		// 读取模板文件目录中的所有文件和子目录
		templateFileList, err := templateFileDir.ReadDir(-1)
		if err != nil {
			// 记录读取模板文件目录失败的日志
			utils.LoggerCaller(fmt.Sprintf("无法读取'%s'目录文件",dir.Name()), err, 1)
		}

		// 初始化存储当前模板链接的数组
		var links []map[string]string
		for _, file := range templateFileList {
			// 跳过子目录
			if file.IsDir() {
				// 记录模板目录包含子目录的日志
				utils.LoggerCaller("分类文件夹下存在子文件夹", fmt.Errorf("'%s'是个子文件夹", file.Name()), 1)
				continue
			}
			// 遍历代理配置中的链接,匹配文件名
			for _, provider := range providers {
				md5Link, err := utils.EncryptionMd5(provider.Name)
				if err != nil {
					// 记录加密链接标签失败的日志
					utils.LoggerCaller("无法加密", err, 1)
				}
				// 如果文件名的MD5与链接标签的MD5匹配,则处理该链接
				if md5Link == strings.Split(file.Name(), ".")[0] {
					// 构建链接的完整路径
					path, _ := url.JoinPath("api", "files", file.Name())
					params := url.Values{}
					params.Add("token", md5Token)
					params.Add("template", dir.Name())
					params.Add("label", provider.Name)
					path += "?" + params.Encode()
					// 将处理后的链接添加到数组中
					links = append(links, map[string]string{"label": provider.Name, "path": path})
					break
				}
			}
			// 将当前模板的链接数组添加到映射中
			templateLinks[dir.Name()] = links
		}
	}

	return templateLinks, nil
}

func VerifyLink(token string) error {
	// 获取配置文件
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