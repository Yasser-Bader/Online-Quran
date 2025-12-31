package models

import "time"

type Students struct {
		ID           uint64    `gorm:"primaryKey"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		FirstName    string    `form:"first_name"`   // أضفنا هذا
		MiddelName   string    `form:"middle_name"`  // أضفنا هذا (لاحظت أنك كتبتها Middel سأتركها كما هي لتجنب الأخطاء)
		LastName     string    `form:"last_name"`    // أضفنا هذا
		Phone        string    `gorm:"unique" form:"phone"` // أضفنا هذا
		Email        string    `gorm:"unique" form:"email"` // أضفنا هذا
		MagicLinkToken string
}
