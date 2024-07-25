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
	for _,provider := range providers{
		jobs.Add(1)
		go func(){
			defer jobs.Done()
			tempTemplate := clone.Clone(template).(models.Template)
			tempTemplate.Route.Rule_set = append(tempTemplate.Route.Rule_set,SetRulesets(serviceMap)...)
			tempTemplate.Route.Rules = append(tempTemplate.Route.Rules,SetRules(serviceMap,false)...)
			tempTemplate.Dns.Rules = append(tempTemplate.Dns.Rules,SetRules(serviceMap,true)...)
			MergeOutbound(provider,key)
			json,_:= json.MarshalIndent(template,"","  ")
			utils.FileWrite(json,filepath.Join(projectDir,"static",key,fmt.Sprintf("%s.json",provider.Name)))
		}()
	}
	jobs.Wait()
	return nil
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
	if err := utils.MemoryDb.Find(&providers).Error; err != nil {
		utils.LoggerCaller("获取provider配置失败",err,1)
		return []error{fmt.Errorf("获取provider配置失败")}
	}
	if len(providers) == 0 {
		utils.LoggerCaller("provider配置为空",nil,1)
		return []error{fmt.Errorf("provider配置为空")}
	}
	var rulesets []models.Ruleset
	if err := utils.MemoryDb.Find(&rulesets).Error; err != nil {
		utils.LoggerCaller("获取ruleset配置失败",err,1)
		return []error{fmt.Errorf("获取ruleset配置失败")}
	}
	newProviders,errs := AddClashTag(providers,specific...)
	newRulesets := SortRulesets(rulesets)
	var workflow sync.WaitGroup
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
			merge(key,projectDir.(string),value,mode,newProviders,newRulesets)
		}()
	}
	workflow.Wait()
	return errs
}