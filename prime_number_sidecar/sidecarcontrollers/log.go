package sidecarcontrollers

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func LogLn(c *gin.Context) {
	var message string
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &message)

	log.Println(message)
}
