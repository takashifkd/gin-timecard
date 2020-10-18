package controller

import (
	"gin-test/crypto"
	"gin-test/db"
	"gin-test/models"

	_ "github.com/jinzhu/gorm/dialects/postgres" // Use PostgreSQL in gorm
)

// ユーザー登録処理
func createUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	// db := gormConnect()
	db := db.GetDB()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil

}

// ユーザーを一件取得
func getUser(username string) models.User {
	db := db.GetDB()
	var user models.User
	db.First(&user, "username = ?", username)
	db.Close()
	return user
}

// 新規作成
func createTimecard(m models.Timecard) {
	db := db.GetDB()
	db.Create(&m)
}

//timecard一覧（ユーザー、月指定）取得
func getTimecardList(UserID string, month string) []models.Timecard {
	db := db.GetDB()
	var timecards []models.Timecard
	db.Where("UserID = ? AND Day LIKE ?", UserID, month+"%").Order("Day").Find(&timecards)
	db.Close()
	return timecards

}

//編集用timecard取得
func getTimecard(UserID string, id string) models.Timecard {
	db := db.GetDB()
	var timecard models.Timecard
	db.Where("UserID = ? AND id = ?", UserID, id).First(&timecard)
	db.Close()
	return timecard
}

func updateTimecard(timecard models.Timecard) {
	db := db.GetDB()
	db.Where("UserID = ? AND id = ?", timecard.UserID, timecard.ID).Updates(&timecard)
	db.Close()

}
