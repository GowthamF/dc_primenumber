package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	os.Exit(2)
}
