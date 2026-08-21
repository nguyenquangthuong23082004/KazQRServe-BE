package seeder

import (
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := SeedStores(db); err != nil {
		return err
	}
	if err := SeedUsers(db); err != nil {
		return err
	}
	if err := SeedCategories(db); err != nil {
		return err
	}
	if err := SeedProducts(db); err != nil {
		return err
	}
	if err := SeedTables(db); err != nil {
		return err
	}
	if err := SeedSessions(db); err != nil {
		return err
	}
	if err := SeedOrders(db); err != nil {
		return err
	}
	if err := SeedOrderItems(db); err != nil {
		return err
	}

	return nil
}
