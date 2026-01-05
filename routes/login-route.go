package routes

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"Online-Quran/config"
	"Online-Quran/models"
)



func Show_Login(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", nil)
}

// التحقق من الدخول
func Check_Login(c *gin.Context) {
	phone := c.PostForm("phone")
	var student models.Students
	
	// محاولة البحث عن الطالب
	result := config.ConnectDB().Where("phone = ?", phone).First(&student)

	// إذا حدث خطأ (يعني الطالب غير موجود)
	if result.Error != nil {
		// نعيد عرض صفحة الدخول مع رسالة خطأ حمراء
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "عفواً، رقم الهاتف هذا غير مسجل لدينا. يرجى التسجيل كطالب جديد.",
		})
		return
	}

	// إذا وجد الطالب، نوجهه لصفحة الحجز الخاصة به
	c.Redirect(http.StatusFound, fmt.Sprintf("/booking?student_id=%d", student.ID))
}
/*
func Check_Login(c *gin.Context) {
    phone := c.PostForm("phone")
    var student models.Students
    
    // البحث عن الطالب بالهاتف
    if err := config.ConnectDB().Where("phone = ?", phone).First(&student).Error; err != nil {
        c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "رقم الهاتف غير مسجل، يرجى التسجيل كطالب جديد"})
        return
    }

    // توجيهه لصفحة الحجز الخاصة به (التي صممناها سابقاً)
    c.Redirect(http.StatusFound, fmt.Sprintf("/booking?student_id=%d", student.ID))
}*/