package config

import (
	"log"
	"os"

	"Online-Quran/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var connectionInstance *gorm.DB

func initConnection() *gorm.DB {
	godotenv.Load()

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	dsn:= os.Getenv("DB_URL")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug() 

        // تأكد أنك تمرر كل الجداول هنا
        err = db.AutoMigrate(
            &models.Students{}, 
            &models.Booking{},  // تأكد أنها مكتوبة هنا
            &models.Slots{},    // ومكتوبة هنا
            &models.Progres{},
			)//
		if err != nil {
            log.Fatal("❌ خطأ أثناء إنشاء الجداول:", err)
        } else {
            // رسالة تأكيد ستظهر لك في التيرمينال لو كله تمام
            log.Println("✅ تم إنشاء جميع الجداول بنجاح!")
        }

        connectionInstance = db
        return connectionInstance
}
/*
	db.AutoMigrate(
		&models.Students{}, 
		&models.Slots{},
		&models.Bookings{}, 
		&models.Progres{})

	return db
}
*/
func ConnectDB() *gorm.DB {
	if connectionInstance == nil {
		connectionInstance = initConnection()
	} else {
		connectionInstance = initConnection()
	}

	return connectionInstance
}
