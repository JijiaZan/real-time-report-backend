package service

import (
	"container/list"
	"database/sql"
	"fmt"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/google/uuid"
)

type SellProductRequest struct {
	ProductId int `json:"productId"`
	ClientId  int `json:"clientId"`
	Quantity  int `json: quantity`
	AccountId int `json: accountId`
}

type ProductRecord struct {
	Id             int
	ProductBatchId int
	EventId        string
	ProductId      int
	Quantity       int
	Cost           float64
	Date           time.Time
	QuantitySold   int
}

type ProductBatchInfo struct {
	ProductId      int
	ProductBatchId int
	Quantity       int
	EventId        string
}

func SellProductQuery(productId int, clientId int, quantity int, accountId int, db *sql.DB, context *gin.Context) (err error) {

	// Create a helper function for preparing failure results.
	fail := func(err error) error {
		return fmt.Errorf("sellProduct: %v", err)
	}

	// Get a Tx for making transaction requests.
	tx, err := db.BeginTx(context, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var count int
	sqlStr := "select sum(quantity) as count from inventory_moving_recorder where product_id=?" //对所有product进行累计计算
	err = tx.QueryRow(sqlStr, productId).Scan(&count)
	if err != nil {
		fmt.Printf("get quantity fail:%s\n", err)
		return err
	}
	if count < quantity {
		return fail(fmt.Errorf("not enough inventory"))
	}

	//找到目前未销售完毕的product_id,product_batch_id,sum(quantity) 分组组合
	sqlStr2 := "select product_id,product_batch_id,sum(quantity),event_id as quantity from inventory_moving_recorder where product_id = ? group by product_batch_id having sum(quantity)>0 order by product_batch_id"
	rows, err := tx.Query(sqlStr2, productId)
	if err != nil {
		fmt.Printf("query data failed，err:%s\n", err)
		return err
	}
	defer rows.Close()
	//获取所有未卖完的product信息
	productBatchInfo := list.New()
	for rows.Next() {
		var productBatchInfo_ ProductBatchInfo
		err := rows.Scan(&productBatchInfo_.ProductId, &productBatchInfo_.ProductBatchId, &productBatchInfo_.Quantity, &productBatchInfo_.EventId)
		if err != nil {
			return err
		}
		fmt.Print(productBatchInfo_)
		productBatchInfo.PushBack(productBatchInfo_)
	}

	//获取商品单价
	var price float64
	sqlStr3 := "select price from product_info where product_id = ?"
	err = tx.QueryRow(sqlStr3, productId).Scan(&price)
	if err != nil {
		fmt.Printf("get product price failed，err:%s\n", err)
		return err
	}

	uuid := uuid.New()
	eventId := uuid.String()
	//eventId := "222222" //获取一个全局唯一id
	cashFlowPrice := float64(quantity) * price

	//循环遍历满足要求的batch_id list，判断数量是否满足要求。若不满足要求，依次遍历后面的数据 product_id,batch_id,quantity
	for i := productBatchInfo.Front(); i != nil; i = i.Next() {
		if i == nil {
			fmt.Print("==================error===============")
		}
		ele := (i.Value).(ProductBatchInfo)
		quan := ele.Quantity

		// 获取当前product_id，batch_id对应的原材料成本
		var materialCost float64
		fmt.Print("!!!",ele.EventId)
		sqlStr := "select cost  from material_moving_recorder where event_id=?" //对所有product进行累计计算
		err = tx.QueryRow(sqlStr, ele.EventId).Scan(&materialCost)
		if err != nil {
			fmt.Printf("get cost failed，err:%s\n", err)
			return err
		}

		//循环插入数据
		sqlStr3 := "insert into inventory_moving_recorder(product_batch_id,event_id,product_id,quantity,cost,date) values(?,?,?,?,?,?)"
		golangDateTime := time.Now().Format("2006-01-02 15:04:05")
		if quan >= quantity {
			// insert inventory opeation
			cost := float64(quantity) * materialCost
			_, err := tx.Exec(sqlStr3, ele.ProductBatchId, eventId, ele.ProductId, -quantity, -cost, golangDateTime)
			if err != nil {
				fmt.Printf("insert data failed, err:%v\n", err)
				return err
			}
			break
		} else {
			// insert inventory opeation
			cost := float64(quan) * materialCost
			_, err := tx.Exec(sqlStr3, ele.ProductBatchId, eventId, ele.ProductId, -quan, -cost, golangDateTime)
			if err != nil {
				fmt.Printf("insert data failed, err:%v\n", err)
				return err
			}
			//更新quantity的值
			quantity -= quan
		}
	}

	// insert cashflow opeation
	sqlStr4 := "insert into cash_flow(event_id,account_id,amount,date) values(?,?,?,?)"
	golangDateTime := time.Now().Format("2006-01-02 15:04:05")
	_, err = tx.Exec(sqlStr4, eventId, accountId, cashFlowPrice, golangDateTime)
	if err != nil {
		fmt.Printf("insert data failed, err:%v\n", err)
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return
}

func SellProduct(context *gin.Context) {
	request := SellProductRequest{}
	context.ShouldBind(&request)
	fmt.Print(context)
	db := dao.DB()
	err := SellProductQuery(request.ProductId, request.ClientId, request.Quantity, request.AccountId, db, context)
	if err != nil {
		fmt.Print("==================error2===============", err)
		context.JSON(500, gin.H{
			"message": err,
		})
	} else {
		context.JSON(200, gin.H{
			"message": "Sell product succeed!",
		})
	}

}
