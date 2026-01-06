package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- جدول الطلاب ---
type Students struct {
	ID             string `gorm:"primaryKey"` // UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string `form:"first_name"`
	LastName       string `form:"last_name"`
	Phone          string `gorm:"unique" form:"phone"`
	Email          string `gorm:"unique" form:"email"`
	MagicLinkToken string
	Level          string `gorm:"default:'new'"` // new, level1, level2...
}

// توليد UUID تلقائياً قبل الإنشاء
func (s *Students) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
