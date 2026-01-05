package models

type Slots struct {
	ID        uint   `gorm:"primaryKey"`
	Day       string // اليوم (السبت، الأحد...)
	Time      string // الساعة (06:00 م)
	IsBooked  bool   `gorm:"default:false"` // هل تم حجزه؟
	// الإضافات الجديدة
	Mode     string `gorm:"default:'private'"` // 'private' (خاص) أو 'group' (مجموعة)
	ZoomLink string // رابط الزوم الخاص بهذا الموعد (يحدده الشيخ عند الإنشاء)
}