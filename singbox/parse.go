package singbox

import (
	"fmt"
	"sifu-clash/utils"
	"strings"
)

func ParseYaml(content []interface{}, name string) ([]map[string]interface{}, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("没有节点信息")
	}
	var proxies []map[string]interface{}
	for _, proxy := range content {
        result, err := formatYaml(proxy.(map[string]interface{}))

        if err == nil {
            proxies = append(proxies, result)
        }
    }

    return proxies, nil
	
}

func formatYaml(proxyMap map[string]interface{}) (proxy map[string]interface{},err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			utils.LoggerCaller("Panic occurred in FormatUrl", err, 1)
			proxy = nil
			return
		}
	}()

	protocolType := proxyMap["type"]


	switch protocolType {
	case "vmess":
		vmess,err := MarshalVmess(proxyMap)
		if err != nil{
			utils.LoggerCaller("解析vmess失败",err,1)
			return nil,err
		}
		proxy = vmess
	case "ss":
		ss,err := MarshalShadowsocks(proxyMap)
		if err != nil{
			utils.LoggerCaller("解析shadowsocks失败",err,1)
			return nil,err
		}
		proxy = ss
	case "trojan":
		trojan,err := MarshalTrojan(proxyMap)
		if err != nil{
			utils.LoggerCaller("解析trojan失败",err,1)
			return nil,err
		}
		proxy = trojan
	default:
		utils.LoggerCaller("协议未预置",fmt.Errorf("没有预置'%s'协议", protocolType),1)
		return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
	}
	return proxy, err
}
func ParseUrl(urls []string, name string) ([]map[string]interface{}, error) {
	if len(urls) == 0 {
        return nil, fmt.Errorf("没有节点信息")
    }
	var proxies []map[string]interface{}

 
    for _, url := range urls {
      
        result, err := formatUrl(url)
        if err == nil {
            proxies = append(proxies, result)
        }
    }
    return proxies, nil
}
func formatUrl(url string)(proxy map[string]interface{},err error){
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			utils.LoggerCaller("Panic occurred in FormatUrl", err, 1)
			proxy = nil
			return
		}
	}()
	protocolType := strings.Split(url, "://")[0]
	switch protocolType {
	case "ss":
		ss,err := Base64Shadowsocks(url)
		if err != nil{
			utils.LoggerCaller("解析shadowsocks失败",err,1)
			return nil,err
		}
		proxy = ss
	case "vmess":
		vmess,err := Base64Vmess(url)
		if err != nil{
			utils.LoggerCaller("解析vmess失败",err,1)
			return nil,err
		}
		proxy = vmess
	case "trojan":
		trojan,err := Base64Trojan(url)
		if err != nil {
			utils.LoggerCaller("解析trojan失败",err,1)
			return nil,err
		}
		proxy = trojan
	default:
		utils.LoggerCaller("协议未预置",fmt.Errorf("没有预置'%s'协议", protocolType),1)
		return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
	}
	return proxy, err
	
}