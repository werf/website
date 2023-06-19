package app

import (
	"030_assets/internal/common"
	"030_assets/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	// [<snippet log_enable>]
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(common.JsonLogger())
	// [<endsnippet log_enable>]

	route.Static("/static/stylesheets", "static/stylesheets")
	route.Static("/static/javascripts", "static/javascripts")
	route.Static("/static/images", "static/images")

	route.LoadHTMLGlob("templates/*")

	route.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong\n")
	})
	route.GET("/", controllers.MainPage)
	route.GET("/image", controllers.ImagePage)

	err := route.Run()
	if err != nil {
		return
	}
}
