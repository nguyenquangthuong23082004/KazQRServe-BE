package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"gorm.io/gorm"
)

func SeedStores(db *gorm.DB) error {
	stores := []store.Store{
		{
			Name:     "Cà phê Kaz",
			Address:  "Hà Nội",
			Phone:    "0901234567",
			Email:    "kazcoffee@example.com",
			Logo:     "",
			IsActive: true,
		},
		{
			Name:     "Nhà hàng Kaz",
			Address:  "Hải Phòng",
			Phone:    "0907654321",
			Email:    "kazrestaurant@example.com",
			Logo:     "",
			IsActive: true,
		},
	}

	for _, s := range stores {
		var existing store.Store

		err := db.
			Where("email = ?", s.Email).
			First(&existing).Error

		if err == nil {
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}

	return nil
}
