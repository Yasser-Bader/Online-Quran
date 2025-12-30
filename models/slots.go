package models

import "gorm.io/gorm"
type Slots struct{
	gorm.Model
	SlotID int
	DayOfWeek string
	StartTime string
	IsActive bool
	
}