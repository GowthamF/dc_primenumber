package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"dc_assignment.com/prime_number/v2/controllers"
	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/nodemessage"
	"dc_assignment.com/prime_number/v2/routes"
	"dc_assignment.com/prime_number/v2/services"
	"dc_assignment.com/prime_number/v2/sidecar"
)

var (
	nodeId            = flag.String("nodeid", "", "ID")
	appPortNumber     = flag.String("portnumber", "8080", "App Port Number")
	sidecarPortNumber = flag.String("sidecarportnumber", "8081", "Sidecar Port Number")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	currentTime := fmt.Sprint(time.Now().UnixMilli())
	randomNumber := fmt.Sprint(100 + rand.Intn(999-100))
	id := currentTime + randomNumber
	models.SidecarPortNumber = sidecarPortNumber
	models.NodeId = nodeId
	hostName := "PRIMENUMBER"
	ipAddress := "localhost"
	port, _ := strconv.ParseInt(*appPortNumber, 0, 64)
	status := models.UpStatus
	enabledPort := "true"
	healthCheckUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 10) + "/healthcheck"
	statusCheckUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 10) + "/status?nodeId=" + *nodeId
	homePageUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 10)
	class := "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"
	name := "MyOwn"
	ins := &models.InstanceModel{
		InstanceId: &id,
		HostName:   &hostName,
		App:        nodeId,
		IpAddress:  &ipAddress,
		Status:     &status,
		Port: &models.PortModel{
			PortNumber: &port,
			Enabled:    &enabledPort,
		},
		HealthCheckUrl: &healthCheckUrl,
		StatusPageUrl:  &statusCheckUrl,
		HomePageUrl:    &homePageUrl,
		DataCenterInfo: &models.DataCenterInfoModel{
			Class: &class,
			Name:  &name,
		},
		MetaData: &models.MetaDataModel{
			Role: &models.PrimeNumberNode,
		},
	}

	services.RegisterInstance(*nodeId, ins)
	ch1 := make(chan string)
	ch2 := make(chan string)
	go nodemessage.ReceiveMessage(ch1, nodemessage.MasterElectionMessage)
	go nodemessage.ReceiveMessage(ch2, nodemessage.NewNodeSpawned)
	go func() {
		id := <-ch1
		models.MasterNodeId = &id
		sidecar.Log(<-ch1)
	}()
	go func() {
		sidecar.Log(<-ch2)
	}()

	go services.UpdateHeartBeat(*nodeId, id)
	go startElection(*nodeId)
	r := routes.SetupRouter(*nodeId)
	r.Run("localhost:" + *appPortNumber)
}

func startElection(nodeId string) {
	durationOfTime := time.Duration(30) * time.Second
	time.AfterFunc(durationOfTime, func() { controllers.GetHigherInstanceIds(nodeId) })
}
