package service

import (
	"database/sql"
	"fmt"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/gin-gonic/gin"
)

type SearchAccountRequest struct {
	AccountId int `json:"accountId"`
}

type Account struct {
	AccountId     int
	AccountNumber string
	BankName      string
	AccountRemain float64
}

func SearchAccountQuery(accountId int, db *sql.DB) (account *Account, err error) {
	var acc Account
	sqlStr := "select b.account_id,b.account_number,b.bank_name,sum(a.amount) from cash_flow a , bank_account b where a.account_id=b.account_id and a.account_id=? group by a.account_id"

	fmt.Printf(string(accountId))
	row := db.QueryRow(sqlStr, accountId)
	err = row.Scan(&acc.AccountId, &acc.AccountNumber, &acc.BankName, &acc.AccountRemain)
	if err != nil {
		return nil, err
	}

	account = &acc
	return
}

func SearchAccount(context *gin.Context) {
	request := SearchAccountRequest{}
	context.ShouldBind(&request)

	accountId := request.AccountId
	db := dao.DB()
	defer db.Close()

	account, err := SearchAccountQuery(accountId, db)
	if err != nil {
		fmt.Printf("err:%s", err)
	}

	context.JSON(200, gin.H{
		"account": account,
	})
}
