package routes

import (
	"Online-Quran/config"
	"Online-Quran/models"
	"Online-Quran/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// عرض لوحة التحكم
func Admin_Dashboard(c *gin.Context) {
	var pendingBookings []models.Booking
	var confirmedBookings []models.Booking
	var students []models.Students
	var availableSlots []models.Slots // <-- متغير جديد للمواعيد المتاحة

	db := config.ConnectDB()

	// 1. الطلبات المعلقة
	db.Preload("Student").Preload("Slot").Where("status = ?", "pending").Find(&pendingBookings)

	// 2. الطلبات المؤكدة
	db.Preload("Student").Preload("Slot").Where("status = ?", "confirmed").Find(&confirmedBookings)

	// 3. قائمة الطلاب
	db.Find(&students)

	// 4. المواعيد المتاحة (التي لم تحجز بعد) <-- الإضافة الجديدة
	db.Where("is_booked = ?", false).Find(&availableSlots)

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"bookings":           pendingBookings,
		"confirmed_bookings": confirmedBookings,
		"students":           students,
		"slots":              availableSlots, // نرسلها للصفحة
	})
}

// --- دالة تعديل الموعد (Update) ---
func Admin_UpdateSlot(c *gin.Context) {
	id := c.Param("id")
	var slot models.Slots

	// البحث عن الموعد
	if err := config.ConnectDB().First(&slot, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}

	// تحديث البيانات
	slot.Day = c.PostForm("day")
	slot.Time = c.PostForm("time")
	slot.Mode = c.PostForm("mode")
	slot.ZoomLink = c.PostForm("zoom_link")

	// حفظ التغييرات
	config.ConnectDB().Save(&slot)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// --- دالة حذف الموعد (جديدة) ---
func Admin_DeleteSlot(c *gin.Context) {
	slotID := c.Param("id")
	// حذف الموعد من قاعدة البيانات
	config.ConnectDB().Delete(&models.Slots{}, slotID)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// الموافقة على الحجز
func Admin_Approve(c *gin.Context) {
	bookingID := c.Param("id")
	var booking models.Booking
	db := config.ConnectDB()

	// البحث مع Preload مهم جداً للإيميل
	if err := db.Preload("Student").Preload("Slot").First(&booking, bookingID).Error; err != nil {
		c.JSON(404, gin.H{"error": "الحجز غير موجود"})
		return
	}

	// تحديث الحالة
	booking.Status = "confirmed"
	db.Save(&booking)

	// إنشاء توكن للطالب لو مش موجود
	if booking.Student.MagicLinkToken == "" {
		booking.Student.MagicLinkToken = uuid.New().String()
		db.Save(&booking.Student)
	}

	// إرسال الإيميل (في الخلفية)
	go utils.SendConfirmationEmail(
		booking.Student.Email,
		booking.Student.FirstName,
		booking.Slot.ZoomLink,
		booking.Student.MagicLinkToken,
		booking.Slot.Mode,
	)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// رفض الحجز وحذفه
func Admin_Reject(c *gin.Context) {
	bookingID := c.Param("id")
	reason := c.PostForm("reason")
	var booking models.Booking
	db := config.ConnectDB()

	if err := db.Preload("Student").Preload("Slot").First(&booking, bookingID).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}

	// إرسال إيميل الرفض
	go utils.SendRejectionEmail(booking.Student.Email, booking.Student.FirstName, reason)

	// تحرير الموعد
	db.Model(&models.Slots{}).Where("id = ?", booking.SlotID).Update("is_booked", false)

	// حذف الحجز (Soft Delete أو Hard Delete حسب الرغبة، هنا Hard Delete للتنظيف)
	db.Delete(&booking)

	// ملحوظة: لا نحذف الطالب هنا، ربما يريد المحاولة مرة أخرى بصورة صحيحة
	// إلا إذا كان التسجيل كله وهمي، يمكنك حذف الطالب أيضاً:
	// db.Delete(&booking.Student)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// إضافة موعد جديد
func Admin_AddSlot(c *gin.Context) {
	slot := models.Slots{
		Day:      c.PostForm("day"),
		Time:     c.PostForm("time"),
		IsBooked: false,
		Mode:     c.PostForm("mode"),
		ZoomLink: c.PostForm("zoom_link"),
	}
	config.ConnectDB().Create(&slot)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// تحديث مستوى الطالب
func Admin_UpdateStudentLevel(c *gin.Context) {
	studentID := c.PostForm("student_id")
	newLevel := c.PostForm("level")
	config.ConnectDB().Model(&models.Students{}).Where("id = ?", studentID).Update("level", newLevel)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// إضافة درجات
func Admin_AddGrade(c *gin.Context) {
	config.ConnectDB().Create(&models.Progress{
		StudentID: c.PostForm("student_id"),
		Date:      time.Now(), // أو خذ التاريخ من الفورم
		Surah:     c.PostForm("surah"),
		Verses:    c.PostForm("verses"),
		Grade:     c.PostForm("grade"),
		Notes:     c.PostForm("notes"),
	})
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func Admin_DeleteBooking(c *gin.Context) {
	bookingID := c.Param("id")

	// الحذف هنا سيقوم بإخفاء الحجز من الاستعلامات العادية (Soft Delete)
	// لأننا نستخدم gorm.DeletedAt في المودل
	config.ConnectDB().Delete(&models.Booking{}, bookingID)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}
