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

		// الصفحة الرئيسية الآن هي التسجيل المدمج
    r.GET("/", routes.Show_Register_Page)
    r.POST("/register-and-book", routes.Register_And_Book)

    // مسارات الطلاب القدامى
    r.GET("/login", routes.Show_Login)
    r.POST("/login-check", routes.Check_Login)
    
    // صفحة الحجز (تستخدم للطلاب القدامى الآن بشكل أساسي)
    r.GET("/booking", routes.Show_Booking)
    r.POST("/booking/upload", routes.Create_Booking)

	

	// --- المسارات العامة ---
	//r.GET("/", routes.Show_Form)                 // الصفحة الرئيسية (تسجيل الطالب)
	//r.GET("/booking", routes.Show_Booking)       // صفحة الحجز (التي تظهر فيها المشكلة)
	//r.POST("/booking/upload", routes.Create_Booking) // رابط رفع الصورة

	// --- مسارات الـ API القديمة ---
	route := r.Group("/api/v1")
	{
		route.POST("/student/create", routes.Create_students)
		route.GET("/student/show", routes.Show_students)
	}
	//---------------------------------------------------------
	
	// تعريف حسابات الأدمن (يمكنك تغيير الاسم وكلمة السر من هنا)
	// "admin" هو اسم المستخدم
	// "123456" هي كلمة المرور
	adminAccounts := gin.Accounts{
		"admin": "123456",
	}

	// إنشاء مجموعة مسارات محمية
	adminGroup := r.Group("/admin", gin.BasicAuth(adminAccounts))
	{
		// أي رابط نضعه هنا سيطلب كلمة سر
		adminGroup.GET("/dashboard", routes.Admin_Dashboard)
		adminGroup.POST("/approve/:id", routes.Admin_Approve)
		adminGroup.POST("/add-slot", routes.Admin_AddSlot)
		adminGroup.POST("/add-grade", routes.Admin_AddGrade)
		adminGroup.POST("/update-level", routes.Admin_UpdateStudentLevel)
	}

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