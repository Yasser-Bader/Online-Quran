package models

import (
	"time"
)

// --- جدول الدرجات ---
type Progress struct {
	ID        uint `gorm:"primaryKey"`
	StudentID string
	Date      time.Time
	Surah     string
	Verses    string
	Grade     string
	Notes     string
}
