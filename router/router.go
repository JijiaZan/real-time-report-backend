package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/JijiaZan/real-time-report-backend/service"
)

var router *gin.Engine

func init()  {
	router = gin.Default()
	router.POST("/purchase_material", service.PurchaseMaterial)
	router.POST("/produce_product", service.ProduceProduct)
	router.POST("/sell_product", service.SellProduct)
}

func Run(port uint16) {
	router.Run(fmt.Sprintf(":%d", port))
}