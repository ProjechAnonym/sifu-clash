package singbox

import (
	"fmt"
	"os"
	"sifu-clash/models"
	"sifu-clash/utils"

	"gopkg.in/yaml.v3"
)

func outboundSelect(tags []string,label string) map[string]interface{} {
    selectMap := map[string]interface{}{"type":"selector","interrupt_exist_connections":false,"tag":label} 
	// 将"auto"标签添加到标签切片中,确保自动选择选项可用
    tags = append(tags, "auto")
    // 将更新后的标签切片赋值回select_map中的"outbounds"字段
    selectMap["outbounds"] = tags
    // 返回更新后的select_map以及nil错误
    return selectMap
}
// outbound_auto 根据给定的标签和模板生成自动外发配置
// 参数tags是用于外发的标签列表,template是配置模板的字符串
// 返回一个包含外发配置的map以及可能的错误
func outboundAuto(tags []string) map[string]interface{}{
    autoMap := map[string]interface{}{"type":"urltest","interrupt_exist_connections":false,"tag":"auto"}
    autoMap["outbounds"] = tags
    // 返回更新后的自动外发配置以及nil错误
    return autoMap
}
func MergeOutbound(provider models.Provider,serviceMap map[string][]models.Ruleset,outbounds []map[string]interface{}) ([]map[string]interface{},error){
	var proxies []map[string]interface{}
    var content []byte
    var err error
    if provider.Remote {
        proxies,err = FetchProxies(provider.Path,provider.Name)
        if err != nil {
            return nil,err
        }
    } else {
        content ,err = os.ReadFile(provider.Path)
        if err != nil {
            utils.LoggerCaller("读取yaml失败",err,1)
            return nil,err
        }
        var data map[string]interface{}
        if err = yaml.Unmarshal(content,&data); err != nil {
            utils.LoggerCaller("解析yaml失败",err,1)
            return nil,err
        }
        if proxiesMsg,ok := data["proxies"].([]interface{}); ok {
            proxies, err = ParseYaml(proxiesMsg, provider.Name)
        }else{
            err = fmt.Errorf("proxies字段不存在")
        }
        if err != nil {
            return nil,err
        }
    }
    outbounds = append(outbounds,proxies...)
    tags := make([]string,len(outbounds))
    for i,outbound := range outbounds {
        tags[i] = outbound["tag"].(string)
    }
    proxies = append(proxies, outboundSelect(tags,"select"))
    for key,rulesets := range serviceMap {
        if key == ""{
            for _,ruleset := range rulesets {
                if !ruleset.China {
                    proxies = append(proxies,outboundSelect(tags,fmt.Sprintf("select-%s",ruleset.Tag)))
                }
            }
        }else{
            if !rulesets[0].China{
                proxies = append(proxies,outboundSelect(tags,fmt.Sprintf("select-%s",key)))
            }
        }
    }
    proxies = append(proxies, outboundAuto(tags))
	return proxies,nil
}