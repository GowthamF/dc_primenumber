package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"dc_assignment.com/prime_number/v2/eurekaservices"
	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
	"github.com/gin-gonic/gin"
)

func StartElection(c *gin.Context) {
	instanceId := c.GetString("nodeId")
	GetHigherInstanceIds(instanceId)
}

func StopElection(c *gin.Context) {}

func RequestElection(c *gin.Context) {
	requestInstanceId, _ := strconv.ParseInt(c.Param("requestInstanceId"), 0, 64)
	instanceId, _ := strconv.ParseInt(c.GetString("nodeId"), 0, 64)
	if instanceId < requestInstanceId {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}

}

func GetHigherInstanceIds(myId string) {
	eurekaservices.GetNodes()
	nodes := eurekaservices.GetNodes()
	var urlToNotify string = ""
	var selectedNode *models.ApplicationModel
	for _, instance := range nodes {
		if instance != nil {
			if *instance.Instance[0].Status == models.UpStatus {
				nodeId, _ := strconv.ParseInt(*instance.Name, 0, 64)
				myNodeId, _ := strconv.ParseInt(myId, 0, 64)

				if myNodeId < nodeId {
					url := SendElectionRequest(instance.Instance[0].HomePageUrl, &myId)
					if urlToNotify == "" && selectedNode != nil {
						urlToNotify = url
						selectedNode = instance
					}
				}
			}

		}
	}
	if urlToNotify != "" {
		SendStartElectionRequest(urlToNotify)
	} else {
		sidecar.Log(myId + "I am the leader")
		// go queue.SendMessage(queue.MasterElectionMessage, myId+" is the Leader")
		eurekaservices.UpdateRole(myId, models.MasterNode)
	}

}

func SendStartElectionRequest(url string) {
	fmt.Println(url)
	_, err := http.Get(url + "/startElection")

	if err != nil {
		log.Fatalln(err)
	}
}

func SendElectionRequest(url *string, myAppId *string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", *url+"/requestElection/"+*myAppId, bytes.NewBuffer([]byte{}))
	var remoteUrl = ""
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	if resp != nil {
		if resp.StatusCode == 200 {
			log.Println("Continue the election")

		} else if resp.StatusCode == 406 {
			log.Println("Do not continue")
			remoteUrl = *url
		}
	}

	return remoteUrl
}
