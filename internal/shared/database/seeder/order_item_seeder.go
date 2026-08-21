package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/order"
	"gorm.io/gorm"
)

func SeedOrderItems(db *gorm.DB) error {
	var pendingOrder order.Order
	if err := db.Where("status = ?", "pending").First(&pendingOrder).Error; err != nil {
		return err
	}

	var completedOrder order.Order
	if err := db.Where("status = ?", "completed").First(&completedOrder).Error; err != nil {
		return err
	}

	var espresso menu.Product
	if err := db.Where("name = ?", "Cà phê Espresso").First(&espresso).Error; err != nil {
		return err
	}

	var cappuccino menu.Product
	if err := db.Where("name = ?", "Cà phê Cappuccino").First(&cappuccino).Error; err != nil {
		return err
	}

	orderItems := []order.OrderItem{
		{
			Quantity:      2,
			PriceSnapshot: espresso.Price,
			Note:          "Ít đường",
			OrderID:       pendingOrder.ID,
			ProductID:     espresso.ID,
		},
		{
			Quantity:      1,
			PriceSnapshot: cappuccino.Price,
			Note:          "Nóng",
			OrderID:       completedOrder.ID,
			ProductID:     cappuccino.ID,
		},
	}

	for _, oi := range orderItems {
		var existing order.OrderItem
		err := db.Where("order_id = ? AND product_id = ?", oi.OrderID, oi.ProductID).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&oi).Error; err != nil {
			return err
		}
	}

	return nil
}
