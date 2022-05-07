package sidecarcontrollers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	file, _ = os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	L       = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func LogLn(c *gin.Context) {
	var message string
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &message)

	L.Println(message)
}
