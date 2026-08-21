package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/app"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/database"
)

func main() {
	// Bắt đầu hàm main() bằng việc đọc file cấu hình:
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalln("Error loading .env file")
	}

	// Khởi  tạo kết nối databases
	db, err := database.InitDb()
	if err != nil {
		log.Fatalln("Error initializing database: ", err)
	}

	// tự động tạo hoặc nâng cấp cấu trúc bảng
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalln("Error migrating database: ", err)
	}

	// Truyền DB vừa kết nối vào App
	application := app.NewApp(db)

	// Khởi chạy Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port :%s...\n", port)

	// 5. Khởi chạy Server
	if err := application.Run(":" + port); err != nil {
		log.Fatalln("Error running application:", err)
	}

}
