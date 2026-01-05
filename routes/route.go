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
	"time"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/gin-gonic/gin"
	"Online-Quran/config"
	"Online-Quran/models"
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

//////////////
/* -- ملغي - 1. دوال المواعيد (الشيخ يضيف موعد) ---
func Admin_AddSlot(c *gin.Context) {
	day := c.PostForm("day")
	timeStr := c.PostForm("time")

	slot := models.Slots{Day: day, Time: timeStr, IsBooked: false}
	config.ConnectDB().Create(&slot)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}
*/
// --- 2. تحديث صفحة الحجز (لعرض المواعيد للطالب) ---
/*func Show_Booking(c *gin.Context) {
	studentID := c.Query("student_id")
	
	// جلب المواعيد المتاحة فقط (غير المحجوزة)
	var slots []models.Slots
	config.ConnectDB().Where("is_booked = ?", false).Find(&slots)

	c.HTML(http.StatusOK, "booking.html", gin.H{
		"student_id": studentID, 
		"slots": slots, // نرسل المواعيد للصفحة
	})
}

// --- 3. تحديث دالة الحجز (لحفظ الموعد المختار) ---
func Create_Booking(c *gin.Context) {
	studentID := c.PostForm("student_id")
	slotID := c.PostForm("slot_id") // نستقبل رقم الموعد

	// ... (كود رفع الصورة كما هو) ...
	file, _ := c.FormFile("receipt")
	filename := fmt.Sprintf("%s_%s", studentID, filepath.Base(file.Filename))
	c.SaveUploadedFile(file, "./uploads/"+filename)

	idUint, _ := strconv.ParseUint(studentID, 10, 64)
	slotUint, _ := strconv.ParseUint(slotID, 10, 64)

	booking := models.Booking{
		StudentID:    uint(idUint),
		SlotID:       uint(slotUint), // حفظ الموعد
		PaymentImage: filename,
		Status:       "pending",
	}
	config.ConnectDB().Create(&booking)
	
	// تحديث الموعد ليصبح محجوزاً
	config.ConnectDB().Model(&models.Slots{}).Where("id = ?", slotUint).Update("is_booked", true)

	c.HTML(http.StatusOK, "booking.html", gin.H{"message": "تم الحجز بنجاح!"})
}
*/
// --- تحديث: عرض صفحة الحجز (بناءً على مستوى الطالب) ---
func Show_Booking(c *gin.Context) {
	studentID := c.Query("student_id")
	
	// 1. جلب بيانات الطالب لنعرف مستواه
	var student models.Students
	if err := config.ConnectDB().First(&student, studentID).Error; err != nil {
		c.String(404, "الطالب غير موجود")
		return
	}

	// 2. جلب المواعيد المتاحة
	var slots []models.Slots
	config.ConnectDB().Where("is_booked = ?", false).Find(&slots)

	// 3. إرسال الطالب والمواعيد لملف HTML
	c.HTML(http.StatusOK, "booking.html", gin.H{
		"student": student, 
		"slots":   slots,
	})
}

// --- تحديث: حفظ الحجز (مع السعر والنوع) ---
func Create_Booking(c *gin.Context) {
	studentID := c.PostForm("student_id")
	slotID := c.PostForm("slot_id")
	
	// استقبال بيانات الدفع الجديدة من الفورم
	bookingType := c.PostForm("booking_type") // مثلاً: assessment, monthly_l1
	amountStr := c.PostForm("amount")         // السعر

	// رفع الصورة
	file, _ := c.FormFile("receipt")
	filename := fmt.Sprintf("%s_%s", studentID, filepath.Base(file.Filename))
	c.SaveUploadedFile(file, "./uploads/"+filename)

	idUint, _ := strconv.ParseUint(studentID, 10, 64)
	slotUint, _ := strconv.ParseUint(slotID, 10, 64)
	amount, _ := strconv.ParseFloat(amountStr, 64) // تحويل السعر لرقم

	booking := models.Booking{
		StudentID:    uint(idUint),
		SlotID:       uint(slotUint),
		PaymentImage: filename,
		Status:       "pending",
		BookingType:  bookingType, // حفظ النوع
		Amount:       amount,      // حفظ المبلغ
	}
	
	config.ConnectDB().Create(&booking)
	
	// تحديث الموعد ليصبح محجوزاً
	config.ConnectDB().Model(&models.Slots{}).Where("id = ?", slotUint).Update("is_booked", true)

	c.HTML(http.StatusOK, "booking.html", gin.H{"message": "تم الحجز بنجاح! سيتم مراجعة الإيصال وتأكيد الموعد."})
}

