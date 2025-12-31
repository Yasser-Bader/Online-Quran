package models

import (
	"time"
	
)
type Slots struct{
	ID             uint64 `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DayOfWeek string `gorm:"type:enum('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');default:'Sunday'"`
	StartTime string  `gorm:"type:time"`
	IsActive bool
	
}