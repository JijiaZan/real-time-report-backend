package router

import (
	"fmt"
	"github.com/JijiaZan/real-time-report-backend/service"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.POST("/purchase_material", service.PurchaseMaterial)
	router.POST("/produce_product", service.ProduceProduct)
	router.POST("/sell_product", service.SellProduct)
	router.GET("/list_material", service.ListMaterial)
	router.GET("/list_product", service.ListProduct)
	router.POST("/search_account", service.SearchAccount)
	router.POST("/finance_report", service.FinanceReport)
}

func Run(port uint16) {
	router.Run(fmt.Sprintf(":%d", port))
	defer dao.DB().Close()
}
