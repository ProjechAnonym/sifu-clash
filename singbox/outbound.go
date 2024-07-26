package singbox

import (
	"sifu-clash/models"
)

func MergeOutbound(provider models.Provider,template string) ([]map[string]interface{},error){
	proxies,err := FetchProxies(provider.Path,provider.Name,template)
	if err != nil {
		return nil,nil
	}
	return proxies,nil
}