package order

// CreateOrderItemDTO cho từng món trong lượt order
type CreateOrderItemDTO struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Note      string `json:"note"`
}

// CreateOrderDTO dành cho Khách vãng lai gửi order qua QR Code (Public API)
type CreateOrderDTO struct {
	TableToken string               `json:"table_token" binding:"required"`
	Note       string               `json:"note"`
	Items      []CreateOrderItemDTO `json:"items" binding:"required,min=1,dive"`
}

// UpdateOrderStatusDTO dành cho Staff/Admin cập nhật trạng thái đơn (confirmed / cancelled)
type UpdateOrderStatusDTO struct {
	OrderID      uint   `json:"-"`
	StoreID      uint   `json:"-"`
	Status       string `json:"status" binding:"required"`
	CancelReason string `json:"cancel_reason"`
}

// CheckoutSessionDTO dành cho Thu ngân thanh toán & đóng phiên bàn
type CheckoutSessionDTO struct {
	TableID       uint   `json:"-"`
	StoreID       uint   `json:"-"`
	PaymentMethod string `json:"payment_method" binding:"required"` // "cash", "banking", "momo"
}
