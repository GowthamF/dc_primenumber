package eurekaservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"dc_assignment.com/m/v2/models"
	"github.com/carlescere/scheduler"
)

func RegisterInstance(appName string, instance *models.InstanceModel) {
	instanceJson := map[string]*models.InstanceModel{"instance": instance}
	fmt.Println(instanceJson)
	json_data, err := json.Marshal(instanceJson)
	client := &http.Client{}
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8761/eureka/apps/"+appName, bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode == 204 {
		log.Println("Successfully Registered")

	} else {
		log.Println(resp.StatusCode)
		log.Println("Error during registering")
	}
}

func GetInstances(appName string) *models.InstancesModel {
	client := &http.Client{}
	var apps *models.InstancesModel
	req, err := http.NewRequest("GET", "http://localhost:8761/eureka/apps/"+appName, bytes.NewBuffer([]byte{}))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &apps)

		if err != nil {
			log.Fatalln(err)
		}

	} else {
		log.Println("Error during retrieving")
	}
	return apps
}

func UpdateHeartBeat(appName string, instanceId string) {
	client := &http.Client{}
	job := func() {

		req, err := http.NewRequest("PUT", "http://localhost:8761/eureka/apps/"+appName+"/"+instanceId, bytes.NewBuffer([]byte{}))

		if err != nil {
			log.Fatalln(err)
		}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode == 200 {
			log.Println("Heart beat updated")

		} else {
			log.Println("Heart beat failed")
		}
	}

	scheduler.Every(15).Seconds().Run(job)
	runtime.Goexit()
}

func GetInstanceIds(app string) []*string {
	durationOfTime := time.Duration(30) * time.Second
	time.Sleep(durationOfTime)
	var instanceIds = []*string{}
	instances := GetInstances(app)
	if instances != nil {
		for _, instance := range *instances.Application.Instance {
			if instance != nil {
				instanceIds = append(instanceIds, instance.InstanceId)
			}
		}
	}

	return instanceIds
}
