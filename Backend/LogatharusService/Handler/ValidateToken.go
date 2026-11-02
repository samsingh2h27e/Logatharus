package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/gin-gonic/gin"
)


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
	userToken := requestBody.token
	_,err = DecryptToken(userToken,[]byte("your_secret_key_here"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid token"})
		return
	}
}

func DecryptToken(userToken string,key []byte) (string,error) {
	token, err := base64.RawStdEncoding.DecodeString(userToken)
	if err != nil {
		return "",err
	}
	block,err := aes.NewCipher(key)
	if err != nil {
		return "",err
	}
	gcm,err := cipher.NewGCM(block)
	if err != nil {
		return "",err
	}
	if len(token) < gcm.NonceSize()+gcm.Overhead() {
		return "",errors.New("invalid token length")
	}
	nonce := token[:gcm.NonceSize()]
	cipherText := token[gcm.NonceSize():]
	plainText,err := gcm.Open(nil,nonce,cipherText,nil)
	if err != nil {
		return "",err
	}
	return string(plainText),nil
}