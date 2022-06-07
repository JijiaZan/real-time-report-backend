package service

import (
	"errors"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/gin-gonic/gin"
	"log"
)

type SellProductRequest struct {
	ProductID int `json:"productID"`
	ClientID  int `json:"clientID"`
	Quantity  int `json: quantity`
	AccountID int `json: accountID`
}

// type ProductRecord struct {
// 	ID             int
// 	ProductBatchID int
// 	EventID        string
// 	ProductID      int
// 	Quantity       int
// 	Cost           float64
// 	Date           time.Time
// 	QuantitySold   int
// }

type ProductBatchInfo struct {
	ProductBatchID int
	Quantity       int
	Cost        float64
}

func SellProductQuery(req SellProductRequest) (error) {

	// Get a Tx for making transaction requests.
	tx, err := dao.DB().Begin()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var count int
	sqlStr := "select sum(quantity) as count from inventory_moving_record where product_id=?" //对所有product进行累计计算
	err = tx.QueryRow(sqlStr, req.ProductID).Scan(&count)
	if err != nil {
		return err
	}
	if count < req.Quantity {
		return errors.New("Not enough inventory")
	}

	var eventID string
	err = tx.QueryRow("SELECT UUID()").Scan(&eventID)
	if err != nil {
		return err
	}

	//找到目前未销售完毕的product_id,product_batch_id,sum(quantity) 分组组合
	sqlStr2 := "select product_batch_id, sum(quantity), sum(cost) from inventory_moving_record where product_id = ? group by product_batch_id having sum(quantity)>0 order by date"
	rows, err := tx.Query(sqlStr2, req.ProductID)
	defer rows.Close()
	if err != nil {
		return err
	}
	
	//sell by date order
	batchInfoList := []ProductBatchInfo{}
	totalCost := 0.0
	for rows.Next() { 
		var batchInfo ProductBatchInfo
		err := rows.Scan(&batchInfo.ProductBatchID, &batchInfo.Quantity, &batchInfo.Cost)
		if err != nil {
			return err
		}
		batchInfoList = append(batchInfoList, batchInfo)
	}

	demandQuantity := req.Quantity
	for _, batchInfo := range(batchInfoList) {
		sqlStr3 := "insert into inventory_moving_record(product_batch_id,event_id,product_id,quantity,cost,date) values(?,?,?,?,?,now())"
		if batchInfo.Quantity >= demandQuantity {
			// insert inventory opeation
			cost := batchInfo.Cost / float64(batchInfo.Quantity)  * float64(demandQuantity)
			totalCost += cost		
			_, err := tx.Exec(sqlStr3, batchInfo.ProductBatchID, eventID, req.ProductID, -demandQuantity, -cost)
			if err != nil {
				return err
			}
			break
		} else {
			// insert inventory opeation
			totalCost += batchInfo.Cost
			_, err := tx.Exec(sqlStr3, batchInfo.ProductBatchID, eventID, req.ProductID, -batchInfo.Quantity, -batchInfo.Cost)
			if err != nil {
				return err
			}
			//更新quantity的值
			demandQuantity -= batchInfo.Quantity
		}
	}

	//获取商品单价
	var price float64
	sqlStr3 := "select price from product_info where product_id = ?"
	err = tx.QueryRow(sqlStr3, req.ProductID).Scan(&price)
	if err != nil {
		return err
	}
	// insert cashflow opeation
	cashFlowPrice := float64(req.Quantity) * price
	sqlStr4 := "insert into cash_flow(event_id,account_id,amount,date) values(?,?,?,now())"
	_, err = tx.Exec(sqlStr4, eventID, req.AccountID, cashFlowPrice)
	if err != nil {
		return err
	}

	// insert sales order
	sqlStr5 := "insert into sales_order(event_id, product_id, client_id, sales_amount, cost, quantity, date) values(?,?,?,?,?,?,now())"
	_, err = tx.Exec(sqlStr5, eventID, req.ProductID, req.ClientID, cashFlowPrice, totalCost, req.Quantity)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func SellProduct(context *gin.Context) {
	request := SellProductRequest{}
	context.ShouldBind(&request)
	err := SellProductQuery(request)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		context.JSON(200, gin.H{
			"message": "Sell product succeed!",
		})
	}

}
