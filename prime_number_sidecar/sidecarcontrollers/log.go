package sidecarcontrollers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var mutex = &sync.Mutex{}

func LogLn(c *gin.Context) {
	mutex.Lock()
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	L := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	var message string
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &message)
	L.Println(message)
	file.Close()
	mutex.Unlock()
	c.Status(http.StatusOK)
}

func PrimeNumberLogLn(c *gin.Context) {
	file, _ := os.OpenFile("primeNumber.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	L := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	var message string
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &message)
	L.Println(message)
	file.Close()
	c.Status(http.StatusOK)
}
