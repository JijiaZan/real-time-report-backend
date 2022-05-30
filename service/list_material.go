package service

import (
	"github.com/gin-gonic/gin"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"log"
	// "fmt"
)

type Material struct {
	MaterialBatchID int `json:"materialID"`
	MaterialID int `json:"materialID"`
	VendorID int `json: vendorID`
	Quantity int `json: quantity`
	Cost float32 `json: cost`
}

var FailMessage = "Something wrong with DB, check log"

// specify at most one materialID
func ListMaterialQuery() ([]Material, string, error) {
	rows, err := dao.DB().Query("SELECT material_batch_id, material_id, vendor_id, sum(quantity), sum(cost) " +
								"FROM material_moving_record " +
								"GROUP BY material_batch_id " +
								"HAVING sum(quantity) > 0 " +
								"ORDER BY date")
	defer rows.Close()
	if err != nil {
		return nil, FailMessage, err
	}

	curMaterialList := make([]Material, 0)
	for rows.Next() {
		var materialBatchID, materialID, vendorID, quantity int
		var cost float32
		if err := rows.Scan(&materialBatchID, &materialID, &vendorID, &quantity, &cost); err != nil {
			return nil, FailMessage, err
		}
		curMaterialList = append(curMaterialList, Material{materialBatchID, materialID, vendorID, quantity, cost})
	}
	return curMaterialList, "OK", nil
}

func ListMaterial(context *gin.Context) {
	curMaterialList, msg, err := ListMaterialQuery()
	if err != nil {
		log.Println(err)
	}
	context.JSON(200, gin.H{
		"message": msg,
		"material": curMaterialList,
	})
}
