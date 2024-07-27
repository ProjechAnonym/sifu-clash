package controller

import (
	"fmt"
	"net"
	"net/url"
	"sifu-clash/utils"
)

func IsLocalhost(input_url string) (bool, error) {
	// 解析输入的URL
	parsedUrl, err := url.Parse(input_url)
	if err != nil {
		// 日志记录URL解析错误
		utils.LoggerCaller("无法解析url", err, 1)
		return false, fmt.Errorf("无法解析url")
	}
	// 获取URL的主机名
	host := parsedUrl.Hostname()
	// 尝试将主机名解析为IP地址
	if ip := net.ParseIP(host); ip != nil {
		// 检查是否为回环地址,回环地址不被允许
		if ip.IsLoopback() {
			utils.LoggerCaller("地址类型错误", fmt.Errorf("不允许设置回环地址"), 1)
			return false, fmt.Errorf("不允许设置回环地址")
		}
		// 检查是否为IPv4地址
		if ip.To4() != nil {
			// 获取本地网络接口地址
			ips, err := net.InterfaceAddrs()
			if err != nil {
				// 日志记录获取接口地址失败
				utils.LoggerCaller("获取接口失败", err, 1)
				return false, fmt.Errorf("获取接口失败")
			}
			// 遍历所有接口地址,检查是否与输入的IP地址匹配
			for _, addr := range ips {
				ip_addr, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					// 日志记录CIDR解析失败
					utils.LoggerCaller("解析地址失败", err, 1)
					return false, fmt.Errorf("解析地址失败")
				}
				// 匹配本机ip,确认输入的url指向本机
				if ip.Equal(ip_addr) {
					return true, nil
				}
			}
			return false, nil
		}
		// 如果不是IPv4地址,则返回错误,不支持IPv6地址
		return false, fmt.Errorf("不支持ipv6")
	}
	// 如果不是IP地址,则返回错误,域名不被允许。
	return false, fmt.Errorf("不支持域名")
}