package db

import (
	"gin-test/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Use PostgreSQL in gorm
)

var (
	db  *gorm.DB
	err error
)

// Init is initialize db from main function
func Init() {
	db, err = gorm.Open("postgres", "host=db port=5432 user=gin-test dbname=gin-test password=gin-test sslmode=disable")
	if err != nil {
		panic("データベースが開けません（Init）")
	}
	autoMigration()
	// timecard := models.Timecard{
	// 	Day:       "2020/10/01",
	// 	Start:     "10:00",
	// 	End:       "20:00",
	// 	BreakTime: "1:30",
	// }
	// db.Create(&timecard)
}

// GetDB is called in models
func GetDB() *gorm.DB {
	return db
}

// Close is closing db
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func autoMigration() {
	db.AutoMigrate(&models.Timecard{})
	db.AutoMigrate(&models.User{})
}