// --- 4. دالة إضافة الدرجات (للشيخ) ---
func Admin_AddGrade(c *gin.Context) {
	studentID := c.PostForm("student_id")
	surah := c.PostForm("surah")
	verses := c.PostForm("verses")
	grade := c.PostForm("grade")
	notes := c.PostForm("notes")

	idUint, _ := strconv.ParseUint(studentID, 10, 64)
	progress := models.Progres{
		StudentID: uint(idUint),
		Date:      time.Now(),
		Surah:     surah,
		Verses:    verses,
		Grade:     grade,
		Notes:     notes,
	}
	config.ConnectDB().Create(&progress)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// --- 5. صفحة الطالب الخاصة (لعرض الدرجات) ---
// دالة عرض ملف الطالب (باستخدام الرابط السري)
func Show_Student_Profile(c *gin.Context) {
	token := c.Param("token")
	var student models.Students
	
	// 1. البحث عن الطالب باستخدام التوكن
	// نستخدم Preload لجلب الحجوزات أيضاً إذا أردت عرضها (اختياري)
	if err := config.ConnectDB().Where("magic_link_token = ?", token).First(&student).Error; err != nil {
		c.HTML(http.StatusNotFound, "profile.html", gin.H{"error": "الرابط غير صحيح أو منتهي الصلاحية"})
		return
	}

	// 2. جلب سجل الدرجات لهذا الطالب مرتبة بالأحدث
	var progress []models.Progres
	config.ConnectDB().Where("student_id = ?", student.ID).Order("date desc").Find(&progress)

	// 3. عرض الصفحة
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"student":  student,
		"progress": progress,
	})
}
// ////////////////////////////////////////////////

// --- 6. تحديث لوحة الأدمن (لجلب الطلاب للمنسدلة) ---
func Admin_Dashboard(c *gin.Context) {
	var bookings []models.Booking
	var students []models.Students // لجلب قائمة الطلاب للدرجات
	
	db := config.ConnectDB()
	db.Preload("Student").Preload("Slot").Where("status = ?", "pending").Find(&bookings) // لاحظ Preload Slot
	db.Find(&students) // هات كل الطلاب

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"bookings": bookings,
		"students": students,
	})
}
//////////////////////////////////////
/*func Admin_Approve(c *gin.Context) {
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
func Admin_Approve(c *gin.Context) {
	bookingID := c.Param("id")
	var booking models.Booking
	
	db := config.ConnectDB()

	// 1. البحث عن الحجز + بيانات الطالب (مهم جداً Preload)
	if err := db.Preload("Student").First(&booking, bookingID).Error; err != nil {
		c.JSON(404, gin.H{"error": "الحجز غير موجود"})
		return
	}

	// 2. تحديث الحالة
	booking.Status = "confirmed"
	db.Save(&booking)

	// 3. إنشاء كود سري للطالب لو مش موجود
	if booking.Student.MagicLinkToken == "" {
		booking.Student.MagicLinkToken = uuid.New().String()
		db.Save(&booking.Student)
	}

	// 4. إرسال الإيميل باستخدام الكود الجديد
	fmt.Println("📧 جاري إرسال الإيميل للطالب:", booking.Student.Email)
	
	// رابط الزوم الثابت (يمكنك تغييره لاحقاً)
	zoomLink := "https://zoom.us/j/123456789"

	err := utils.SendConfirmationEmail(
		booking.Student.Email,
		booking.Student.FirstName,
		zoomLink,
		booking.Student.MagicLinkToken,
	)

	if err != nil {
		fmt.Println("❌ فشل الإرسال:", err)
	} else {
		fmt.Println("✅ تم إرسال الإيميل بنجاح!")
	}

	// العودة لصفحة الأدمن
	c.Redirect(http.StatusFound, "/admin/dashboard")
}*/