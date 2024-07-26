package singbox

import (
	"encoding/base64"
	"encoding/json"
	"sifu-clash/models"
	"strconv"
	"strings"
)

func MarshalVmess(proxyMap map[string]interface{}) (map[string]interface{},error){
	tlsEnable, err := GetMapValue(proxyMap, "tls")
	if err != nil {
		tlsEnable = false
	}
	var vmess models.Vmess
	if tlsEnable.(bool) {
		skipCertVerify, err := GetMapValue(proxyMap, "skip-cert-verify")
		if err != nil {
			skipCertVerify = true
		}
		vmess = models.Vmess{
			Tag:         proxyMap["name"].(string),
			Type:        "vmess",
			Server:      proxyMap["server"].(string),
			Server_port: proxyMap["port"].(int),
			Uuid:        proxyMap["uuid"].(string),
			Alter_id:    proxyMap["alterId"].(int),
			Security:    proxyMap["cipher"].(string),
			Tls: &models.Tls{
				Enabled:  tlsEnable.(bool),
				Insecure: skipCertVerify.(bool),
				Server_name: proxyMap["servername"].(string),
			},
		}
	} else {
		vmess = models.Vmess{
			Tag:         proxyMap["name"].(string),
			Type:        "vmess",
			Server:      proxyMap["server"].(string),
			Server_port: proxyMap["port"].(int),
			Uuid:        proxyMap["uuid"].(string),
			Alter_id:    proxyMap["alterId"].(int),
			Security:    proxyMap["cipher"].(string),
			Tls: nil,
		}
	}
	switch proxyMap["network"].(string) {
	case "grpc":
		service_name, err := GetMapValue(proxyMap, "grpc-opts", "grpc-service-name")
		if err != nil {
			return nil, err
		}
		transport := models.Grpc{
			Type:                  proxyMap["network"].(string),
			Service_name:          service_name.(string),
			Idle_timeout:          "15s",
			Ping_timeout:          "15s",
			Permit_without_stream: false,
		}
		vmess.Transport = transport
	case "ws":
		transport := models.WebSocket{
			Type:                   proxyMap["network"].(string),
			Path:                   proxyMap["ws-path"].(string),
			Headers:                map[string]string{"Host": proxyMap["ws-headers"].(map[string]interface{})["Host"].(string)},
			Early_data_header_name: "Sec-WebSocket-Protocol",
		}
		vmess.Transport = transport
	}
	vmessMap, err := Struct2map(vmess, "vmess")
	if err != nil {
		return nil, err
	}
	return vmessMap, nil
}
func Base64Vmess(link string) (map[string]interface{},error){
    // 移除链接前缀"vmess://",获取base64编码的配置信息
	info := strings.TrimPrefix(link, "vmess://")
    // 解码base64编码的配置信息
	var decodedInfo []byte
	var err error
	decodedInfo, err = base64.URLEncoding.DecodeString(info)
	if err != nil {
		return nil,err
	}
    // 将解码后的信息反序列化为map格式
	var proxyMap map[string]interface{}
	if err := json.Unmarshal(decodedInfo,&proxyMap);err != nil{
		return nil,err
	}
    // 从map中提取并转换端口号和alter_id
	port,err := strconv.Atoi(proxyMap["port"].(string))
	if err != nil {
		return nil,err
	}
	alterId,err := strconv.Atoi(proxyMap["aid"].(string))
	if err != nil {
		return nil,err
	}
    // 判断tls是否启用,以及是否需要跳过证书验证
	var tlsEnable bool
	if _,err := GetMapValue(proxyMap,"tls"); err != nil {
		tlsEnable = false
	}else{
		tlsEnable = true
	}
	
	skipCert,err := GetMapValue(proxyMap,"skip-cert-verify")
	if err != nil {
        // 记录未找到skip_cert_verify键的日志,设置为跳过证书验证,并返回错误
		skipCert = true
	}
	sni,err := GetMapValue(proxyMap,"sni")
	if err != nil{
        // 记录未找到sni键的日志,设置sni为空字符串,并返回错误
		sni = ""
	}
    // 根据提取的信息,构建tls配置
	tls := models.Tls{
		Enabled: tlsEnable,
		Insecure: skipCert.(bool),
		Server_name: sni.(string),
	}
    // 根据提取的信息,构建vmess配置
	vmess := models.Vmess{
		Type: "vmess",
		Tag: proxyMap["ps"].(string),
		Server: proxyMap["add"].(string),
		Server_port: port,
		Uuid: proxyMap["id"].(string),
		Alter_id: alterId,
		Security: "auto",
		Tls: &tls,
	}

    // 根据网络类型设置传输协议
	switch proxyMap["net"].(string) {
		case "ws":
			transport := models.WebSocket{
				Type: proxyMap["net"].(string),
				Path: proxyMap["path"].(string),
				Headers: map[string]string{"host":proxyMap["host"].(string)},
				Early_data_header_name: "Sec-WebSocket-Protocol",
			}
			vmess.Transport = transport
		}
    // 将vmess配置转换为map格式,并返回
	vmessMap,err := Struct2map(vmess,"vmess")
	if err != nil {
		return nil,err
	}
	return vmessMap,nil
}