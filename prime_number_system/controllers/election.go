package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"dc_assignment.com/prime_number/v2/eurekaservices"
	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/nodemessage"
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
	masterNodes := eurekaservices.GetNodesByRole(models.MasterNode)

	if len(masterNodes) > 0 {
		isMasterNodeRunning := false
		for _, masterNode := range masterNodes {
			isMasterNodeRunning = *eurekaservices.CheckIfNodeIsAlive(masterNode.Instance[0].StatusPageUrl)
		}

		if isMasterNodeRunning {
			return
		}
	}
	hasLockFileCreated := *createElectionLockFile()
	if hasLockFileCreated {
		nodes := eurekaservices.GetNodes()
		var urlToNotify string = ""
		var selectedNode *models.ApplicationModel
		for _, instance := range nodes {
			if instance != nil {
				if *instance.Instance[0].Status == models.UpStatus {
					nodeId, _ := strconv.ParseInt(*instance.Name, 0, 64)
					myNodeId, _ := strconv.ParseInt(myId, 0, 64)
					fmt.Println(myNodeId)
					fmt.Println(nodeId)
					if myNodeId < nodeId {
						url := SendElectionRequest(instance.Instance[0].HomePageUrl, &myId)
						if urlToNotify == "" && selectedNode == nil {
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
			removeElectionLockFile()
			go nodemessage.SendMessage(nodemessage.MasterElectionMessage, myId+" is the Leader")
			go sidecar.Log(myId + " I am the leader")
			go eurekaservices.UpdateRole(myId, models.MasterNode)
			go assignRoles(&myId)
		}
	}
}

func assignRoles(masterNodeId *string) {
	nodes := eurekaservices.GetNodes()
	var acceptorNodeIds []*string = []*string{}
	var learnerNodeIds []*string = []*string{}

	for _, node := range nodes {
		if *node.Name != *masterNodeId {
			if len(acceptorNodeIds) < 2 {
				acceptorNodeIds = append(acceptorNodeIds, node.Name)
				go eurekaservices.UpdateRole(*node.Name, models.AcceptorNode)
			} else if len(learnerNodeIds) < 1 {
				learnerNodeIds = append(learnerNodeIds, node.Name)
				go eurekaservices.UpdateRole(*node.Name, models.LearnerNode)
			} else {
				go eurekaservices.UpdateRole(*node.Name, models.ProposerNode)
			}
		}
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
			removeElectionLockFile()
			log.Println("Do not continue")
			remoteUrl = *url
		}
	}
	return remoteUrl
}

func createElectionLockFile() *bool {
	var (
		lockstate bool = false
	)

	if _, err := os.Stat("ms.lock"); err == nil {
		return &lockstate

	} else if os.IsNotExist(err) {
		var file, err = os.Create("ms.lock")
		if err != nil {
			return &lockstate
		}
		file.Close()
		lockstate = true
	}

	return &lockstate
}

func removeElectionLockFile() {
	_, err := os.Stat("ms.lock")
	if err == nil || os.IsExist(err) {
		var err = os.Remove("ms.lock")
		if err != nil {
			fmt.Println("Error removing file: ", err)
		}
	}
}
