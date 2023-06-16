package app

import (
	"020_logging/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(common.JsonLogger())

	route.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong\n")
	})

	err := route.Run()
	if err != nil {
		return
	}
}
