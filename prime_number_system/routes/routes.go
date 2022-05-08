package routes

import (
	"net/http"

	"dc_assignment.com/prime_number/v2/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(nodeId string) *gin.Engine {
	routeEngine := gin.Default()
	routeEngine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<div><h1>Welcome to Prime Number Distributed Sytem.</h1><h2>Everything is fine...</h2><h3>Made in &#128151; with Go </h3></div>"))
	})

	routeEngine.Use(func(ctx *gin.Context) {
		ctx.Set("nodeId", nodeId)
		ctx.Next()
	})

	routeEngine.GET("/status", controllers.Status)

	routeEngine.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Running")
	})

	routeEngine.GET("/startElection", controllers.StartElection)
	routeEngine.POST("/stopElection", controllers.StopElection)
	routeEngine.GET("/requestElection/:requestInstanceId", controllers.RequestElection)
	routeEngine.POST("/stopApplication", controllers.StopApplication)

	return routeEngine
}
