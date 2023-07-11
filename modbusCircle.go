package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/modbus"
	"log"
	"net/http"
)

func main() {
	// Create a new Modbus RTU server
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = 1 // Set the timeout in seconds
	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	// Create a new Gin router
	router := gin.Default()

	// Define a route for reading a holding register
	router.GET("/holding-register/:address", func(c *gin.Context) {
		address := c.Param("address")
		client := modbus.NewClient(handler)
		result, err := client.ReadHoldingRegisters(0, 1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"address": address, "value": result[0]})
	})

	// Define a route for writing a holding register
	router.PUT("/holding-register/:address/:value", func(c *gin.Context) {
		address := c.Param("address")
		value := c.Param("value")
		client := modbus.NewClient(handler)
		// Convert value to int16 (assuming your controller uses 16-bit registers)
		convertedValue := int16(0)
		_, err := fmt.Sscanf(value, "%d", &convertedValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value"})
			return
		}
		_, err = client.WriteSingleRegister(0, uint16(convertedValue))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"address": address, "value": convertedValue})
	})

	// Start the HTTP server
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
