package database

import (
	"log"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store/entity"
	"gorm.io/gorm"
)

// Migrate sẽ chạy tất cả các file migration
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting auto migration...")

	// Khai báo tất cả các model muốn tạo bảng ở đây
	err := db.AutoMigrate(
		&entity.Store{},
	)

	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	log.Println("Auto migration completed successfully")

	return nil
}
