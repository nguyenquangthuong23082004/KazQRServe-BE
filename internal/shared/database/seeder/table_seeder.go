package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"gorm.io/gorm"
)

func SeedTables(db *gorm.DB) error {
	var coffeeStore store.Store
	if err := db.Where("email = ?", "kazcoffee@example.com").First(&coffeeStore).Error; err != nil {
		return err
	}

	var restaurantStore store.Store
	if err := db.Where("email = ?", "kazrestaurant@example.com").First(&restaurantStore).Error; err != nil {
		return err
	}

	tables := []table.Table{
		{
			Name:    "Bàn 1",
			Area:    "Trong nhà",
			UUID:    "e025da2c-80fc-4c4f-9e79-bf7151e2f811",
			StoreID: coffeeStore.ID,
		},
		{
			Name:    "Bàn 2",
			Area:    "Sân vườn",
			UUID:    "f51e3b08-3cb0-4e78-bc4c-f0cf468fe090",
			StoreID: coffeeStore.ID,
		},
		{
			Name:    "Bàn 101",
			Area:    "Tầng 1",
			UUID:    "74a8d461-9c6a-4b08-9dfc-27a6f272a297",
			StoreID: restaurantStore.ID,
		},
		{
			Name:    "Bàn 102",
			Area:    "Tầng 1",
			UUID:    "0c5dbb5c-44de-4d7a-b9c1-7de941bc4d91",
			StoreID: restaurantStore.ID,
		},
	}

	for _, t := range tables {
		var existing table.Table
		err := db.Where("uuid = ?", t.UUID).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&t).Error; err != nil {
			return err
		}
	}

	return nil
}
