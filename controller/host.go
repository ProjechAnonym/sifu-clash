package controller

import (
	"fmt"
	"net"
	"net/url"
	"sifu-clash/utils"
)

func IsLocalhost(input_url string) (bool, error) {
	
	parsedUrl, err := url.Parse(input_url)
	if err != nil {
		
		utils.LoggerCaller("无法解析url", err, 1)
		return false, fmt.Errorf("无法解析url")
	}
	
	host := parsedUrl.Hostname()
	
	if ip := net.ParseIP(host); ip != nil {
		
		if ip.IsLoopback() {
			utils.LoggerCaller("地址类型错误", fmt.Errorf("不允许设置回环地址"), 1)
			return false, fmt.Errorf("不允许设置回环地址")
		}
		
		if ip.To4() != nil {
			
			ips, err := net.InterfaceAddrs()
			if err != nil {
				
				utils.LoggerCaller("获取接口失败", err, 1)
				return false, fmt.Errorf("获取接口失败")
			}
			
			for _, addr := range ips {
				ip_addr, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					
					utils.LoggerCaller("解析地址失败", err, 1)
					return false, fmt.Errorf("解析地址失败")
				}
				
				if ip.Equal(ip_addr) {
					return true, nil
				}
			}
			return false, nil
		}
		
		return false, fmt.Errorf("不支持ipv6")
	}
	
	return false, fmt.Errorf("不支持域名")
}