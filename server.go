package main
import(
	//"fmt"
	"os"
	"Online-Quran/config"
	"Online-Quran/routes"
	"github.com/gin-gonic/gin"


)
/*
func main(){
	config.ConnectDB()
    
	r := gin.Default()

	route:= r.Group("/api/v1")
	{
	route.POST("/student/create", routes.Create_students)
	route.GET("/student/show", routes.Show_students)
	}
	r.Run(":8080")
}

func main() {
	config.ConnectDB()

	r := gin.Default()

	// 1. تحميل ملفات الـ HTML
	r.LoadHTMLGlob("templates/*")

	// 2. الصفحة الرئيسية تعرض الفورم
	r.GET("/", routes.Show_Form)

	route := r.Group("/api/v1")
	{
		// لاحظ أننا نستخدم نفس الدالة Create_students التي عدلناها
		route.POST("/student/create", routes.Create_students)
		route.GET("/student/show", routes.Show_students)
	}

	r.Run(":8080")
}
func main() {
    // ... الاتصال بقاعدة البيانات ...
    config.ConnectDB()
    // إنشاء مجلد uploads لو مش موجود
    os.Mkdir("./uploads", 0755)

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    
    // إتاحة مجلد الصور للمتصفح (عشان الأدمن يشوف الإيصال)
    r.Static("/uploads", "./uploads")

    // مسارات الطالب
    r.GET("/", routes.Show_Form)
    r.GET("/booking", routes.Show_Booking)       // صفحة الحجز
    r.POST("/booking/upload", routes.Create_Booking) // رفع الإيصال

    // مسارات الأدمن
    r.GET("/admin/dashboard", routes.Admin_Dashboard)
    r.POST("/admin/approve/:id", routes.Admin_Approve)

    // ... باقي الكود ...
    r.Run(":8080")
}
/////////////////////////////////////////////////////
*/
func main() {
	config.ConnectDB()

	// إنشاء مجلد الصور لو مش موجود (عشان رفع الإيصالات)
	os.Mkdir("./uploads", 0755)

	r := gin.Default()

	// تحميل ملفات HTML
	r.LoadHTMLGlob("templates/*")
	
	// جعل مجلد الصور متاحاً للمتصفح (عشان الأدمن يشوف الصور)
	r.Static("/uploads", "./uploads")

	// --- المسارات العامة ---
	r.GET("/", routes.Show_Form)                 // الصفحة الرئيسية (تسجيل الطالب)
	r.GET("/booking", routes.Show_Booking)       // صفحة الحجز (التي تظهر فيها المشكلة)
	r.POST("/booking/upload", routes.Create_Booking) // رابط رفع الصورة

	// --- مسارات الـ API القديمة ---
	route := r.Group("/api/v1")
	{
		route.POST("/student/create", routes.Create_students)
		route.GET("/student/show", routes.Show_students)
	}

	// --- مسارات الأدمن (اختياري الآن) ---
	r.GET("/admin/dashboard", routes.Admin_Dashboard)
	r.POST("/admin/approve/:id", routes.Admin_Approve)
// مسارات الأدمن الجديدة
    r.POST("/admin/add-slot", routes.Admin_AddSlot)
    r.POST("/admin/add-grade", routes.Admin_AddGrade)

    // مسار بروفايل الطالب (السحري)
    r.GET("/student/:token", routes.Show_Student_Profile)
	r.Run(":8080")
}
/*
import (
	"fmt"
	"net/smtp"
)

func main() {
	// ضع بياناتك هنا للتجربة
	from := "yasserbadr76@gmail.com"
	password := "vnkbrcvndzqmjlue"
	to := "yasserbader010@gmail.com" // ارسل لنفسك

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Test Email form Go\r\n" +
		"\r\n" +
		"This is a test email.")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)

	if err != nil {
		fmt.Println("❌ الخطأ هو:", err)
	} else {
		fmt.Println("✅ الإيميل شغال تمام!")
	}
}*/