package main

import (
	"gin-test/db"
	"gin-test/router"

	_ "github.com/jinzhu/gorm/dialects/postgres" // Use PostgreSQL in gorm
)

// DB設定
func Initialize() {
}

func main() {
	db.Init()
	router.Init()
	defer db.Close()

}
