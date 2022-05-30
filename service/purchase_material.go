package service

import (
	"github.com/gin-gonic/gin"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"log"
)

type PurchaseMaterialRequest struct {
	MaterialID int `json:"materialID"`
	VendorID int `json: vendorID`
	Quantity int `json: quantity`
	Cost float32 `json: cost`
	AccountID int `json: accountID`
}

func PurchaseMaterialQuery(req PurchaseMaterialRequest) (string, error) {
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
	
	// get next batchID
	batchID := -1
	err = tx.QueryRow("SELECT IFNULL(MAX(material_batch_id), 0) FROM material_moving_record").Scan(&batchID)
	if err != nil {
		return FailMessage, err
	}
	batchID += 1
	
	// buy in material
	_, err = tx.Exec("INSERT INTO material_moving_record(event_id, material_batch_id, material_id, vendor_id, quantity, cost, date) VALUES (?, ?, ?, ?, ?, ?, NOW())",
						eventID, batchID, req.MaterialID, req.VendorID, req.Quantity, req.Cost,)
	if err != nil {
		return FailMessage, err
	}
	
	// cash flow out
	_, err = tx.Exec("INSERT INTO cash_flow(event_id,account_id,amount,date) VALUES (?,?,?,NOW())",
						eventID, req.AccountID, -1 * req.Cost)

	
	if err = tx.Commit(); err != nil {
		return FailMessage, err
	}
	return "OK", nil

}

func PurchaseMaterial(context *gin.Context) {
	request := PurchaseMaterialRequest{}
	context.ShouldBind(&request)
	msg, err := PurchaseMaterialQuery(request)
	if err != nil {
		log.Println(err)
	}
	context.JSON(200, gin.H{
		"message": msg,
	})
}