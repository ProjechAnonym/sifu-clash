package route

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sifu-clash/controller"
	"sifu-clash/middleware"
	"sifu-clash/utils"

	"github.com/gin-gonic/gin"
)

func SettingFiles(group *gin.RouterGroup) {
	route := group.Group("/files")
	route.GET("/:file", func(ctx *gin.Context) {
		file := ctx.Param("file")
		template := ctx.Query("template")
		token := ctx.Query("token")
		label := ctx.Query("label")
		projectDir, err := utils.GetValue("project-dir")
		if err != nil {
			utils.LoggerCaller("获取工作目录失败", err, 1)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取工作目录失败"})
		}
		if err := controller.VerifyLink(token); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, label))
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.File(filepath.Join(projectDir.(string), "static", template, file))
	})
	route.GET("fetch",middleware.TokenAuth(),func(ctx *gin.Context){
		links, err := controller.FetchLinks()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, links)
	},
	)
}