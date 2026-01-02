package models

import( 
	"time"
)

type Progres struct{
	ID             uint64 `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
    StudentID uint
	Date time.Time
	Surah string
	Verses string
	Grade string
	Notes string
}