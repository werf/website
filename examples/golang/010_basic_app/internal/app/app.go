package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	route := gin.Default()

	route.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong\n")
	})

	err := route.Run()
	if err != nil {
		return
	}
}
