package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"dc_assignment.com/m/v2/eurekaservices"
	"dc_assignment.com/m/v2/models"
	"dc_assignment.com/m/v2/queue"
	"github.com/gin-gonic/gin"
)

func StartElection(c *gin.Context) {
	instanceId := c.GetString("instanceId")
	GetHigherInstanceIds(instanceId, models.PrimeNumberNode)
}

func StopElection(c *gin.Context) {}

func RequestElection(c *gin.Context) {
	requestInstanceId, _ := strconv.ParseInt(c.Param("requestInstanceId"), 0, 64)
	instanceId, _ := strconv.ParseInt(c.GetString("instanceId"), 0, 64)
	if instanceId < requestInstanceId {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}

}

func GetHigherInstanceIds(myId string, appName string) {
	instances := eurekaservices.GetInstances(appName)
	var urlToNotify string = ""
	for _, instance := range *instances.Application.Instance {
		if instance != nil {
			insId, _ := strconv.ParseInt(*instance.InstanceId, 0, 64)
			myAppId, _ := strconv.ParseInt(myId, 0, 64)

			if myAppId < insId {
				url := SendElectionRequest(instance.HomePageUrl, &myId)
				fmt.Println(url)
				if urlToNotify == "" {
					urlToNotify = url
				}
			}
		}
	}
	if urlToNotify != "" {
		SendStartElectionRequest(urlToNotify)
	} else {
		log.Println(myId, " I am the Leader")
		go queue.SendMessage(queue.MasterElectionMessage, myId+" is the Leader")
		for _, instance := range *instances.Application.Instance {
			if instance != nil {
				insId, _ := strconv.ParseInt(*instance.InstanceId, 0, 64)
				myAppId, _ := strconv.ParseInt(myId, 0, 64)
				if insId == myAppId {
					app := models.MasterNode
					ins := &models.InstanceModel{
						InstanceId: instance.InstanceId,
						HostName:   instance.HostName,
						App:        &app,
						IpAddress:  instance.IpAddress,
						Status:     instance.Status,
						Port: &models.PortModel{
							PortNumber: instance.Port.PortNumber,
							Enabled:    instance.Port.Enabled,
						},
						HealthCheckUrl: instance.HealthCheckUrl,
						StatusPageUrl:  instance.StatusPageUrl,
						HomePageUrl:    instance.HomePageUrl,
						DataCenterInfo: &models.DataCenterInfoModel{
							Class: instance.DataCenterInfo.Class,
							Name:  instance.DataCenterInfo.Name,
						},
					}
					eurekaservices.RegisterInstance(app, ins)
				}
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
			log.Println("Do not continue")
			remoteUrl = *url
		}
	}

	return remoteUrl
}
