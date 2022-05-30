package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"fmt"
	"log"
)

type ProduceProdctRequest struct {
	ProductID int `json:"productID"`
	Quantity int `json: quantity`
	// key - material id, value - quantity 
	MaterialDemands map[int] int `json: materialDemands`
}



func ProduceProdctQuery (req ProduceProdctRequest) (string, error) {
	tx, err := dao.DB().Begin()
	if err != nil {
		return FailMessage, err
	}
	defer tx.Rollback()

	// generate eventID
	var eventID string
	err = tx.QueryRow("SELECT UUID()").Scan(&eventID)
	if err != nil {
		return FailMessage, err
	}

	var totalCost float32
	totalCost = 0
	// consume material
	for materialID, demandQuantity := range req.MaterialDemands {
		rows, err := dao.DB().Query("SELECT material_batch_id, vendor_id, sum(quantity), sum(cost) " +
								"FROM material_moving_record " +
								"WHERE material_id = ? " +
								"GROUP BY material_batch_id " +
								"HAVING sum(quantity) > 0 " +
								"ORDER BY date", materialID)
		defer rows.Close()
		if err != nil {
			return FailMessage, err
		}
		demandQuantity *= req.Quantity
		for rows.Next() {
			var materialBatchID, vendorID, quantity int
			var cost float32
			if err := rows.Scan(&materialBatchID, &vendorID, &quantity, &cost); err != nil {
				return FailMessage, err
			}
			reqQuantity := demandQuantity
			// cur batch is not enough
			if demandQuantity > quantity{
				reqQuantity = quantity
			}
			demandQuantity -= reqQuantity

			curBatchCost := cost * (float32(reqQuantity) / float32(quantity))

			_, err = tx.Exec("INSERT INTO material_moving_record(event_id, material_batch_id, material_id, vendor_id, quantity, cost, date) VALUES (?, ?, ?, ?, ?, ?, NOW())",
			eventID, materialBatchID, materialID, vendorID, -1 * reqQuantity, -1 * curBatchCost,)
			if err != nil {
				return FailMessage, err
			}
			totalCost += curBatchCost

			// meet demand
			if demandQuantity == 0 {
				break
			}
		}

		fmt.Printf("material %d, remain: %d \n", materialID, demandQuantity )
		if demandQuantity > 0 {
			//material demand not meet and roll back
			msg := fmt.Sprintf("Not enough material ID: %v", materialID)
			return msg, errors.New(msg)
		}
	}
	// get next batchID
	var productBatchID int
	err = tx.QueryRow("SELECT IFNULL(MAX(product_batch_id), 0) FROM inventory_moving_record").Scan(&productBatchID)
	if err != nil {
		return FailMessage, err
	}
	productBatchID += 1

	// produce product
	_, err = tx.Exec("INSERT INTO inventory_moving_record(product_batch_id, event_id, product_id, quantity, cost, date) VALUES (?, ?, ?, ?, ?, NOW())",
	productBatchID, eventID, req.ProductID, req.Quantity, totalCost)
	if err != nil {
		return FailMessage, err
	}

	if err = tx.Commit(); err != nil {
		return FailMessage, err
	}

	return "OK", nil
}

func ProduceProduct(context *gin.Context) {
	request := ProduceProdctRequest{}
	context.ShouldBind(&request)
	msg, err := ProduceProdctQuery(request)
	if err != nil {
		log.Println(err)
	}
	context.JSON(200, gin.H{
		"message": msg,
	})
}