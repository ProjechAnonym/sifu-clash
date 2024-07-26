package singbox

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"sifu-clash/models"
	"strconv"
	"strings"
)

func MarshalShadowsocks(proxyMap map[string]interface{}) (map[string]interface{}, error) {
	// 创建一个shadowsocks配置结构体实例,初始化其字段值从proxy_map中获取
	ss := models.ShadowSocks{
		Type:        "shadowsocks",
		Tag:         proxyMap["name"].(string),
		Server:      proxyMap["server"].(string),
		Server_port: proxyMap["port"].(int),
		Method:      proxyMap["cipher"].(string),
		Password:    proxyMap["password"].(string),
	}

	// 将shadowsocks配置结构体转换为map,便于后续处理或返回
	// 这里使用了Struct2map函数进行转换,如果转换失败,则记录错误日志并返回错误
	ssMap, err := Struct2map(ss, "ss")
	if err != nil {
		return nil, err
	}

	// 转换成功,返回转换后的map以及nil错误
	return ssMap, nil
}
func Base64Shadowsocks(link string) (map[string]interface{}, error) {
    // 移除链接前缀"ss://"并解码URL编码的部分
	info, err := url.QueryUnescape(strings.TrimPrefix(link, "ss://"))
	if err != nil {
		return nil, err
	}

    // 根据"@"分割解码后的信息和服务器信息
	parts := strings.Split(info, "@")
    // 解码信息部分
	decodedInfo, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

    // 分割解码后的信息为方法和密码
	infoParts := strings.Split(string(decodedInfo), ":")
    // 分割服务器信息为地址和标签
	serverInfo := strings.Split(parts[1], "#")

    // 解析服务器URL
	serverUrl, err := url.Parse("ss://" + serverInfo[0])
	if err != nil {
        // 返回解析服务器URL失败的错误
		return nil, fmt.Errorf("failed to parse server URL: %v", err)
	}

    // 获取服务器端口
	port, err := strconv.Atoi(serverUrl.Port())
	if err != nil {
		return nil, err
	}

    // 构建shadowsocks配置结构体
	ss := models.ShadowSocks{
		Type: "shadowsocks",
		Tag: serverInfo[1],
		Server: serverUrl.Hostname(),
		Server_port: port,
		Method: infoParts[0],
		Password: infoParts[1],
	}

    // 将shadowsocks结构体转换为map
	ssMap, err := Struct2map(ss, "shadowsocks")
	if err != nil {
        // 记录日志并返回错误
		return nil, err
	}

    // 返回转换后的map和nil错误
	return ssMap, nil
}