package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"dc_assignment.com/prime_number/v2/controllers"
	"dc_assignment.com/prime_number/v2/eurekaservices"
	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/routes"
)

var (
	nodeId            = flag.String("nodeid", "", "ID")
	appPortNumber     = flag.String("appPortNumber", "8080", "App Port Number")
	sideCarPortNumber = flag.String("sideCarPortNumber", "0", "Sidecar Port Number")
)

func main() {

	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	currentTime := fmt.Sprint(time.Now().UnixMilli())
	randomNumber := fmt.Sprint(rand.Int31n(10000))
	id := currentTime + randomNumber

	hostName := "PRIMENUMBER"
	ipAddress := "localhost"
	port, _ := strconv.ParseInt(*appPortNumber, 0, 64)
	status := "UP"
	enabledPort := "true"
	healthCheckUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 16) + "/healthcheck"
	statusCheckUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 16) + "/status"
	homePageUrl := "http://" + ipAddress + ":" + strconv.FormatInt(port, 16)
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
	}

	eurekaservices.RegisterInstance(*nodeId, ins)
	// go queue.ReceiveMessage(queue.MasterElectionMessage)
	go startElection(id, *nodeId)
	go eurekaservices.UpdateHeartBeat(*nodeId, id)
	r := routes.SetupRouter(id)
	r.Run(":" + *appPortNumber)
}

func startElection(id string, app string) {
	durationOfTime := time.Duration(30) * time.Second
	time.Sleep(durationOfTime)
	controllers.GetHigherInstanceIds(id, app)
}
