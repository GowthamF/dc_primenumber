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
	"dc_assignment.com/prime_number/v2/nodemessage"
	"dc_assignment.com/prime_number/v2/routes"
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
	randomNumber := fmt.Sprint(rand.Int31n(10000))
	id := currentTime + randomNumber
	models.SidecarPortNumber = sidecarPortNumber
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

	eurekaservices.RegisterInstance(*nodeId, ins)
	go startElection(*nodeId)
	go nodemessage.ReceiveMessage(nodemessage.MasterElectionMessage)
	go eurekaservices.UpdateHeartBeat(*nodeId, id)
	r := routes.SetupRouter(*nodeId)
	r.Run("localhost:" + *appPortNumber)
}

func startElection(nodeId string) {
	durationOfTime := time.Duration(30) * time.Second
	time.Sleep(durationOfTime)
	controllers.GetHigherInstanceIds(nodeId)
}

// func spawnProcess() {
// 	durationOfTime := time.Duration(10) * time.Second
// 	time.Sleep(durationOfTime)
// 	cmd := exec.Command("bash", "app.sh")

// 	e := cmd.Run()
// 	eurekaservices.L.Println(e)
// }
