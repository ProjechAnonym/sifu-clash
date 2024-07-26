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
	// 遍历提供的键序列
	for i, key := range keys {
		if tempMap[key] != nil {
			// 如果是最后一个键,返回对应的值的副本和nil错误
			if i == len(keys)-1 {
				return clone.Clone(tempMap[key]), nil
			}
			// 如果当前值是映射,更新tempGlobalVars以继续下一层查找
			if subMap, ok := tempMap[key].(map[string]interface{}); ok {
				tempMap = subMap
			} else {
				// 如果当前值不是映射,返回错误
				return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		} else {
			// 如果当前键不存在,返回错误
			return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
		}
	}
	// 如果所有键都成功找到,返回最终值
	return nil, fmt.Errorf("参数不足,缺少键值参数")
}
func Struct2map[P models.Vmess | models.ShadowSocks | models.Trojan](proxy P,class string) (map[string]interface{},error){
	// 将s配置结构体序列化为JSON格式
	proxyBytes, err := json.Marshal(proxy)
	if err != nil{
		utils.LoggerCaller(fmt.Sprintf("json序列化'%s'失败",class),err,1)
		return nil,err
	}
	// 反序列化JSON数据回map[string]interface{}格式,以便于后续处理
	var proxyMap map[string]interface{}
	
	err = json.Unmarshal(proxyBytes, &proxyMap)
	if err != nil {
		utils.LoggerCaller(fmt.Sprintf("'%s'json序列转换为字典失败",class),err,1)
		return nil,err
	}
	return proxyMap,nil
}