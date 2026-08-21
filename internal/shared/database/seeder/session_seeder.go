package seeder

import (
	"time"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/order"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"gorm.io/gorm"
)

func SeedSessions(db *gorm.DB) error {
	var table1 table.Table
	if err := db.Where("uuid = ?", "e025da2c-80fc-4c4f-9e79-bf7151e2f811").First(&table1).Error; err != nil {
		return err
	}

	var table2 table.Table
	if err := db.Where("uuid = ?", "f51e3b08-3cb0-4e78-bc4c-f0cf468fe090").First(&table2).Error; err != nil {
		return err
	}

	sessions := []order.Session{
		{
			Status:      "active",
			TotalAmount: 0.0,
			TableID:     table1.ID,
		},
		{
			Status:      "inactive",
			TotalAmount: 70000.0,
			PaidMethod:  "tiền mặt",
			PaidAt:      time.Now(),
			TableID:     table2.ID,
		},
	}

	for _, s := range sessions {
		var existing order.Session
		err := db.Where("table_id = ? AND status = ?", s.TableID, s.Status).First(&existing).Error
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
