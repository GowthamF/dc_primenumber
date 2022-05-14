package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	isSpawn = flag.Bool("isspawn", false, "count")
)

// taskkill /IM main.exe /F
// docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management

func main() {
	//Do the election based on nodeId not based on InstanceId
	flag.Parse()
	if *isSpawn {
		for i := 0; i < 2; i++ {
			spawnProcess()
		}
	} else {
		for i := 0; i < 7; i++ {
			spawnProcess()
		}
	}
	// ReadFile()
	// WriteFile()
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
		number, _ := strconv.ParseInt(fscanner.Text(), 0, 32)
		numberToCheck := int32(number)
		log.Println(numberToCheck)
		NumbersToCheck = append(NumbersToCheck, numberToCheck)
	}

	// numberToCheck := int32(54322)

	// for _, number := range NumbersToCheck {
	// 	if *number != numberToCheck {
	// 		log.Println(*number)
	// 		NumbersToCheck = append(NumbersToCheck, number)
	// 	}
	// }
}

func WriteFile() {

	numberToCheck := 54322
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

func spawnProcess() {
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
