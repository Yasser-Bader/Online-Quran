package models

import "gorm.io/gorm"

type Progres struct{
	gorm.Model
	ProgresID int
	StudentID int
	Date string
	Surah string
	Verses string
	Grade string
	Notes string
}