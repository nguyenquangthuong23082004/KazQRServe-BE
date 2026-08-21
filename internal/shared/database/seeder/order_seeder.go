package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/order"
	"gorm.io/gorm"
)

func SeedOrders(db *gorm.DB) error {
	var activeSession order.Session
	if err := db.Where("status = ?", "active").First(&activeSession).Error; err != nil {
		return err
	}

	var inactiveSession order.Session
	if err := db.Where("status = ?", "inactive").First(&inactiveSession).Error; err != nil {
		return err
	}

	orders := []order.Order{
		{
			Status:    "pending",
			SessionID: activeSession.ID,
		},
		{
			Status:    "completed",
			SessionID: inactiveSession.ID,
		},
	}

	for _, o := range orders {
		var existing order.Order
		err := db.Where("session_id = ? AND status = ?", o.SessionID, o.Status).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.Create(&o).Error; err != nil {
			return err
		}
	}

	return nil
}
