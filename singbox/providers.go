package singbox

import (
	"encoding/base64"
	"fmt"
	"sifu-clash/utils"
	"strings"

	"github.com/gocolly/colly/v2"
	"gopkg.in/yaml.v3"
)

func FetchProxies(url,name,template string) ([]map[string]interface{},error) {
	var proxies []map[string]interface{}
	var err error
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		var results []map[string]interface{}
		content := map[string]interface{}{}
		if err = yaml.Unmarshal(r.Body, &content); err != nil {
			utils.LoggerCaller(fmt.Sprintf("解析'%s'yaml配置文件失败",name), err, 1)
			base64msg, err := base64.StdEncoding.DecodeString(string(r.Body))
			if err != nil {
				utils.LoggerCaller(fmt.Sprintf("'%s'base64解码失败",name), err, 1)
				return
			}
			results, err = ParseUrl(strings.Split(string(base64msg), "\n"), name, template)
			if err != nil {
				utils.LoggerCaller(fmt.Sprintf("生成'%s'配置文件失败",name), err, 1)
			}
		} else {
			results, err = ParseYaml(content, name, template)
			if err != nil {
				utils.LoggerCaller(fmt.Sprintf("生成'%s'配置文件失败",name), err, 1)
			}
		}
		proxies = results
	})

	c.OnError(func(r *colly.Response, e error) {
		utils.LoggerCaller(fmt.Sprintf("连接%s失败", r.Request.URL.Host), e, 1)
		err = e
		request_url := r.Request.URL
		params := request_url.Query()
		for k, v := range params {
			if k == "flag" && v[0] == "clash" {
				params.Del("flag")
				request_url.RawQuery = params.Encode()
				c.Visit(request_url.String())
			}
		}
	})
	c.Visit(url)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}