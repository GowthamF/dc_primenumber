package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

func main() {
	//Do the election based on

	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		currentTime := fmt.Sprint(time.Now().UnixMilli())
		randomNumber := fmt.Sprint(rand.Int31n(10000))
		id := currentTime + randomNumber

		cmd := exec.Command("bash", "app.sh", id)

		go cmd.Run()
	}
}
