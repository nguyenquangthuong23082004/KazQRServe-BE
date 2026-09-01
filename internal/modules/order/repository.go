package order

import (
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// --- SESSION REPOSITORY ---

// FindActiveSessionByTableID tìm Session đang active của bàn
func (r *OrderRepository) FindActiveSessionByTableID(tableID uint) (*Session, error) {
	var s Session
	err := r.db.Where("table_id = ? AND status = ?", tableID, SessionStatusActive).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// SaveSession lưu thông tin phiên bàn
func (r *OrderRepository) SaveSession(session *Session) error {
	return r.db.Save(session).Error
}

// FindSessionByIDAndStoreID lấy chi tiết session kèm danh sách orders và items
func (r *OrderRepository) FindSessionByIDAndStoreID(sessionID uint, storeID uint) (*Session, error) {
	var s Session
	err := r.db.Preload("Table").
		Preload("Orders.Items.Product").
		Where("id = ? AND store_id = ?", sessionID, storeID).
		First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// --- ORDER REPOSITORY ---

// SaveOrder tạo hoặc cập nhật đơn hàng
func (r *OrderRepository) SaveOrder(order *Order) error {
	return r.db.Save(order).Error
}

// FindOrderByIDAndStoreID lấy chi tiết 1 đơn hàng theo ID và storeID
func (r *OrderRepository) FindOrderByIDAndStoreID(orderID uint, storeID uint) (*Order, error) {
	var o Order
	err := r.db.Preload("Table").
		Preload("Items.Product").
		Where("id = ? AND store_id = ?", orderID, storeID).
		First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// FindOrdersByStoreID lấy danh sách các đơn hàng của store (có thể lọc theo status)
func (r *OrderRepository) FindOrdersByStoreID(storeID uint, status string) ([]Order, error) {
	var orders []Order
	query := r.db.Preload("Table").Preload("Items.Product").Where("store_id = ?", storeID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("id DESC").Find(&orders).Error
	return orders, err
}

// FindOrdersBySessionID lấy tất cả các đơn hàng thuộc session
func (r *OrderRepository) FindOrdersBySessionID(sessionID uint) ([]Order, error) {
	var orders []Order
	err := r.db.Preload("Items").Where("session_id = ?", sessionID).Find(&orders).Error
	return orders, err
}
