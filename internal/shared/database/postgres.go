package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// khởi tạo và trả về kết nối GORM DB
func InitDb() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable" // Mặc định disable cho môi trường dev local
	}

	// Tạo chuỗi kết nối DSN (Data Source Name) chuẩn cho GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	// Mở kết nối
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// In log câu lệnh SQL ra terminal
		// Ẩn log SQL cho môi trường production
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		// Bắt lỗi kết nối databases, in ra thông báo lỗi và exit luôn
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Kiểm tra kết nối thực tế qua Ping
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("Failed to get underlying SQL database: %v", err)
	// }

	// if err := sqlDB.Ping(); err != nil {
	// 	log.Fatalf("Failed to ping database: %v", err)
	// }

	log.Println("Database connected successfully")

	// Trả về đối tượng DB
	return db, nil
}
