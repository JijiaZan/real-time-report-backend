package router

import (
	"fmt"
	"github.com/JijiaZan/real-time-report-backend/service"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.POST("/purchase_material", service.PurchaseMaterial)
	router.POST("/produce_product", service.ProduceProduct)
	router.POST("/sell_product", service.SellProduct)
	router.POST("/list_material_and_product", service.ListMaterialAndProduct)
	router.POST("/search_account", service.SearchAccount)
	router.POST("/finance_report", service.FinanceReport)
}

func Run(port uint16) {
	router.Run(fmt.Sprintf(":%d", port))
}
