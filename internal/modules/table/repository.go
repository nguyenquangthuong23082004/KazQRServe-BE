package table

import (
	"gorm.io/gorm"
)

type TableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Save(table *Table) error {
	return r.db.Save(table).Error
}

func (r *TableRepository) FindByIDAndStoreID(id uint, storeID uint) (*Table, error) {
	var t Table
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TableRepository) FindAllByStoreID(storeID uint) ([]Table, error) {
	var tables []Table
	err := r.db.Where("store_id = ?", storeID).Order("id ASC").Find(&tables).Error
	return tables, err
}

func (r *TableRepository) FindByUUID(uuid string) (*Table, error) {
	var t Table
	err := r.db.Preload("Store").Where("uuid = ?", uuid).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TableRepository) CountByStoreID(storeID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Table{}).Where("store_id = ?", storeID).Count(&count).Error
	return count, err
}

func (r *TableRepository) Delete(id uint, storeID uint) error {
	return r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&Table{}).Error
}
