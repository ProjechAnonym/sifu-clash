package middleware

import (
	"sifu-clash/models"
	"sifu-clash/utils"
	"time"

	"github.com/gin-contrib/cors"
)

func Cors() cors.Config {
	serverConfig, err := utils.GetValue("mode")
	if err != nil{
		utils.LoggerCaller("获取运行模式失败",err,1)
	}
	origins := serverConfig.(models.Server).Cors
	// var allow_origins = make([]string, len(origins))
	// copy(allow_origins, origins)
	coresConfig := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET","DELETE"},
		AllowHeaders:     []string{"Origin", "domain", "scheme", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return coresConfig
}