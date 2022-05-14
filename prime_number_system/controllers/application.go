package controllers

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/nodemessage"
	"dc_assignment.com/prime_number/v2/services"
	"dc_assignment.com/prime_number/v2/sidecar"
	"github.com/gin-gonic/gin"
)

func StopApplication(c *gin.Context) {
	services.ShutdownNode(c.GetString("nodeId"))
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
	go sidecar.Log(*models.NodeId + " Notification from  Proposer ")
	learnerNode := services.GetNodesByRole(models.LearnerNode)
	go sidecar.Log(strconv.Itoa(len(learnerNode)) + " Made it here")
	if len(learnerNode) > 0 {
		go sidecar.Log(" Made it here 1")
		if !*outcome.IsPrimeNumber {
			go sidecar.Log(" Made it here 2")
			if outcome.NumberToCheck != outcome.DivisibleNumber && outcome.DivisibleNumber != 0 && outcome.DivisibleNumber != 1 {
				reminder := outcome.NumberToCheck % outcome.DivisibleNumber
				if reminder == 0 {
					go sidecar.Log(" Made it here 3")
					go services.NotifyLearnerNodeProcess(learnerNode[0].Instance[0].HomePageUrl, outcome)
				}
			}

		} else {
			go sidecar.Log(" Made it here 4")
			go services.NotifyLearnerNodeProcess(learnerNode[0].Instance[0].HomePageUrl, outcome)
		}
	}

	c.JSON(http.StatusOK, outcome)
}

var outcomes []*models.PrimeNumbersValidationMessage = make([]*models.PrimeNumbersValidationMessage, 0)

func NotifyLearner(c *gin.Context) {
	var outcome *models.PrimeNumbersValidationMessage

	err := c.BindJSON(&outcome)

	if err != nil {
		go sidecar.Log(err.Error())
		c.AbortWithError(500, err)
	}
	outcomes = append(outcomes, outcome)
	go sidecar.Log(*models.NodeId + " Notification from  Acceptors ")
	go sidecar.Log("Total Number of Messages" + strconv.Itoa(len(outcomes)))
	isPrimeNumber := false
	numberToCheck := 0

	if len(outcomes) == int(*models.NumberOfProposers) {
		for _, o := range outcomes {
			if !*o.IsPrimeNumber {
				log.Println(strconv.Itoa(int(o.NumberToCheck)) + " is not a Prime Number")
				go sidecar.PrimeNumberLog(strconv.Itoa(int(o.NumberToCheck)) + " is not a Prime Number")
				break
			} else {
				isPrimeNumber = *o.IsPrimeNumber
				numberToCheck = int(o.NumberToCheck)
			}
		}

		if isPrimeNumber {
			log.Println(strconv.Itoa(numberToCheck) + " is a Prime Number")
			go sidecar.PrimeNumberLog(strconv.Itoa(numberToCheck) + " is a Prime Number")
		}
	}
	go SpawnProcess()
	nodemessage.SendMessage(nodemessage.NewNodeSpawned, "TRUE")
	c.JSON(http.StatusOK, outcomes)
}

func SpawnProcess() {
	durationOfTime := time.Duration(10) * time.Second
	time.Sleep(durationOfTime)
	cmd := exec.Command("bash", "app.sh")
	go cmd.Run()
}
