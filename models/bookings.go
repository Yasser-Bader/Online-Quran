package models

/*import "time"

type Bookings struct{
	ID             uint64 `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	BookingID int
	StudentId int
	SlotID int 
	BookingDate string
	PaymentScreenshot string
	Status string
}*/
import (
	"time"
	"gorm.io/gorm"
)

type Bookings struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	StudentID uint           // مفتاح أجنبي للطالب
	Student   Students       `gorm:"foreignKey:StudentID"` // العلاقة
	PaymentImage string      // مسار الصورة المحفوظة
	Status       string      `gorm:"default:'pending'"` // pending, confirmed
	ZoomLink     string      // سنخزنه هنا أو نرسله من الإعدادات
}