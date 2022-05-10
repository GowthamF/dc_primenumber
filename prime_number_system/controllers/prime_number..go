package controllers

import (
	"net/http"

	"dc_assignment.com/prime_number/v2/models"
	"dc_assignment.com/prime_number/v2/services"
	"dc_assignment.com/prime_number/v2/sidecar"
	"github.com/gin-gonic/gin"
)

func CheckPrimeNumber(c *gin.Context) {
	var primeNumber *models.PrimeNumbers
	err := c.BindJSON(&primeNumber)

	if err != nil {
		c.AbortWithError(500, err)
	}

	message := services.CheckIfPrimeNumber(primeNumber.NumberToCheck, primeNumber.StartRange, primeNumber.EndRange)
	go sidecar.Log("SENT TO ACCEPTOR")
	c.JSON(http.StatusOK, message)
}
