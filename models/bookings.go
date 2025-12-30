package models

import "gorm.io/gorm"

type Bookings struct{
	gorm.Model
	BookingID int
	StudentId int
	SlotID int 
	BookingDate string
	PaymentScreenshot string
	Status string
}