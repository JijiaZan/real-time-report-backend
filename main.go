package main
import (
	"github.com/JijiaZan/real-time-report-backend/router"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func main() {
	db, err := sql.Open("mysql", "root:12345678@tcp(139.196.30.199:3306)/graduate_project")
	if err != nil {
        panic(err.Error())
    }
	defer db.Close()
	router.Run(8080)
}