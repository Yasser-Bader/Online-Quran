/*package routes

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"Online-Quran/config"
	"Online-Quran/models"
)

	func Create_students(c *gin.Context) {
		var students models.Students
		if err := c.ShouldBindJSON(&students); err != nil {
		  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
		}
		config.ConnectDB().Create(&students)
		c.JSON(200, gin.H{
		  "view": students,
		})
	  }

	  func Show_students(c *gin.Context) {
				var students []models.Students
				config.ConnectDB().Find(&students)
				c.JSON(200, gin.H{
					"view": students,
					})
		}*/
package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"Online-Quran/config"
	"Online-Quran/models"
	"Online-Quran/utils" // استيراد ملف الإيميل
)
// --- صفحة الحجز (للطالب) ---


func Create_students(c *gin.Context) {
	var students models.Students

	// استخدام ShouldBind بدلاً من ShouldBindJSON لكي يقبل البيانات من الفورم
	if err := c.ShouldBind(&students); err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "تأكد من صحة البيانات: " + err.Error()})
		return
		}

		// حفظ الطالب في قاعدة البيانات
		result := config.ConnectDB().Create(&students)
		if result.Error != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "حدث خطأ أثناء الحفظ (ربما الإيميل أو الهاتف مكرر)"})
			return
		}

		// العودة للصفحة مع رسالة نجاح
		/*c.HTML(http.StatusOK, "index.html", gin.H{
			"message":"مرحباً " + students.FirstName + "\n تم تسجيل بياناتك بنجاح ! ",
		})*/
		// بدلاً من عرض رسالة، وجهه لصفحة الحجز
       // c.Redirect(http.StatusFound, fmt.Sprintf("/booking?student_id=%d", students.ID))
	    redirectPath := fmt.Sprintf("/booking?student_id=%d", students.ID)
	    c.Redirect(http.StatusFound, redirectPath)
}
	

// دالة عرض صفحة الفورم لأول مرة
func Show_Form(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func Show_students(c *gin.Context) {
	var students []models.Students
	config.ConnectDB().Find(&students)
	c.JSON(200, gin.H{
		"view": students,
	})
}

// عرض صفحة رفع الإيصال (تستقبل ID الطالب)
func Show_Booking(c *gin.Context) {
	studentID := c.Query("student_id")
	c.HTML(http.StatusOK, "booking.html", gin.H{"student_id": studentID})
}

// استقبال الصورة وحفظ الحجز
func Create_Booking(c *gin.Context) {
	studentID := c.PostForm("student_id")
	
	// استقبال الملف
	file, err := c.FormFile("receipt")
	if err != nil {
		c.HTML(http.StatusBadRequest, "booking.html", gin.H{"error": "يجب رفع صورة الإيصال"})
		return
	}

	// حفظ الملف في مجلد uploads
	filename := fmt.Sprintf("%s_%s", studentID, filepath.Base(file.Filename))
	dst := "./uploads/" + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.HTML(http.StatusInternalServerError, "booking.html", gin.H{"error": "فشل حفظ الصورة"})
		return
	}

	// حفظ في قاعدة البيانات
	idUint, _ := strconv.ParseUint(studentID, 10, 64)
	booking := models.Bookings{
		StudentID:    uint(idUint),
		PaymentImage: filename,
		Status:       "pending",
	}
	config.ConnectDB().Create(&booking)

	c.HTML(http.StatusOK, "booking.html", gin.H{"message": "تم إرسال الإيصال بنجاح! سيتم مراجعته وإرسال الرابط لك قريباً."})
}

// --- لوحة الأدمن (لك وللشيخ) ---

// عرض كل الحجوزات المعلقة
func Admin_Dashboard(c *gin.Context) {
	var bookings []models.Bookings
	// هات الحجوزات واعمل Join مع جدول الطلاب عشان نعرف اسم الطالب
	config.ConnectDB().Preload("Student").Where("status = ?", "pending").Find(&bookings)
	
	c.HTML(http.StatusOK, "admin.html", gin.H{"bookings": bookings})
}

// زر الموافقة
func Admin_Approve(c *gin.Context) {
	bookingID := c.Param("id")
	var booking models.Bookings
	
	db := config.ConnectDB()
	// نجيب الحجز ونحمل بيانات الطالب المرتبط بيه
	if err := db.Preload("Student").First(&booking, bookingID).Error; err != nil {
		c.JSON(404, gin.H{"error": "الحجز غير موجود"})
		return
	}

	// 1. تحديث الحالة
	booking.Status = "confirmed"
	db.Save(&booking)

	// 2. إنشاء Magic Token للطالب لو مش عنده
	if booking.Student.MagicLinkToken == "" {
		booking.Student.MagicLinkToken = uuid.New().String()
		db.Save(&booking.Student)
	}

	// 3. إرسال الإيميل (زوم + تقويم + توكن)
	zoomLink := "https://zoom.us/j/123456789" // رابط الشيخ الثابت
	go utils.SendConfirmationEmail(booking.Student.Email, booking.Student.FirstName, zoomLink, booking.Student.MagicLinkToken)

	// إعادة توجيه للأدمن
	c.Redirect(http.StatusFound, "/admin/dashboard")
}