package models

import "gorm.io/gorm"

type Students struct{
	gorm.Model
	StudentID int
	FirstName string
	MiddelName string
	LastName string
	Phone string
	Email string
	MagicLinkToken string
}

