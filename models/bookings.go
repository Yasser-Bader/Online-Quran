package models

import (
	"time"

	"gorm.io/gorm"
)

// --- جدول الحجوزات ---
type Booking struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// الربط مع الطالب (مهم جداً لظهور الاسم)
	StudentID string
	Student   Students `gorm:"foreignKey:StudentID;references:ID"`

	// الربط مع الموعد
	SlotID uint
	Slot   Slots `gorm:"foreignKey:SlotID"`

	PaymentImage string
	Status       string `gorm:"default:'pending'"` // pending, confirmed
	BookingType  string // assessment, group_monthly, private_session
	Amount       float64
}
