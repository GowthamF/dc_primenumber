package controllers

import (
	"os"
	"strconv"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/services"
	"dc_assignment.com/prime_number/v2/sidecar"
	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	services.ShutdownNode(c.GetString("nodeId"))
	RemoveElectionLockFile(models.MasterLock)
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

func ProposerCount(c *gin.Context) {
	var proposersCount *int32

	err := c.BindJSON(&proposersCount)

	if err != nil {
		sidecar.Log(err.Error())
		c.AbortWithError(500, err)
	}

	models.NumberOfProposers = proposersCount

	go sidecar.Log("Number of Proposers " + strconv.Itoa(int(*proposersCount)))
}

func NotifyAcceptor(c *gin.Context) {
	var outcome *models.PrimeNumbersValidationMessage

	err := c.BindJSON(&outcome)

	if err != nil {
		go sidecar.Log(err.Error())
		c.AbortWithError(500, err)
	}
	go sidecar.Log(*models.NodeId + " Notification from  Proposer" + strconv.FormatBool(*outcome.IsPrimeNumber))
}
