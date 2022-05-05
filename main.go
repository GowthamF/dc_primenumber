package main

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	//Do the election based on nodeId not based on InstanceId

	for i := 0; i < 5; i++ {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		appPort := listener.Addr().(*net.TCPAddr).Port
		listener.Close()

		listener1, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		sideCarPort := listener1.Addr().(*net.TCPAddr).Port
		listener1.Close()
		rand.Seed(time.Now().UnixNano())
		currentTime := fmt.Sprint(time.Now().UnixMilli())
		randomNumber := fmt.Sprint(rand.Int31n(10000))
		id := currentTime + randomNumber

		cmd := exec.Command("bash", "app.sh", id, strconv.Itoa(appPort), strconv.Itoa(sideCarPort))

		go cmd.Run()
	}
}
