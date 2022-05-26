package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var db_connection *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:12345678@tcp(139.196.30.199:3306)/graduate_project")
	if err != nil {
        panic(err.Error())
    }
	db_connection = db
}

func CloesDB() {
	db_connection.Close()
}

// Return DB connection
func DB() *sql.DB {
	return db_connection
}