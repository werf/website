package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MainPage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{})
}

func ImagePage(context *gin.Context) {
	context.HTML(http.StatusOK, "image.html", gin.H{})
}
