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

	// [<snippet assets>]
	route.Static("/static/stylesheets", "static/stylesheets")
	route.Static("/static/javascripts", "static/javascripts")
	route.Static("/static/images", "static/images")
	// [<endsnippet assets>]

	// [<snippet templates>]
	route.LoadHTMLGlob("templates/*")
	// [<endsnippet templates>]

	route.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong\n")
	})

	// [<snippet routes>]
	route.GET("/", controllers.MainPage)
	route.GET("/image", controllers.ImagePage)
	// [<endsnippet routes>]

	err := route.Run()
	if err != nil {
		return
	}
}
