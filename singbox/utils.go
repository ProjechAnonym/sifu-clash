package singbox

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sifu-clash/models"
	"sifu-clash/utils"

	"github.com/huandu/go-clone"
)

func AddClashTag(providers []models.Provider,specific ...int) ([]models.Provider,[]error){
	var errors []error
	for i,provider := range providers {
		if provider.Remote{
			parsedUrl,err := url.Parse(provider.Path)
			if err != nil {
				utils.LoggerCaller(fmt.Sprintf("解析'%s'链接失败",provider.Name),err,1)
				errors = append(errors,fmt.Errorf("解析'%s'链接失败",provider.Name))
			}
			params := parsedUrl.Query()
			clashTag := false
			for key,values := range params {
				if key == "flag" && values[0] == "clash"{
					clashTag = true
					break
				}
			}
			if !clashTag{
				params.Add("flag","clash")
				parsedUrl.RawQuery = params.Encode()
				providers[i].Path = parsedUrl.String()
			}
		}
	}
	if len(specific) != 0{
		return []models.Provider{providers[specific[0]]},errors
	}
	return providers,errors
}

func SortRulesets(rulesets []models.Ruleset) map[string][]models.Ruleset{
	if len(rulesets) == 0{
		return nil
	}
	serviceMap := make(map[string][]models.Ruleset)
	for _,ruleset := range rulesets{
		if serviceMap[ruleset.Label] == nil{
			serviceMap[ruleset.Label] = []models.Ruleset{ruleset}
		}else{
			serviceMap[ruleset.Label] = append(serviceMap[ruleset.Label], ruleset)
		}
	}
	return serviceMap
}

func GetMapValue(dstMap map[string]interface{},keys ...string)(interface{},error){
	tempMap := dstMap
	
	for i, key := range keys {
		if tempMap[key] != nil {
			
			if i == len(keys)-1 {
				return clone.Clone(tempMap[key]), nil
			}
			
			if subMap, ok := tempMap[key].(map[string]interface{}); ok {
				tempMap = subMap
			} else {
				
				return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		} else {
			
			return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
		}
	}
	
	return nil, fmt.Errorf("参数不足,缺少键值参数")
}
func Struct2map[P models.Vmess | models.ShadowSocks | models.Trojan](proxy P,class string) (map[string]interface{},error){
	
	proxyBytes, err := json.Marshal(proxy)
	if err != nil{
		utils.LoggerCaller(fmt.Sprintf("json序列化'%s'失败",class),err,1)
		return nil,err
	}
	
	var proxyMap map[string]interface{}
	
	err = json.Unmarshal(proxyBytes, &proxyMap)
	if err != nil {
		utils.LoggerCaller(fmt.Sprintf("'%s'json序列转换为字典失败",class),err,1)
		return nil,err
	}
	return proxyMap,nil
}