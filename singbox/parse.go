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
        // 如果格式化成功,则将格式化后的代理信息添加到列表中
        if err == nil {
            proxies = append(proxies, result)
        }
    }
    // 返回格式化后的代理信息列表和nil错误,表示处理成功
    return proxies, nil
	
}

func formatYaml(proxyMap map[string]interface{}) (proxy map[string]interface{},err error) {
	// 使用defer和recover处理函数内部可能出现的panic,确保函数能够安全返回
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			utils.LoggerCaller("Panic occurred in FormatUrl", err, 1)
			proxy = nil
			return
		}
	}()
	// 从proxy_map中提取协议类型和标签信息
	// 获取协议类型
	protocolType := proxyMap["type"]

	// 根据协议类型切换不同的处理逻辑
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
		// 如果协议类型不在支持的范围内,返回错误
		return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
	}
	return proxy, err
}
func ParseUrl(urls []string, name, template string) ([]map[string]interface{}, error) {
	if len(urls) == 0 {
        return nil, fmt.Errorf("没有节点信息")
    }
	var proxies []map[string]interface{}

    // 遍历URL列表,尝试对每个URL进行格式化处理
    for _, url := range urls {
        // 格式化URL并获取结果,如果格式化成功,则将结果添加到代理配置切片中
        result, err := formatUrl(url)
        if err == nil {
            proxies = append(proxies, result)
        }
    }

    // 返回处理后的代理配置切片和nil错误
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
	// 解析链接的协议类型
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
		// 如果协议类型不在支持的范围内,则返回错误
		return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
	}
	return proxy, err
	
}