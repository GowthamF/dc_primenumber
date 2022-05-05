package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	currentTime := fmt.Sprint(time.Now().UnixMilli())
	randomNumber := fmt.Sprint(rand.Int31n(10000))
	id := currentTime + randomNumber

	cmd := exec.Command("bash", "app.sh", id)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
