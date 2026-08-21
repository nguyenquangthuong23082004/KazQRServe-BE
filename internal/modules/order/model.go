package order

import (
	"time"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Session struct {
	model.Base
	Status      string    `gorm:"column:status;type:varchar(255);not null;default:active" json:"status"`
	TotalAmount float64   `gorm:"column:total_amount;type:decimal(10,2);not null" json:"total_amount"`
	PaidMethod  string    `gorm:"column:paid_method;type:varchar(255);null" json:"paid_method"`
	PaidAt      time.Time `gorm:"column:paid_at;type:timestamp;null" json:"paid_at"`
	//FK
	TableID uint         `gorm:"column:table_id;not null" json:"table_id"`
	Table   *table.Table `gorm:"foreignKey:TableID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"table,omitempty"`
}

type Order struct {
	model.Base
	Status       string `gorm:"column:status;type:varchar(255);not null;default:pending" json:"status"`
	CancelReason string `gorm:"column:cancel_reason;type:text;null" json:"cancel_reason"`
	//FK
	SessionID uint     `gorm:"column:session_id;not null" json:"session_id"`
	Session   *Session `gorm:"foreignKey:SessionID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"session,omitempty"`
}

type OrderItem struct {
	model.Base
	Quantity      int     `gorm:"column:quantity;type:int;not null" json:"quantity"`
	PriceSnapshot float64 `gorm:"column:price_snapshot;type:decimal(10,2);not null" json:"price_snapshot"`
	Note          string  `gorm:"column:notes;type:text;null" json:"notes"`
	//FK
	OrderID   uint          `gorm:"column:order_id;not null" json:"order_id"`
	Order     *Order        `gorm:"foreignKey:OrderID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"order,omitempty"`
	ProductID uint          `gorm:"column:product_id;not null" json:"product_id"`
	Product   *menu.Product `gorm:"foreignKey:ProductID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product,omitempty"`
}
