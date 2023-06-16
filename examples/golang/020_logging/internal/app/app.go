package app

import (
	"020_logging/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	// [<snippet log_enable>]
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(common.JsonLogger())
	// [<endsnippet log_enable>]

	route.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong\n")
	})

	err := route.Run()
	if err != nil {
		return
	}
}
