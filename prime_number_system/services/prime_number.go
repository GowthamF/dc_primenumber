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
	"sync"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
)

func CheckIfPrimeNumber(numberToCheck int32, startRange int32, endRange int32) *models.PrimeNumbersValidationMessage {
	isPrimeNumber := false
	var message *models.PrimeNumbersValidationMessage = &models.PrimeNumbersValidationMessage{NumberToCheck: numberToCheck, IsPrimeNumber: &isPrimeNumber}
	go sidecar.Log(*models.NodeId + " Start Range :" + strconv.FormatInt(int64(startRange), 10) + "End Range :" + strconv.FormatInt(int64(endRange), 10))
	for i := startRange; i < endRange; i++ {
		if numberToCheck != i && i != 0 && i != 1 {
			reminder := numberToCheck % i

			if reminder == 0 {
				go sidecar.Log("Node ID " + *models.NodeId + " says " + strconv.FormatInt(int64(numberToCheck), 10) + " can be divided by" + strconv.FormatInt(int64(i), 10))
				isPrimeNumber = false
				message.IsPrimeNumber = &isPrimeNumber
				message.DivisibleNumber = i
				break
			} else if reminder == 1 {
				isPrimeNumber = true
				message.IsPrimeNumber = &isPrimeNumber
			}
		}
	}

	return message
}

var NumbersToCheck []int32 = []int32{}

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
		log.Println(fscanner.Text())
		number, _ := strconv.ParseInt(fscanner.Text(), 0, 32)
		numberToCheck := int32(number)
		NumbersToCheck = append(NumbersToCheck, numberToCheck)
	}
}

func WriteFile(numberToCheck int32) {
	file, err := os.Create("prime_numbers.txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)

	for _, number := range NumbersToCheck {
		if number != int32(numberToCheck) {
			_, _ = datawriter.WriteString(strconv.Itoa(int(number)) + "\n")
		}
	}

	datawriter.Flush()
	file.Close()
}

var isWritten = false
var mutex = &sync.Mutex{}

func StartProcess(index int) {

	if index == 4 {
		go sidecar.Log("Everything is done")
		return
	}

	ReadFile()
	nodes := GetNodesByRole(models.ProposerNode)
	learnerNodes := GetNodesByRole(models.LearnerNode)
	go sidecar.Log("Made it here in Start Process 1")
	if len(learnerNodes) > 0 {
		go sidecar.Log("Made it here in Start Process 2")
		numberToCheck := NumbersToCheck[index]
		go sidecar.Log(strconv.Itoa(int(numberToCheck)))
		log.Println(numberToCheck, "NUMBERS")
		NotifyLearnerNode(learnerNodes[0].Instance[0].HomePageUrl, int32(len(nodes)))

		rangeToAssign := math.Round(float64(numberToCheck) / float64(len(nodes)))

		startRange := 0

		for _, node := range nodes {
			endRange := startRange + int(rangeToAssign) - 1
			go sendRequest(node.Instance[0].HomePageUrl, numberToCheck, int32(startRange), int32(endRange))
			startRange = endRange + 1
		}

		// mutex.Lock()
		// WriteFile(numberToCheck)
		// mutex.Unlock()
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
