package seeder

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/order"
	"gorm.io/gorm"
)

func SeedOrders(db *gorm.DB) error {
	var activeSession order.Session
	if err := db.Where("status = ?", order.SessionStatusActive).First(&activeSession).Error; err != nil {
		return err
	}

	var closedSession order.Session
	if err := db.Where("status = ?", order.SessionStatusClosed).First(&closedSession).Error; err != nil {
		return err
	}

	orders := []order.Order{
		{
			Status:      order.OrderStatusPending,
			TotalAmount: 60000.0,
			Note:        "Đơn đặt hàng chờ duyệt",
			SessionID:   activeSession.ID,
			TableID:     activeSession.TableID,
			StoreID:     activeSession.StoreID,
		},
		{
			Status:      order.OrderStatusCompleted,
			TotalAmount: 45000.0,
			Note:        "Đơn hàng đã hoàn thành",
			SessionID:   closedSession.ID,
			TableID:     closedSession.TableID,
			StoreID:     closedSession.StoreID,
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
