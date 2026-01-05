package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"Online-Quran/config"
	"Online-Quran/models"
	"Online-Quran/utils" // استيراد ملف الإيميل

)


// --- دالة إضافة موعد جديد (معدلة) ---
func Admin_AddSlot(c *gin.Context) {
	day := c.PostForm("day")
	timeStr := c.PostForm("time")
	mode := c.PostForm("mode")         // private OR group
	zoomLink := c.PostForm("zoom_link") // رابط الزوم لهذا الموعد

	slot := models.Slots{
		Day:      day, 
		Time:     timeStr, 
		IsBooked: false,
		Mode:     mode,
		ZoomLink: zoomLink,
	}
	config.ConnectDB().Create(&slot)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// --- دالة جديدة: تحديث مستوى الطالب ---
func Admin_UpdateStudentLevel(c *gin.Context) {
	studentID := c.PostForm("student_id")
	newLevel := c.PostForm("level") // level1, level2, ...

	config.ConnectDB().Model(&models.Students{}).Where("id = ?", studentID).Update("level", newLevel)
	
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func Admin_Approve(c *gin.Context) {
    bookingID := c.Param("id")
    var booking models.Booking
    
    db := config.ConnectDB()

    // 1. مهم جداً: Preload("Student") AND Preload("Slot")
    // نحتاج Slot عشان نجيب رابط الزوم الخاص بهذا الموعد
    if err := db.Preload("Student").Preload("Slot").First(&booking, bookingID).Error; err != nil {
        c.JSON(404, gin.H{"error": "الحجز غير موجود"})
        return
    }

    booking.Status = "confirmed"
    db.Save(&booking)

    if booking.Student.MagicLinkToken == "" {
        booking.Student.MagicLinkToken = uuid.New().String()
        db.Save(&booking.Student)
    }

    // 2. نأخذ الرابط من الموعد نفسه (سواء كان خاص، مجموعة، أو تحديد مستوى)
    zoomLink := booking.Slot.ZoomLink 
    if zoomLink == "" {
        zoomLink = "https://zoom.us/default-link" // رابط احتياطي
    }

    // 3. الإرسال
    utils.SendConfirmationEmail(
        booking.Student.Email,
        booking.Student.FirstName,
        zoomLink,
        booking.Student.MagicLinkToken,
    )

    c.Redirect(http.StatusFound, "/admin/dashboard")
}