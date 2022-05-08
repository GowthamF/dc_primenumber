package services

import (
	"fmt"
	"math"
	"strconv"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/sidecar"
)

func CheckIfPrimeNumber(numberToCheck int32, startRange int32, endRange int32) {

	for i := startRange; i < endRange; i++ {
		if numberToCheck != i && i != 0 && i != 1 {
			reminder := numberToCheck % i

			if reminder == 0 {
				go sidecar.Log(strconv.FormatInt(int64(numberToCheck), 10) + " can be divided by" + strconv.FormatInt(int64(i), 10))
				break
			}
		}
	}
}

func StartProcess(numberToCheck int32) {
	nodes := GetNodesByRole(models.ProposerNode)

	rangeToAssign := math.Round(float64(numberToCheck) / float64(len(nodes)))

	startRange := 0

	for _, node := range nodes {
		endRange := startRange + int(rangeToAssign) - 1
		go CheckIfPrimeNumber(numberToCheck, int32(startRange), int32(endRange))
		startRange = endRange + 1
		fmt.Println(node.Name)
	}

}
