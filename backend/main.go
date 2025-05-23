package main

import (
	"runner-api/controllers"
	"runner-api/docs"
	"runner-api/state"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vue dev server default port
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/"

	containerState := state.NewContainerState()

	containerState.CreateContainers()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("api/ping", controllers.HealthCheck)
	r.GET("api/languages", controllers.GetSupportedLanguages)
	r.POST("api/runcode", controllers.RunCode)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
