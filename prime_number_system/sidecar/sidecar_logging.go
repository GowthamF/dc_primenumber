package sidecar

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"dc_assignment.com/prime_number/v2/models"
)

func Log(message string) {
	json_data, err := json.Marshal(message)
	client := &http.Client{}

	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+*models.SidecarPortNumber+"/log", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	go client.Do(req)
}

func PrimeNumberLog(message string) {
	json_data, err := json.Marshal(message)
	client := &http.Client{}

	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+*models.SidecarPortNumber+"/primeNumberLog", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	go client.Do(req)
}
