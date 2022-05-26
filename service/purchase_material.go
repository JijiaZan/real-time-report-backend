package service

import (
	"github.com/gin-gonic/gin"
)

type PurchaseMaterialRequest struct {
	PurchaseOrderId string `json:"purchaseOrderId"`
	MaterialId string `json:"materialId"`
	VendorId string `json: vendorId`
	Quantity string `json: quantity`
	Value string `json: value`
}

func PurchaseMaterial(context *gin.Context) {
	request := PurchaseMaterialRequest{}
	context.ShouldBind(&request)


	context.JSON(200, gin.H{
		"purchaseOrderId": request.PurchaseOrderId,
	})
}