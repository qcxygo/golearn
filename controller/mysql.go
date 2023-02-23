package controller

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	d, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/douyin"+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}
	db = d
}
