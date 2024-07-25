package singbox

import (
	"fmt"
	"net/url"
	"sifu-clash/models"
	"sifu-clash/utils"
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