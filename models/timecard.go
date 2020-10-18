package models

import "gorm.io/gorm"

type Timecard struct {
	gorm.Model
	UserID uint `gorm:"primaryKey"`
	// Year      string
	// Month     string
	Day       string `gorm:"primaryKey"`
	Start     string `form:"Start"`
	End       string `form:"End"`
	BreakTime string `form:"BreakTime"`
}
