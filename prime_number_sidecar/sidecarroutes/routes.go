package sidecarroutes

import (
	"net/http"

	"dc_assignment.com/sidecar/v2/sidecarcontrollers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routeEngine := gin.Default()
	routeEngine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<div><h1>Welcome to Prime Number Side Car.</h1><h2>Everything is fine...</h2><h3>Made in &#128151; with Go </h3></div>"))
	})

	routeEngine.POST("/log", sidecarcontrollers.LogLn)

	return routeEngine
}
