package controller

import (
	"fmt"
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"

	"gopkg.in/yaml.v3"
)

func FetchItems() (*models.Proxy, error) {
	var providers []models.Provider
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		utils.LoggerCaller("获取机场配置失败", err, 1)
		return nil, fmt.Errorf("获取机场配置失败")
	}
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取规则集配置失败", err, 1)
		return nil, fmt.Errorf("获取规则集配置失败")
	}
	return &models.Proxy{Providers: providers,Rulesets: rulesets}, nil
}

func AddItems(newProxy models.Proxy) error {
    projectDir, err := utils.GetValue("project-dir")
    if err != nil {
       
        utils.LoggerCaller("获取工作目录失败", err, 1)
        return fmt.Errorf("获取工作目录失败")
    }
	
    if len(newProxy.Providers) == 0 && len(newProxy.Rulesets) == 0{
        return fmt.Errorf("没有有效配置")
    }
	var addMsg string
	if len(newProxy.Providers) != 0 {
		if err := utils.MemoryDb.Create(&newProxy.Providers).Error; err != nil {
			utils.LoggerCaller("添加机场配置失败", err, 1)
			addMsg = "添加机场配置失败"
		}
	}
	if len(newProxy.Rulesets) != 0 {
		if err := utils.MemoryDb.Create(&newProxy.Rulesets).Error; err != nil {
			utils.LoggerCaller("添加规则集配置失败", err, 1)
			addMsg = addMsg + ",添加规则集配置失败"
		}
	}
	var newProviders []models.Provider
	var newRulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&newProviders).Error; err != nil {
		utils.LoggerCaller("获取机场配置失败", err, 1)
		return fmt.Errorf("获取机场配置失败")
    }
	if err := utils.MemoryDb.Find(&newRulesets).Error; err != nil {
		utils.LoggerCaller("获取规则集配置失败", err, 1)
		return fmt.Errorf("获取规则集配置失败")
    }
    
    proxyYaml, err := yaml.Marshal(models.Proxy{Providers: newProviders,Rulesets: newRulesets})
    if err != nil {  
        utils.LoggerCaller("序列化yaml文件失败", err, 1)
        return fmt.Errorf("序列化yaml文件失败")
    }

    if err := utils.FileWrite(proxyYaml, filepath.Join(projectDir.(string), "config", "proxy.config.yaml")); err != nil { 
        utils.LoggerCaller("写入proxy配置文件失败", err, 1)
        return fmt.Errorf("写入proxy配置文件失败")
    }

	if addMsg != ""{
		return fmt.Errorf(addMsg)
    }
    
    return nil
}

func DeleteProxy(proxy map[string][]int) error{
    // 获取项目目录路径,用于确定生成文件的路径
    projectDir, err := utils.GetValue("project-dir")
    if err != nil {
        // 记录获取项目目录失败的日志并返回错误
        utils.LoggerCaller("获取工作目录失败", err, 1)
        return fmt.Errorf("获取工作目录失败")
    }
    var deletemsg string
	if len(proxy["providers"]) != 0 {
		var tempProviders []models.Provider
		var deleteProviders []models.Provider
		if err := utils.MemoryDb.Find(&tempProviders,proxy["providers"]).Error; err != nil {
			utils.LoggerCaller("获取待删除机场配置失败", err, 1)
			return fmt.Errorf("获取待删除机场配置失败")
		}
		for _,tempProvider := range(tempProviders){
			md5Label,err := utils.EncryptionMd5(tempProvider.Name)
			if err != nil {
				utils.LoggerCaller("加密md5失败",err,1)
				return fmt.Errorf("加密md5失败")
			}
			templates,err := utils.GetValue("templates")
			if err != nil {
				utils.LoggerCaller("获取模板配置失败", err, 1)
				return fmt.Errorf("获取模板配置失败")
			}
			for key := range(templates.(map[string]models.Template)){
				if err := utils.FileDelete(filepath.Join(projectDir.(string), "static", key, md5Label + ".json")); err != nil {
					utils.LoggerCaller(fmt.Sprintf("删除'%s'目录下的'%s'配置文件失败",key,tempProvider.Name),err,1)
				}
			}
			if !tempProvider.Remote{
				if err := utils.FileDelete(tempProvider.Path); err != nil {
					utils.LoggerCaller("删除yaml文件失败",err,1)
				}else {
					deleteProviders = append(deleteProviders, tempProvider)
				}
			}
        }
		if err := utils.MemoryDb.Delete(&deleteProviders).Error; err != nil {
			utils.LoggerCaller("删除机场配置失败", err, 1)
			deletemsg = deletemsg + err.Error()
		}
    }

	if len(proxy["rulesets"]) != 0 {
		if err := utils.MemoryDb.Delete(&models.Ruleset{},proxy["rulesets"]).Error; err != nil {
			utils.LoggerCaller("删除规则集配置失败", err, 1)
			deletemsg = deletemsg + "," + err.Error()
		}
	}
	var providers []models.Provider
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		utils.LoggerCaller("获取机场配置失败", err, 1)
        return fmt.Errorf("获取机场配置失败")
    }
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取规则集配置失败", err, 1)
        return fmt.Errorf("获取规则集配置失败")
    }
	// 将新的代理配置转换为YAML格式
	proxyYaml, err := yaml.Marshal(models.Proxy{Providers: providers,Rulesets: rulesets})
	if err != nil {
		// 记录转换配置失败的日志并返回错误
		utils.LoggerCaller("", err, 1)
		return fmt.Errorf("marshal Proxy failed")
	}

	// 更新代理配置文件
	if err := utils.FileWrite(proxyYaml, filepath.Join(projectDir.(string), "config", "proxy.config.yaml")); err != nil {
		// 记录写入配置文件失败的日志并返回错误
		utils.LoggerCaller("生成代理配置文件失败", err, 1)
		return fmt.Errorf("生成代理配置文件失败")
	}
	if deletemsg != ""{
		return fmt.Errorf(deletemsg)
    }
	
    if len(proxy["providers"]) != 0 {
        var hosts []models.Host
        if err := utils.DiskDb.Find(&hosts).Error; err != nil {
            // 记录查询服务器列表失败的日志并返回错误
            utils.LoggerCaller("查询主机列表失败", err, 1)
            return fmt.Errorf("查询主机列表失败")
        }
        for _,host := range(hosts){
            changeTag := true
            // 如果url列表为0则直接设置改变标签为真
            if len(providers) == 0{
                changeTag = true
            }else{
                for _,provider := range(providers){
                    if host.Config == provider.Name{
                        changeTag = false
                        break
                    }
                }
            }
            if changeTag{
                if err := utils.DiskDb.Model(&models.Host{}).Where("url = ?",host.Url).Update("config","").Error; err != nil{
					utils.LoggerCaller("更换主机配置失败",err,1)
                    return fmt.Errorf("更换主机配置失败")
                }
            }
        }
    }
    // 删除操作成功，返回nil
    return nil
}