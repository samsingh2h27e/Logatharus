package handler

import "github.com/gin-gonic/gin"


type ValidationRequest struct {
	token string `json:"token" binding:"required"`
}

func ValidateToken(ctx *gin.Context) {
	// Implement your token validation logic here
	var requestBody ValidationRequest
	err := ctx.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

}