package controllers

import (
	"os"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/services"
	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	services.ShutdownNode(c.GetString("nodeId"))
	CreateElectionLockFile(models.MasterLock)
	os.Exit(2)
}

func Status(c *gin.Context) {
	nodeId := c.Query("nodeId")
	instances := services.GetInstances(nodeId)

	c.JSON(200, gin.H{
		"Status": instances.Application.Instance[0].Status,
		"NodeId": nodeId,
		"Role":   instances.Application.Instance[0].MetaData.Role,
	})
}
