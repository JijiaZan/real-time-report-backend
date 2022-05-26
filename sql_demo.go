package main

import (
	// _ "github.com/go-sql-driver/mysql"
	// "database/sql"
	"github.com/JijiaZan/real-time-report-backend/dao"
	"log"
	"fmt"
)

func main () {
	// db, err := sql.Open("mysql", "root:12345678@tcp(139.196.30.199:3306)/graduate_project")
	// defer db.Close()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// Printf("hasfsa")
	
	rows, err2 := dao.DB().Query("SELECT VERSION()")
	defer rows.Close()
	if err2 != nil {
		log.Fatal(err2)
	}

	for rows.Next() {
		var table string
		rows.Scan(&table)
		fmt.Printf(table)
	}
}

