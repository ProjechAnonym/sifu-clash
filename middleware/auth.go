package middleware

import (
	"net/http"
	"sifu-clash/models"
	"sifu-clash/utils"

	"github.com/gin-gonic/gin"
)



func TokenAuth() gin.HandlerFunc {
	
	return func(ctx *gin.Context) {
		
		header := ctx.GetHeader("Authorization")
		
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		
		serverConfig, err := utils.GetValue("mode")
		
		if err != nil {
			utils.LoggerCaller("Get key failed!", err, 1)
			ctx.Abort()
			return
		}
		
		key := serverConfig.(models.Server).Token

		
		if key == header {
			
			ctx.Set("token", header)
			return
		}
		
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Abort()
	}
}