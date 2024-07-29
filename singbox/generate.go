package singbox

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sifu-clash/models"
	"sifu-clash/utils"
	"sync"

	"github.com/huandu/go-clone"
)
func merge(key,projectDir string,template models.Template,mode bool,providers []models.Provider,serviceMap map[string][]models.Ruleset) []error{
	var jobs sync.WaitGroup
	errChannel := make(chan error,len(providers))
	for _,provider := range providers{
		jobs.Add(1)
		go func(){
			defer jobs.Done()
			tempTemplate := clone.Clone(template).(models.Template)
			tempTemplate.Route.Rule_set = append(tempTemplate.Route.Rule_set,SetRulesets(serviceMap)...)
			tempTemplate.Route.Rules = append(tempTemplate.Route.Rules,SetRules(serviceMap,provider)...)
			tempTemplate.Dns.Rules = append(tempTemplate.Dns.Rules,SetDnsRules(serviceMap)...)
			proxies,err := MergeOutbound(provider,serviceMap,tempTemplate.CustomOutbounds)
			if err != nil {
				utils.LoggerCaller(fmt.Sprintf("模板'%s'与'%s'节点合并失败",key,provider.Name),err,1)
				errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败",key,provider.Name)
				return
			}
			tempTemplate.Outbounds = append(tempTemplate.Outbounds,proxies...)
			json,err := json.MarshalIndent(tempTemplate,"","  ")
			if err != nil {
				utils.LoggerCaller("json序列化失败",err,1)
				errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败",key,provider.Name)
				return
			}
			var label string
			if mode {
				md5label,err := utils.EncryptionMd5(provider.Name)
				if err != nil {
					utils.LoggerCaller("md5加密失败",err,1)
					errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败",key,provider.Name)
					return
				}
				label = md5label
			}else{
				label = provider.Name
			}
			if err := utils.FileWrite(json,filepath.Join(projectDir,"static",key,fmt.Sprintf("%s.json",label)));err != nil {
				utils.LoggerCaller("写入文件失败",err,1)
				errChannel <- fmt.Errorf("模板'%s'与'%s'节点合并失败",key,provider.Name)
				return
			}
			utils.LoggerCaller(fmt.Sprintf("模板'%s'与'%s'节点合并成功",key,provider.Name),nil,1)
		}()
	}
	jobs.Wait()
	close(errChannel)
	var errs []error
	for err := range errChannel{
		errs = append(errs,err)
	}
	return errs
}
func Workflow(specific ...int) []error {
	projectDir,err := utils.GetValue("project-dir")
	if err != nil {
		utils.LoggerCaller("获取项目目录失败",err,1)
		return []error{fmt.Errorf("获取项目目录失败")}
	}
	templates, err := utils.GetValue("templates")
	if err != nil {
		utils.LoggerCaller("获取模板配置失败",err,1)
		return []error{fmt.Errorf("获取模板配置失败")}
	}
	var providers []models.Provider
	if len(specific) == 0 {
		if err := utils.MemoryDb.Find(&providers).Error; err != nil {
			utils.LoggerCaller("获取provider配置失败",err,1)
			return []error{fmt.Errorf("获取provider配置失败")}
		}
	}else{
		if err := utils.MemoryDb.Find(&providers,specific).Error; err != nil {
			utils.LoggerCaller("获取provider配置失败",err,1)
			return []error{fmt.Errorf("获取provider配置失败")}
		}
	}
	fmt.Println(providers)
	if len(providers) == 0 {
		utils.LoggerCaller("provider配置为空",nil,1)
		return []error{fmt.Errorf("provider配置为空")}
	}
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取ruleset配置失败",err,1)
		return []error{fmt.Errorf("获取ruleset配置失败")}
	}
	newProviders,errs := AddClashTag(providers)
	newRulesets := SortRulesets(rulesets)
	var workflow sync.WaitGroup
	errChannel := make(chan error,len(newProviders) * len(templates.(map[string]models.Template)))
	server,err := utils.GetValue("mode")
    if err != nil{
        utils.LoggerCaller("获取服务模式配置失败",err,1)
        return []error{fmt.Errorf("获取服务模式配置失败")}
    }
	mode := server.(models.Server).Mode
	for key,value := range(templates.(map[string]models.Template)){
		workflow.Add(1)
		go func(){
			defer workflow.Done()
			errors := merge(key,projectDir.(string),value,mode,newProviders,newRulesets)
			for _,err := range errors{
				errChannel <- err
			}
		}()
	}
	workflow.Wait()
	close(errChannel)
	for err := range errChannel{
		errs = append(errs,err)
	}
	return errs
}