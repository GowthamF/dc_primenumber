package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	c.JSON(200, "Application terminated")
	os.Exit(2)
}
