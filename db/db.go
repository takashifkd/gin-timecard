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
	// db, err = gorm.Open("postgres", "host=db port=5432 user=gin-test dbname=gin-test password=gin-test sslmode=disable")//ローカル
	db, err = gorm.Open("postgres", "host=ec2-18-213-176-229.compute-1.amazonaws.com port=5432 user=qgwmxbbrwbumhl dbname=d6341p2hhcno6h password=3c938c17cb93c27eae22c9e151062231057c5bf04ff53e7ff04ca87046ca61f9 sslmode=require") //本番
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
