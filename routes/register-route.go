package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/gin-gonic/gin"
	"Online-Quran/config"
	"Online-Quran/models"
)

// عرض صفحة التسجيل (للطلاب الجدد)
func Show_Register_Page(c *gin.Context) {
	var slots []models.Slots
	
	// التعديل هنا: نضيف شرط AND mode = 'assessment'
	// لكي لا تظهر مواعيد الدروس الخاصة أو المجموعات للطالب الجديد
	config.ConnectDB().Where("is_booked = ? AND mode = ?", false, "assessment").Find(&slots)

	c.HTML(http.StatusOK, "register.html", gin.H{
		"slots": slots,
	})
}

/*/ عرض الصفحة المدمجة
func Show_Register_Page(c *gin.Context) {
    // نرسل المواعيد المتاحة للصفحة
    var slots []models.Slots
    config.ConnectDB().Where("is_booked = ?", false).Find(&slots)
    c.HTML(http.StatusOK, "register.html", gin.H{"slots": slots})
}*/

// المعالجة المدمجة (Transaction)
func Register_And_Book(c *gin.Context) {
    db := config.ConnectDB()
    
    // 1. استقبال البيانات (يدوياً أو بالربط)
    firstName := c.PostForm("first_name")
    lastName := c.PostForm("last_name")
    phone := c.PostForm("phone")
    email := c.PostForm("email")
    slotID := c.PostForm("slot_id")

    // 2. معالجة الصورة قبل الدخول في قاعدة البيانات
    file, err := c.FormFile("receipt")
    if err != nil {
        c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "يجب رفع الصورة"})
        return
    }
    filename := fmt.Sprintf("%s_%s", phone, filepath.Base(file.Filename)) // نستخدم الهاتف مؤقتاً في الاسم

    // 3. بدء الـ Transaction (أهم خطوة)
    tx := db.Begin()

    // أ) حفظ الطالب
    student := models.Students{
        FirstName: firstName, 
        LastName: lastName, 
        Phone: phone, 
        Email: email,
        Level: "new",
    }
    
    if err := tx.Create(&student).Error; err != nil {
        tx.Rollback() // الغاء العملية
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "هذا الإيميل أو الهاتف مسجل مسبقاً، يرجى تسجيل الدخول"})
        return
    }

    // ب) حفظ الحجز
    slotUint, _ := strconv.ParseUint(slotID, 10, 64)
    booking := models.Booking{
        StudentID:    uint(student.ID), // نأخذ الـ ID من الطالب الذي تم إنشاؤه للتو
        SlotID:       uint(slotUint),
        PaymentImage: filename,
        Status:       "pending",
        BookingType:  "assessment",
        Amount:       50,
    }

    if err := tx.Create(&booking).Error; err != nil {
        tx.Rollback() // لو فشل الحجز، نحذف الطالب تلقائياً
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "فشل حفظ الحجز"})
        return
    }

    // ج) تحديث الموعد
    if err := tx.Model(&models.Slots{}).Where("id = ?", slotUint).Update("is_booked", true).Error; err != nil {
        tx.Rollback()
        return
    }

    // 4. اعتماد العملية نهائياً وحفظ الصورة
    tx.Commit()
    c.SaveUploadedFile(file, "./uploads/"+filename)

    c.HTML(http.StatusOK, "register.html", gin.H{"message": "تم التسجيل والحجز بنجاح! سيتم التواصل معك."})
}