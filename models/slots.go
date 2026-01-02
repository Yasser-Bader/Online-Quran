package models

type Slots struct {
	ID        uint   `gorm:"primaryKey"`
	Day       string // اليوم (السبت، الأحد...)
	Time      string // الساعة (06:00 م)
	IsBooked  bool   `gorm:"default:false"` // هل تم حجزه؟
}