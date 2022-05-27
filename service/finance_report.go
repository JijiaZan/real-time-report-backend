package service

import (
	"fmt"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/gin-gonic/gin"
)

//还需要查看哪些报表，可能需要进一步明确需求》》》
func FinanceReport(context *gin.Context) {

	db := dao.DB()
	//查看报表1
	sqlStr := "select * from material_moving_recorder"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query data failed，err:%s\n", err)
	}
	defer rows.Close()

	context.JSON(200, gin.H{
		"message": "OK",
		"data":    rows,
	})
}
