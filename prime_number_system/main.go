package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"dc_assignment.com/m/v2/controllers"
	"dc_assignment.com/m/v2/eurekaservices"
	"dc_assignment.com/m/v2/models"
	"dc_assignment.com/m/v2/queue"
	"dc_assignment.com/m/v2/routes"
)

var (
	nodeId = flag.String("nodeid", "", "ID")
)

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	portNumber := listener.Addr().(*net.TCPAddr).Port
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	currentTime := fmt.Sprint(time.Now().UnixMilli())
	randomNumber := fmt.Sprint(rand.Int31n(10000))
	id := currentTime + randomNumber

	hostName := "PRIMENUMBER"
	ipAddress := "localhost"
	port := portNumber
	status := "UP"
	enabledPort := "true"
	healthCheckUrl := "http://" + ipAddress + ":" + strconv.Itoa(port) + "/healthcheck"
	statusCheckUrl := "http://" + ipAddress + ":" + strconv.Itoa(port) + "/status"
	homePageUrl := "http://" + ipAddress + ":" + strconv.Itoa(port)
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
	go queue.ReceiveMessage(queue.MasterElectionMessage)
	go startElection(id, *nodeId)
	go eurekaservices.UpdateHeartBeat(*nodeId, id)
	r := routes.SetupRouter(id)
	r.RunListener(listener)
}

func startElection(id string, app string) {
	durationOfTime := time.Duration(30) * time.Second
	time.Sleep(durationOfTime)
	controllers.GetHigherInstanceIds(id, app)
}
