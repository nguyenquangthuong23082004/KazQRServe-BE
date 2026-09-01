package order

import (
	"time"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

// Hằng số trạng thái Session
const (
	SessionStatusActive = "active"
	SessionStatusClosed = "closed"
)

// Hằng số trạng thái Order
const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

type Session struct {
	model.Base
	Status      string       `gorm:"column:status;type:varchar(50);not null;default:'active'" json:"status"`
	TotalAmount float64      `gorm:"column:total_amount;type:decimal(10,2);not null;default:0" json:"total_amount"`
	PaidMethod  string       `gorm:"column:paid_method;type:varchar(50);null" json:"paid_method"`
	PaidAt      *time.Time   `gorm:"column:paid_at;type:timestamp;null" json:"paid_at"`
	TableID     uint         `gorm:"column:table_id;not null;index" json:"table_id"`
	Table       *table.Table `gorm:"foreignKey:TableID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"table,omitempty"`
	StoreID     uint         `gorm:"column:store_id;not null;default:1;index" json:"store_id"`
	Store       *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
	Orders      []Order      `gorm:"foreignKey:SessionID" json:"orders,omitempty"`
}

type Order struct {
	model.Base
	Status       string       `gorm:"column:status;type:varchar(50);not null;default:'pending'" json:"status"`
	TotalAmount  float64      `gorm:"column:total_amount;type:decimal(10,2);not null;default:0" json:"total_amount"`
	Note         string       `gorm:"column:note;type:text;null" json:"note"`
	CancelReason string       `gorm:"column:cancel_reason;type:text;null" json:"cancel_reason"`
	SessionID    uint         `gorm:"column:session_id;not null;index" json:"session_id"`
	Session      *Session     `gorm:"foreignKey:SessionID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"session,omitempty"`
	TableID      uint         `gorm:"column:table_id;not null;default:1;index" json:"table_id"`
	Table        *table.Table `gorm:"foreignKey:TableID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"table,omitempty"`
	StoreID      uint         `gorm:"column:store_id;not null;default:1;index" json:"store_id"`
	Store        *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
	Items        []OrderItem  `gorm:"foreignKey:OrderID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"items"`
}

type OrderItem struct {
	model.Base
	Quantity            int           `gorm:"column:quantity;type:int;not null" json:"quantity"`
	PriceSnapshot       float64       `gorm:"column:price_snapshot;type:decimal(10,2);not null" json:"price_snapshot"`
	ProductNameSnapshot string        `gorm:"column:product_name_snapshot;type:varchar(255);not null;default:''" json:"product_name_snapshot"`
	Note                string        `gorm:"column:note;type:text;null" json:"note"`
	OrderID             uint          `gorm:"column:order_id;not null;index" json:"order_id"`
	ProductID           uint          `gorm:"column:product_id;not null;index" json:"product_id"`
	Product             *menu.Product `gorm:"foreignKey:ProductID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"product,omitempty"`
}
