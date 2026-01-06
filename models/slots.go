package models

// --- جدول المواعيد ---
type Slots struct {
	ID       uint `gorm:"primaryKey"`
	Day      string
	Time     string
	IsBooked bool   `gorm:"default:false"`
	Mode     string // private, group, assessment
	ZoomLink string
}
