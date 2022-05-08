package controllers

import (
	"os"

	"dc_assignment.com/prime_number/v2/eurekaservices"
	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	eurekaservices.ShutdownNode(c.GetString("nodeId"))
	os.Exit(2)
}

func Status(c *gin.Context) {
	nodeId := c.Query("nodeId")
	instances := eurekaservices.GetInstances(nodeId)

	c.JSON(200, gin.H{
		"Status": instances.Application.Instance[0].Status,
		"NodeId": nodeId,
		"Role":   instances.Application.Instance[0].MetaData.Role,
	})
}
