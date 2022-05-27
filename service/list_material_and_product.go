package service

import (
	"github.com/gin-gonic/gin"
)

type ListMaterialAndProductRequest struct {
}

func ListMaterialAndProduct(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "OK",
	})
}
