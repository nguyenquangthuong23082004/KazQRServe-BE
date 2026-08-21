package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) error {
	var coffeeStore store.Store
	if err := db.Where("email = ?", "kazcoffee@example.com").First(&coffeeStore).Error; err != nil {
		return err
	}

	var restaurantStore store.Store
	if err := db.Where("email = ?", "kazrestaurant@example.com").First(&restaurantStore).Error; err != nil {
		return err
	}

	// Helper to find category ID
	getCategoryID := func(name string, storeID uint) uint {
		var cat menu.Category
		db.Where("name = ? AND store_id = ?", name, storeID).First(&cat)
		return cat.ID
	}

	coffeeCatID := getCategoryID("Cà phê", coffeeStore.ID)
	teaCatID := getCategoryID("Trà", coffeeStore.ID)
	dessertCatID := getCategoryID("Bánh ngọt", coffeeStore.ID)

	starterCatID := getCategoryID("Khai vị", restaurantStore.ID)
	mainCatID := getCategoryID("Món chính", restaurantStore.ID)
	drinksCatID := getCategoryID("Đồ uống", restaurantStore.ID)

	products := []menu.Product{
		{
			Name:        "Cà phê Espresso",
			Description: "Cà phê đen đậm đà kiểu Ý",
			Price:       30000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  coffeeCatID,
		},
		{
			Name:        "Cà phê Cappuccino",
			Description: "Espresso kết hợp sữa nóng và bọt sữa béo mịn",
			Price:       45000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  coffeeCatID,
		},
		{
			Name:        "Trà đào",
			Description: "Trà đào thơm ngon cùng đào miếng giòn ngọt",
			Price:       40000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  teaCatID,
		},
		{
			Name:        "Bánh Tiramisu",
			Description: "Bánh ngọt hương vị cà phê và kem phô mai kiểu Ý",
			Price:       50000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  dessertCatID,
		},
		{
			Name:        "Salad Caesar",
			Description: "Salad xà lách tươi giòn kèm bánh mì nướng và sốt Caesar",
			Price:       65000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  starterCatID,
		},
		{
			Name:        "Bít tết bò",
			Description: "Thịt bò áp chảo thượng hạng kèm sốt tiêu đen đặc biệt",
			Price:       250000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  mainCatID,
		},
		{
			Name:        "Nước cam vắt",
			Description: "Nước cam tươi nguyên chất giàu vitamin C",
			Price:       45000.0,
			ImageURL:    "",
			IsAvailable: true,
			CategoryID:  drinksCatID,
		},
	}

	for _, p := range products {
		var existing menu.Product
		err := db.Where("name = ? AND category_id = ?", p.Name, p.CategoryID).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}
