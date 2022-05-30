package service

import (
	"github.com/gin-gonic/gin"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"log"
	// "fmt"
)

type Product struct {
	ProductBatchID int `json:"productBatchID"`
	ProductID int `json:"productID"`
	Quantity int `json: quantity`
	Cost float32 `json: cost`
}

func ListProductQuery() ([]Product, string, error) {
	rows, err := dao.DB().Query("SELECT product_batch_id, product_id, sum(quantity), sum(cost) " +
								"FROM inventory_moving_record " +
								"GROUP BY product_batch_id " +
								"HAVING sum(quantity) > 0 " +
								"ORDER BY date")
	defer rows.Close()
	if err != nil {
		return nil, FailMessage, err
	}

	curProductList := make([]Product, 0)
	for rows.Next() {
		var productBatchID, productID, quantity int
		var cost float32
		if err := rows.Scan(&productBatchID, &productID, &quantity, &cost); err != nil {
			return nil, FailMessage, err
		}
		curProductList = append(curProductList, Product{productBatchID, productID, quantity, cost})
	}
	
	return curProductList, "OK", err
}

func ListProduct(context *gin.Context) {
	var curProductList []Product
	curProductList, msg, err := ListProductQuery()
	if err != nil {
		log.Println(err)
	}
	context.JSON(200, gin.H{
		"message": msg,
		"product": curProductList,
	})
}