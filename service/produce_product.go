package service

import (
	"github.com/gin-gonic/gin"
)

func ProduceProduct(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "OK",
	})
}