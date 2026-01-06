package main

import (
	"Online-Quran/config"
	"Online-Quran/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	os.Mkdir("./uploads", 0755)

	r := gin.Default()

	// تحميل ملفات HTML
	r.LoadHTMLGlob("templates/*.html")

	r.Static("/images", "./templates/images")
	r.Static("/uploads", "./uploads")

	// الصفحة الرئيسية (التصميم الجديد)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})

	// المسارات العامة (تسجيل ودخول)
	r.GET("/register", routes.Show_Register_Page)
	r.POST("/register-and-book", routes.Register_And_Book)
	r.GET("/login", routes.Show_Login)
	r.POST("/login-check", routes.Check_Login)

	// مسار الحجز (محمي نوعاً ما بالـ ID)
	r.GET("/booking", routes.Show_Booking)
	r.POST("/booking/upload", routes.Create_Booking) // هذا يستخدم للتجديد

	// مسار البروفايل (عام لكن سري بالتوكن)
	r.GET("/student/:token", routes.Show_Student_Profile)

	// منطقة الأدمن
	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": "123456"}))
	{
		admin.GET("/dashboard", routes.Admin_Dashboard)
		admin.POST("/approve/:id", routes.Admin_Approve)
		admin.POST("/reject/:id", routes.Admin_Reject)
		admin.POST("/add-slot", routes.Admin_AddSlot)
		admin.POST("/add-grade", routes.Admin_AddGrade)
		admin.POST("/update-level", routes.Admin_UpdateStudentLevel)
		admin.POST("/delete-booking/:id", routes.Admin_DeleteBooking)
		admin.POST("/delete-slot/:id", routes.Admin_DeleteSlot)
		admin.POST("/update-slot/:id", routes.Admin_UpdateSlot)
	}

	r.Run(":8080")
}
