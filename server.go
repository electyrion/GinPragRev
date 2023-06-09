package main

import (
	"golang-gin-poc/controller"
	"golang-gin-poc/middlewares"
	"golang-gin-poc/repository"
	"golang-gin-poc/service"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	
	videoService service.VideoService = service.New(videoRepository)
	loginService service.LoginService = service.NewLoginService()
	jwtService service.JWTService = service.NewJWTService()
	
	videoController controller.VideoController = controller.New(videoService)
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	// this will be executed when main func finished
	defer videoRepository.CloseDB()

	setupLogOutput()

	server := gin.New()

	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("templates/*.html")

	server.Use(gin.Recovery(), middlewares.Logger(), gindump.Dump())

	// Login Endpoint: Authentication + Token creation
	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	// JWT Authorization Middleware applies to "/api" only
	apiRoutes := server.Group("/api", middlewares.AuthorizeJWT())
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

		apiRoutes.PUT("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video updated successfully"})
			}
		})

		apiRoutes.DELETE("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
			}
		})
	}

	// the "/view" endpoints are public (no Authorization requirec)
	viewRoutes := server.Group("/view") 
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}
	server.Run(":8080")
}