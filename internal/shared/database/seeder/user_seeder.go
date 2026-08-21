package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/auth"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	var coffeeStore store.Store
	if err := db.Where("email = ?", "kazcoffee@example.com").First(&coffeeStore).Error; err != nil {
		return err
	}

	var restaurantStore store.Store
	if err := db.Where("email = ?", "kazrestaurant@example.com").First(&restaurantStore).Error; err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passStr := string(hashedPassword)

	users := []auth.User{
		{
			Email:    "coffee_admin@example.com",
			Password: passStr,
			Role:     "admin",
			StoreID:  coffeeStore.ID,
		},
		{
			Email:    "coffee_staff@example.com",
			Password: passStr,
			Role:     "staff",
			StoreID:  coffeeStore.ID,
		},
		{
			Email:    "restaurant_admin@example.com",
			Password: passStr,
			Role:     "admin",
			StoreID:  restaurantStore.ID,
		},
		{
			Email:    "restaurant_staff@example.com",
			Password: passStr,
			Role:     "staff",
			StoreID:  restaurantStore.ID,
		},
	}

	for _, u := range users {
		var existing auth.User
		err := db.Where("email = ?", u.Email).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
