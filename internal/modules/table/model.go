package table

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type Table struct {
	model.Base
	Name    string       `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Area    string       `gorm:"column:area;type:varchar(255);null" json:"area"`
	UUID    string       `gorm:"column:uuid;type:uuid;unique;not null" json:"uuid"`
	StoreID uint         `gorm:"column:store_id;not null" json:"store_id"`
	Store   *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
}
