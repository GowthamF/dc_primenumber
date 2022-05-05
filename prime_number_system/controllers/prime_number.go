package controllers

import "fmt"

func CheckIfPrimeNumber(numberToCheck int32, startRange int32, endRange int32) {

	for i := startRange; i < endRange; i++ {
		if numberToCheck != i && i != 0 && i != 1 {
			reminder := numberToCheck % i

			if reminder == 0 {
				fmt.Println(numberToCheck, "can be divided by", i)
				break
			}
		}
	}
}
