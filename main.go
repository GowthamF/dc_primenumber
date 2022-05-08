package main

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"time"
)

// taskkill /IM main.exe /F
// docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management

func main() {
	//Do the election based on nodeId not based on InstanceId

	for i := 0; i < 7; i++ {
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
		randomNumber := fmt.Sprint(100 + rand.Intn(999-100))
		id := currentTime + randomNumber
		fmt.Println(appPort, sideCarPort)

		cmd := exec.Command("bash", "app.sh", id, strconv.Itoa(appPort), strconv.Itoa(sideCarPort))
		cmd1 := exec.Command("bash", "sidecar.sh", strconv.Itoa(appPort), strconv.Itoa(sideCarPort))

		go cmd.Run()
		go cmd1.Run()
	}
}
