package table

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type TableStatus string

const (
	StatusAvailable TableStatus = "available"
	StatusOccupied  TableStatus = "occupied"
	StatusReserved  TableStatus = "reserved"
)

type Table struct {
	model.Base
	Name      string       `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Area      string       `gorm:"column:area;type:varchar(255);null" json:"area"`
	Capacity  int          `gorm:"column:capacity;type:int;not null;default:4" json:"capacity"`
	Status    TableStatus  `gorm:"column:status;type:varchar(50);not null;default:'available'" json:"status"`
	UUID      string       `gorm:"column:uuid;type:varchar(255);uniqueIndex;not null" json:"uuid"`
	QRCodeURL string       `gorm:"column:qr_code_url;type:varchar(255);null" json:"qr_code_url"`
	StoreID   uint         `gorm:"column:store_id;not null;index" json:"store_id"`
	Store     *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
}

