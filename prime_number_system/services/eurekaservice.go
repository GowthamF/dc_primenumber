package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
	"github.com/carlescere/scheduler"
)

func RegisterInstance(appName string, instance *models.InstanceModel) {
	instanceJson := map[string]*models.InstanceModel{"instance": instance}
	json_data, err := json.Marshal(instanceJson)
	client := &http.Client{}
	if err != nil {
		sidecar.Log(appName + err.Error())
	}

	req, err := http.NewRequest("POST", "http://localhost:8761/eureka/apps/"+appName, bytes.NewBuffer(json_data))

	if err != nil {
		sidecar.Log(appName + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		sidecar.Log(appName + err.Error())
	}

	if resp.StatusCode == 204 {
		sidecar.Log(appName + " Successfully Registered")

	} else {
		sidecar.Log("Error during registering")
	}
}

func GetNodes() []*models.ApplicationModel {
	client := &http.Client{}
	apps := map[string]*models.NodesModel{}
	var nodes []*models.ApplicationModel = []*models.ApplicationModel{}
	req, err := http.NewRequest("GET", "http://localhost:8761/eureka/apps", bytes.NewBuffer([]byte{}))

	if err != nil {
		sidecar.Log(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		sidecar.Log(err.Error())
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			sidecar.Log(err.Error())
		}

		err = json.Unmarshal(body, &apps)
		nodes = apps["applications"].Nodes

		if err != nil {
			sidecar.Log(err.Error())
		}

	} else {
		sidecar.Log("Error during retrieving")
	}

	return nodes
}

func GetInstances(appName string) *models.InstancesModel {
	client := &http.Client{}
	var apps *models.InstancesModel
	req, err := http.NewRequest("GET", "http://localhost:8761/eureka/apps/"+appName, bytes.NewBuffer([]byte{}))

	if err != nil {
		sidecar.Log(appName + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		sidecar.Log(appName + err.Error())
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			sidecar.Log(appName + err.Error())
		}

		err = json.Unmarshal(body, &apps)

		if err != nil {
			sidecar.Log(appName + err.Error())
		}

	} else {
		sidecar.Log(appName + " Error during retrieving")
	}
	return apps
}

func UpdateHeartBeat(appName string, instanceId string) {
	client := &http.Client{}
	job := func() {

		req, err := http.NewRequest("PUT", "http://localhost:8761/eureka/apps/"+appName+"/"+instanceId, bytes.NewBuffer([]byte{}))

		if err != nil {
			sidecar.Log(appName + ":" + instanceId + err.Error())
		}
		resp, err := client.Do(req)

		if err != nil {
			sidecar.Log(appName + ":" + instanceId + err.Error())
		}

		if resp.StatusCode == 200 {
			// sidecar.Log(appName + ":" + instanceId + " Heart beat updated")

		} else {
			sidecar.Log(appName + ":" + instanceId + " Heart beat failed")
		}
	}

	scheduler.Every(15).Seconds().Run(job)
	runtime.Goexit()
}

func GetInstanceIds(app string) []*string {
	// durationOfTime := time.Duration(30) * time.Second
	// time.Sleep(durationOfTime)
	var instanceIds = []*string{}
	instances := GetInstances(app)
	if instances != nil {
		for _, instance := range instances.Application.Instance {
			if instance != nil {
				instanceIds = append(instanceIds, instance.InstanceId)
			}
		}
	}

	return instanceIds
}

func UpdateRole(app string, role string) {
	instanceIds := GetInstanceIds(app)
	client := &http.Client{}

	updateRole := func(instanceId string) {
		req, err := http.NewRequest("PUT", "http://localhost:8761/eureka/apps/"+app+"/"+instanceId+"/metadata?role="+role, bytes.NewBuffer([]byte{}))

		if err != nil {
			sidecar.Log(app + ":" + instanceId + err.Error())
		}
		resp, err := client.Do(req)

		if err != nil {
			sidecar.Log(app + ":" + instanceId + err.Error())
		}

		if resp.StatusCode == 200 {
			sidecar.Log(app + ":" + instanceId + " Role updated")

		} else {
			sidecar.Log(app + ":" + instanceId + " Role update failed")
		}
	}

	for _, instanceId := range instanceIds {
		if instanceId != nil {
			go updateRole(*instanceId)
		}
	}
}

func GetNodesByRole(role string) []*models.ApplicationModel {
	nodes := GetNodes()
	var nodesByRoles []*models.ApplicationModel = []*models.ApplicationModel{}
	for _, node := range nodes {
		nodeRole := node.Instance[0].MetaData.Role
		if *nodeRole == role {
			if *CheckIfNodeIsAlive(node.Instance[0].StatusPageUrl) {
				nodesByRoles = append(nodesByRoles, node)
			}
		}
	}

	return nodesByRoles
}

func ShutdownNode(appName string) {
	instances := GetInstanceIds(appName)

	client := &http.Client{}
	deleteIntance := func(instanceId string) {
		req, err := http.NewRequest("DELETE", "http://localhost:8761/eureka/apps/"+appName+"/"+instanceId, bytes.NewBuffer([]byte{}))

		if err != nil {
			fmt.Println(err)
			sidecar.Log(appName + ":" + instanceId + err.Error())
		}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			sidecar.Log(appName + ":" + instanceId + err.Error())
		}

		if resp.StatusCode == 200 {
			fmt.Println("DELETED")
			sidecar.Log(appName + ":" + instanceId + " Deleted")

		} else {
			fmt.Println("FAILED")
			sidecar.Log(appName + ":" + instanceId + " Deleting failed")
		}
	}

	for _, instanceId := range instances {
		if instanceId != nil {
			go deleteIntance(*instanceId)
		}
	}
}

func CheckIfNodeIsAlive(nodeUrl *string) *bool {
	resp, err := http.Get(*nodeUrl)
	var isRunning bool = true
	if err != nil {
		isRunning = false
		sidecar.Log(*nodeUrl + ": " + err.Error())
	} else {
		if resp.StatusCode == 200 {
			isRunning = true
			sidecar.Log(*nodeUrl + ": Running")
		}
	}
	return &isRunning
}
