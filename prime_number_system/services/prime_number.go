package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
)

func CheckIfPrimeNumber(numberToCheck int32, startRange int32, endRange int32) *models.PrimeNumbersValidationMessage {
	isPrimeNumber := true
	var message *models.PrimeNumbersValidationMessage = &models.PrimeNumbersValidationMessage{NumberToCheck: numberToCheck, IsPrimeNumber: &isPrimeNumber}
	for i := startRange; i < endRange; i++ {
		if numberToCheck != i && i != 0 && i != 1 {
			reminder := numberToCheck % i

			if reminder == 0 {
				go sidecar.Log("Node ID " + *models.NodeId + " says" + strconv.FormatInt(int64(numberToCheck), 10) + " can be divided by" + strconv.FormatInt(int64(i), 10))
				isPrimeNumber = false
				message.IsPrimeNumber = &isPrimeNumber
				message.DivisibleNumber = i
				break
			}
		}
	}

	return message
}

var NumbersToCheck []*int32 = []*int32{}

func ReadFile() {
	file, err := os.Open("prime_numbers.txt")

	if err != nil {
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		number, _ := strconv.ParseInt(fscanner.Text(), 0, 32)
		numberToCheck := int32(number)
		NumbersToCheck = append(NumbersToCheck, &numberToCheck)
	}

}

func StartProcess() {
	nodes := GetNodesByRole(models.ProposerNode)
	learnerNodes := GetNodesByRole(models.LearnerNode)

	if len(learnerNodes) > 0 {
		numberToCheck := *NumbersToCheck[0]
		NotifyLearnerNode(learnerNodes[0].Instance[0].HomePageUrl, int32(len(nodes)))

		rangeToAssign := math.Round(float64(numberToCheck) / float64(len(nodes)))

		startRange := 0

		for _, node := range nodes {
			endRange := startRange + int(rangeToAssign) - 1
			go sendRequest(node.Instance[0].HomePageUrl, numberToCheck, int32(startRange), int32(endRange))
			startRange = endRange + 1
		}

		for _, number := range NumbersToCheck {
			if number != &numberToCheck {
				NumbersToCheck = append(NumbersToCheck, number)
			}
		}
	}

}

func sendRequest(url *string, numberToCheck int32, startRange int32, endRange int32) {
	client := &http.Client{}
	primeNumberModel := models.PrimeNumbers{
		NumberToCheck: numberToCheck,
		StartRange:    startRange,
		EndRange:      endRange,
	}

	json_data, err := json.Marshal(primeNumberModel)

	if err != nil {
		log.Println(err)
		go sidecar.Log(err.Error())
	}

	sendPrimeNumberCheckRequest := func() {
		req, err := http.NewRequest("POST", *url+"/checkPrimeNumber", bytes.NewBuffer(json_data))

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
				go sidecar.Log(*models.NodeId + " Check Prime Number Request has been sent")

			} else {
				log.Println("Failed")
				go sidecar.Log(*models.NodeId + " Check Prime Number Request sent failed")
			}
		}

	}

	go sendPrimeNumberCheckRequest()
}
