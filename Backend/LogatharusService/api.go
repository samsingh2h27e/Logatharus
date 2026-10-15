package logatharusService

import (
	"github.com/gin-gonic/gin"
	handler "github.com/samsingh2h27/Logatharus/LogatharusService/Handler"
)

func RegisterRoutes() {
	// Register your routes here

	//logger class routes
	server := gin.Default()
	defer server.Run()

	server.GET("/Logger/ValidateToken",handler.ValidateToken)
	server.POST("/Logger/Logs",handler.Log)
}