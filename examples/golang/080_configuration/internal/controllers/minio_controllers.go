package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"werf_guide_app/internal/services"
)

func UploadController(c *gin.Context) {
	ctx := context.Background()

	file, _ := c.FormFile("file")
	buffer, err := file.Open()

	if err != nil {
		log.Fatalln(err)
	}
	defer buffer.Close()

	minioClient, err := services.MinioConnection()
	if err != nil {
		log.Fatalln(err)
	}
	bucketName := services.GetBucketName()

	objectName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	_, err = minioClient.PutObject(ctx, bucketName, objectName, fileBuffer,
		fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		log.Fatalln(err)
	}

	c.String(http.StatusOK, "File uploaded.\n")
}

func DownloadController(c *gin.Context) {
	minioClient, err := services.MinioConnection()
	if err != nil {
		log.Fatalln(err)
	}
	bucketName := services.GetBucketName()

	object, err := minioClient.GetObject(context.Background(), bucketName,
		"file.txt", minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
	}
	defer object.Close()

	b, err := io.ReadAll(object)
	if err != nil {
		c.String(http.StatusOK, "You haven't uploaded anything yet.\n")
	}

	c.String(http.StatusOK, string(b))
}
