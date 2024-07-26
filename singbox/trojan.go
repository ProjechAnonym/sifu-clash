package singbox

import (
	"net/url"
	"sifu-clash/models"
	"strconv"
	"strings"
)

func MarshalTrojan(proxyMap map[string]interface{}) (map[string]interface{}, error) {
	
	skipCertVerify, err := GetMapValue(proxyMap, "skip-cert-verify")
	if err != nil {
		skipCertVerify = false
	}

	sni, err := GetMapValue(proxyMap, "sni")
	if err != nil {
		return nil, err
	}

	trojan := models.Trojan{
		Type:        "trojan",
		Tag:         proxyMap["name"].(string),
		Server:      proxyMap["server"].(string),
		Server_port: proxyMap["port"].(int),
		Password:    proxyMap["password"].(string),
		Tls: &models.Tls{
			Enabled:     true,
			Insecure:    skipCertVerify.(bool),
			Server_name: sni.(string),
		},
	}
	trojanMap, err := Struct2map(trojan, "trojan")
	if err != nil {
		return nil, err
	}
	return trojanMap, nil
}

func Base64Trojan(link string) (map[string]interface{}, error) {
    // 移除链接前缀"trojan://",以便后续处理
    info := strings.TrimPrefix(link, "trojan://")
    // 使用"@"分割链接字符串,获取密码和服务器信息
    parts := strings.Split(info, "@")
    // 从分割后的第一部分获取密码
    password := parts[0]
    // 使用"#"分割服务器信息,获取服务器URL和标签
    urlParts := strings.Split(parts[1], "#")
    // 解析服务器URL,为后续获取端口和参数做准备
    serverUrl, err := url.Parse("trojan://" + urlParts[0])
    if err != nil {       
        return nil, err
    }
    // 解码标签信息,以便正确使用
    tag, err := url.QueryUnescape(urlParts[1])
    if err != nil {
        return nil, err
    }
    // 获取服务器端口
    port, err := strconv.Atoi(serverUrl.Port())
    if err != nil {
        return nil, err
    }
    
    // 从服务器URL中获取参数
    params := serverUrl.Query()
    // 初始化是否跳过证书验证的变量
    var skipCert bool
    // 根据参数"allowInsecure"的值,确定是否跳过证书验证
    if skipCertVerify := params.Get("allowInsecure"); skipCertVerify != "" {
        if skipCertVerify == "1" {
            skipCert = true
        } else {
            skipCert = false
        }
    } else {
        skipCert = true
    }
    // 构建trojan配置结构体
    trojan := models.Trojan{
        Type: "trojan",
        Tag: tag,
        Password: password,
        Server: serverUrl.Hostname(),
        Server_port: port,
        Tls: &models.Tls{
            Enabled: true,
            Insecure: skipCert,
            Server_name: params.Get("sni"),
        },
    }
    // 将trojan配置结构体转换为map格式
    trojanMap, err := Struct2map(trojan, "trojan")
    if err != nil {
        // 日志记录结构体转换失败
        return nil, err
    }
    // 返回转换后的map和nil错误
    return trojanMap, nil
}