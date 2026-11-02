package main

import (
	"github.com/gin-gonic/gin"
	handler "github.com/samsingh2h27/Logatharus/LogatharusService/Handler"
)

func RegisterRoutes() {
	// Register your routes here

	//logger class routes
	server := gin.Default()

	server.POST("/Logger/ValidateToken", handler.ValidateToken)
	server.GET("/Logger/Logs", handler.Log)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}