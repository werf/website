package app

import (
	"net/http"
	"werf_guide_app/internal/common"
	"werf_guide_app/internal/controllers"

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

	// [<snippet db_routes>]
	route.GET("/remember", controllers.RememberController)
	route.GET("/say", controllers.SayController)
	// [<endsnippet db_routes>]

	// [<snippet minio_routes>]
	route.POST("/upload", controllers.UploadController)
	route.GET("/download", controllers.DownloadController)
	// [<endsnippet minio_routes>]

	err := route.Run()
	if err != nil {
		return
	}
}
