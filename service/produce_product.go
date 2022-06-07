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

type MaterialBatchInfo struct {
	MaterialBatchID int
	VendorID int
	Quantity int
	Cost float32
}



func ProduceProdctQuery (req ProduceProdctRequest) (error) {
	tx, err := dao.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// generate eventID
	var eventID string
	err = tx.QueryRow("SELECT UUID()").Scan(&eventID)
	if err != nil {
		return err
	}

	var totalCost float32
	totalCost = 0
	// consume material
	for materialID, demandQuantity := range req.MaterialDemands {
		rows, err := tx.Query("SELECT material_batch_id, vendor_id, sum(quantity), sum(cost) " +
								"FROM material_moving_record " +
								"WHERE material_id = ? " +
								"GROUP BY material_batch_id " +
								"HAVING sum(quantity) > 0 " +
								"ORDER BY date", materialID)
		defer rows.Close()
		if err != nil {
			return err
		}
		// total demand of one mat = per product demand * quantity of product
		demandQuantity *= req.Quantity
		materialBatchInfoList := []MaterialBatchInfo{}
		for rows.Next() {
			// var materialBatchID, vendorID, quantity int
			// var cost float32
			materialBatchInfo := MaterialBatchInfo{}
			if err := rows.Scan(&materialBatchInfo.MaterialBatchID, &materialBatchInfo.VendorID, &materialBatchInfo.Quantity, &materialBatchInfo.Cost); err != nil {
				return err
			}
			materialBatchInfoList = append(materialBatchInfoList, materialBatchInfo)
		}
		for _, materialBatchInfo := range(materialBatchInfoList) {
			reqQuantity := demandQuantity
			// cur batch is not enough
			if demandQuantity > materialBatchInfo.Quantity{
				reqQuantity = materialBatchInfo.Quantity
			}
			demandQuantity -= reqQuantity

			curBatchCost := materialBatchInfo.Cost / float32(materialBatchInfo.Quantity) * float32(reqQuantity)

			_, err = tx.Exec("INSERT INTO material_moving_record(event_id, material_batch_id, material_id, vendor_id, quantity, cost, date) VALUES (?, ?, ?, ?, ?, ?, NOW())",
			eventID, materialBatchInfo.MaterialBatchID, materialID, materialBatchInfo.VendorID, -1 * reqQuantity, -1 * curBatchCost,)
			if err != nil {
				return err
			}
			totalCost += curBatchCost
			// meet demand
			if demandQuantity == 0 {
				break
			}
		}

		// fmt.Printf("material %d, remain: %d \n", materialID, demandQuantity )
		if demandQuantity > 0 {
			//material demand not meet and roll back
			errMsg := fmt.Sprintf("Not enough material ID: %v", materialID)
			return errors.New(errMsg)
		}
	}
	// get next ProductBatchID
	var productBatchID int
	err = tx.QueryRow("SELECT IFNULL(MAX(product_batch_id), 0) FROM inventory_moving_record").Scan(&productBatchID)
	if err != nil {
		return err
	}
	productBatchID += 1

	// produce product
	_, err = tx.Exec("INSERT INTO inventory_moving_record(product_batch_id, event_id, product_id, quantity, cost, date) VALUES (?, ?, ?, ?, ?, NOW())",
	productBatchID, eventID, req.ProductID, req.Quantity, totalCost)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ProduceProduct(context *gin.Context) {
	request := ProduceProdctRequest{}
	context.ShouldBind(&request)
	err := ProduceProdctQuery(request)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		context.JSON(200, gin.H{
			"message": "Produce product succeed!",
		})
	}

}