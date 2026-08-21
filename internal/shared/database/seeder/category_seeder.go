package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	var coffeeStore store.Store
	if err := db.Where("email = ?", "kazcoffee@example.com").First(&coffeeStore).Error; err != nil {
		return err
	}

	var restaurantStore store.Store
	if err := db.Where("email = ?", "kazrestaurant@example.com").First(&restaurantStore).Error; err != nil {
		return err
	}

	categories := []menu.Category{
		{
			Name:    "Cà phê",
			Rank:    1,
			StoreID: coffeeStore.ID,
		},
		{
			Name:    "Trà",
			Rank:    2,
			StoreID: coffeeStore.ID,
		},
		{
			Name:    "Bánh ngọt",
			Rank:    3,
			StoreID: coffeeStore.ID,
		},
		{
			Name:    "Khai vị",
			Rank:    1,
			StoreID: restaurantStore.ID,
		},
		{
			Name:    "Món chính",
			Rank:    2,
			StoreID: restaurantStore.ID,
		},
		{
			Name:    "Đồ uống",
			Rank:    3,
			StoreID: restaurantStore.ID,
		},
	}

	for _, c := range categories {
		var existing menu.Category
		err := db.Where("name = ? AND store_id = ?", c.Name, c.StoreID).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}

	return nil
}
