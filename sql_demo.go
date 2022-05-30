package main

import (
	// _ "github.com/go-sql-driver/mysql"
	// "database/sql"
	// "github.com/JijiaZan/real-time-report-backend/dao"
	"github.com/JijiaZan/real-time-report-backend/service"
	// "log"
	// "fmt"
)

func main () {
	// db, err := sql.Open("mysql", "root:12345678@tcp(139.196.30.199:3306)/graduate_project")
	// defer db.Close()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// Printf("hasfsa")
	service.ListMaterialQuery()
	// _, err := dao.DB().Exec("INSERT INTO material_moving_recorder VALUES (UUID(), 1, 2, 3, 4, 5, NOW(), 3)")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err2 := dao.DB().Query("SELECT VERSION()")
	// defer rows.Close()
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	// for rows.Next() {
	// 	var table string
	// 	rows.Scan(&table)
	// 	fmt.Printf(table)
	// }
}

