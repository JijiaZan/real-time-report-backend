package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"fmt"
)

func main () {
	db, err := sql.Open("mysql", "root:12345678@tcp(139.196.30.199:3306)/graduate_project")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	rows, err2 := db.Query("SHOW TABLES")
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

