package routes

import (
	"Online-Quran/config"
	"Online-Quran/models"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// عرض صفحة التسجيل (الرئيسية)
func Show_Register_Page(c *gin.Context) {
	var slots []models.Slots
	// للطلاب الجدد: نظهر فقط مواعيد الـ assessment
	config.ConnectDB().Where("is_booked = ? AND mode = ?", false, "assessment").Find(&slots)
	c.HTML(http.StatusOK, "register.html", gin.H{"slots": slots})
}

// تسجيل وحجز في خطوة واحدة
func Register_And_Book(c *gin.Context) {
	db := config.ConnectDB()

	// استقبال البيانات
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	slotID := c.PostForm("slot_id")

	file, err := c.FormFile("receipt")
	if err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "يجب رفع إيصال التحويل"})
		return
	}
	filename := fmt.Sprintf("%s_%s", phone, filepath.Base(file.Filename))

	// Transaction لضمان سلامة البيانات
	tx := db.Begin()

	// 1. إنشاء الطالب
	student := models.Students{
		FirstName: firstName, LastName: lastName, Phone: phone, Email: email, Level: "new",
	}
	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusOK, "register.html", gin.H{"error": "البيانات مسجلة مسبقاً، يرجى تسجيل الدخول"})
		return
	}

	// 2. إنشاء الحجز
	slotUint, _ := strconv.ParseUint(slotID, 10, 64)
	booking := models.Booking{
		StudentID:    student.ID, // UUID
		SlotID:       uint(slotUint),
		PaymentImage: filename,
		Status:       "pending",
		BookingType:  "assessment",
		Amount:       50,
	}
	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "فشل في إنشاء الحجز"})
		return
	}

	// 3. حجز الموعد
	tx.Model(&models.Slots{}).Where("id = ?", slotUint).Update("is_booked", true)

	tx.Commit()
	c.SaveUploadedFile(file, "./uploads/"+filename)

	c.HTML(http.StatusOK, "register.html", gin.H{
		"success": true, "student_name": firstName,
	})
}

// صفحة تسجيل الدخول
func Show_Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// التحقق من الدخول
func Check_Login(c *gin.Context) {
	phone := c.PostForm("phone")
	var student models.Students

	if err := config.ConnectDB().Where("phone = ?", phone).First(&student).Error; err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "رقم الهاتف غير مسجل لدينا."})
		return
	}

	// التحقق هل لديه حجز معلق؟
	var lastBooking models.Booking
	config.ConnectDB().Where("student_id = ?", student.ID).Order("created_at desc").First(&lastBooking)

	if lastBooking.ID != 0 && lastBooking.Status == "pending" {
		c.HTML(http.StatusOK, "login.html", gin.H{"warning": "طلبك السابق ما زال قيد المراجعة."})
		return
	}

	// توجيه للتجديد
	c.Redirect(http.StatusFound, fmt.Sprintf("/booking?student_id=%s", student.ID))
}

// صفحة الحجز (للتجديد)
func Show_Booking(c *gin.Context) {
	studentID := c.Query("student_id")
	var student models.Students

	// استخدام UUID في البحث
	if err := config.ConnectDB().Where("id = ?", studentID).First(&student).Error; err != nil {
		c.HTML(http.StatusNotFound, "booking.html", gin.H{"error": "حدث خطأ، يرجى إعادة تسجيل الدخول"})
		return
	}

	// جلب المواعيد المتاحة
	var slots []models.Slots
	query := config.ConnectDB().Where("is_booked = ?", false)
	// لو الطالب جديد، نظهر له فقط الـ Assessment
	if student.Level == "new" {
		query = query.Where("mode = ?", "assessment")
	}
	query.Find(&slots)

	c.HTML(http.StatusOK, "booking.html", gin.H{"student": student, "slots": slots})
}

// إنشاء حجز التجديد
func Create_Booking(c *gin.Context) {
	studentID := c.PostForm("student_id")
	slotID := c.PostForm("slot_id")
	bookingType := c.PostForm("booking_type")
	amountStr := c.PostForm("amount")

	file, _ := c.FormFile("receipt")
	filename := fmt.Sprintf("%s_renew_%s", studentID, filepath.Base(file.Filename))
	c.SaveUploadedFile(file, "./uploads/"+filename)

	slotUint, _ := strconv.ParseUint(slotID, 10, 64)
	amount, _ := strconv.ParseFloat(amountStr, 64)

	booking := models.Booking{
		StudentID:    studentID,
		SlotID:       uint(slotUint),
		PaymentImage: filename,
		Status:       "pending",
		BookingType:  bookingType,
		Amount:       amount,
	}
	config.ConnectDB().Create(&booking)

	// تحديث الموعد
	config.ConnectDB().Model(&models.Slots{}).Where("id = ?", slotUint).Update("is_booked", true)

	// جلب بيانات الطالب لكي نرسل اسمه لصفحة النجاح
	var student models.Students
	config.ConnectDB().Where("id = ?", studentID).First(&student)

	c.HTML(http.StatusOK, "booking.html", gin.H{
		"success_message": "تم إرسال طلب التجديد بنجاح! سيتم مراجعته وإرسال التأكيد عبر الإيميل.",
		"student":         student, // 👈 إرسال الطالب الحقيقي بدلاً من الـ Placeholder
	})
}

// --- 5. صفحة الطالب الخاصة (لعرض الدرجات) ---
// --- 5. عرض ملف الطالب (باستخدام الرابط السري) ---
func Show_Student_Profile(c *gin.Context) {
	token := c.Param("token")
	var student models.Students

	// 1. البحث عن الطالب باستخدام التوكن
	// ملاحظة: الـ ID أصبح String (UUID) وهذا لا يؤثر هنا لأننا نبحث بالتوكن
	if err := config.ConnectDB().Where("magic_link_token = ?", token).First(&student).Error; err != nil {
		c.HTML(http.StatusNotFound, "profile.html", gin.H{"error": "الرابط غير صحيح أو منتهي الصلاحية"})
		return
	}

	// 2. جلب سجل الدرجات (Progress)
	// قمت بتصحيح اسم الجدول ليكون Progress (بدل Progres) ليكون Clean Code
	var progress []models.Progress
	config.ConnectDB().Where("student_id = ?", student.ID).Order("date desc").Find(&progress)

	// 3. عرض الصفحة
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"student":  student,
		"progress": progress,
	})
}
