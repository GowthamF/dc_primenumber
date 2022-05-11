package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
)

func NotifyLearnerNode(url *string, nProposers int32) {
	json_data, err := json.Marshal(nProposers)
	if err != nil {
		go sidecar.Log(*models.NodeId + err.Error())
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", *url+"/proposersCount", bytes.NewBuffer(json_data))

	if err != nil {
		log.Println(err)
		go sidecar.Log(*models.NodeId + err.Error())
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		go sidecar.Log(*models.NodeId + err.Error())
	}
	if resp != nil {
		if resp.StatusCode == 200 {
			log.Println("Sent")
			go sidecar.Log(*models.NodeId + " Proposers count has been notified")

		} else {
			log.Println("Failed")
			go sidecar.Log(*models.NodeId + " Proposers count has been notified failed")
		}
	}
}

func NotifyAcceptorNode(url *string, outcome *models.PrimeNumbersValidationMessage) {
	json_data, err := json.Marshal(outcome)
	if err != nil {
		go sidecar.Log(*models.NodeId + err.Error())
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", *url+"/notifyAcceptor", bytes.NewBuffer(json_data))

	if err != nil {
		log.Println(err)
		go sidecar.Log(*models.NodeId + err.Error())
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		go sidecar.Log(*models.NodeId + err.Error())
	}
	if resp != nil {
		if resp.StatusCode == 200 {
			log.Println("Sent")
			go sidecar.Log(*models.NodeId + " outcome has been sent to Acceptor")

		} else {
			log.Println("Failed")
			go sidecar.Log(*models.NodeId + " outcome has been sent to Acceptor failed")
		}
	}
}
