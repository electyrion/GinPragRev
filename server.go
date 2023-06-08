package main

import (
	"golang-gin-poc/controller"
	"golang-gin-poc/middlewares"
	"golang-gin-poc/service"
	"io"
	"net/http"
	"os"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoService service.VideoService = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func extractVideoID(url string) (string, error) {
	regexPattern := `v=([^&]+)`

	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(url)

	if len(matches) >= 2 {
		videoID := matches[1]
		return videoID, nil
	}

	return "", fmt.Errorf("video ID not found in URL")
}

func main() {

	url := "https://www.youtube.com/watch?v=sDJLQMZzzM4&list=PL3eAkoh7fypr8zrkiygiY1e9osoqjoV9w&index=4"

	videoID, err := extractVideoID(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Video ID:", videoID)

	setupLogOutput()

	server := gin.New()

	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("templates/*.html")

	server.Use(gin.Recovery(), middlewares.Logger(), 
		middlewares.BasicAuth(), gindump.Dump())

	server.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	apiRoutes := server.Group("/api")
	{
		apiRoutes.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST("/videos", func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video input is valid"})
			}
		})
	}

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}
	server.Run(":8080")
}