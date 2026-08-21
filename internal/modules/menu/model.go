package menu

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type Category struct {
	model.Base
	Name    string       `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Rank    int          `gorm:"column:rank;type:int;not null" json:"rank"`
	StoreID uint         `gorm:"column:store_id;not null" json:"store_id"`
	Store   *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
}

type Product struct {
	model.Base
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string    `gorm:"column:description;type:text;null" json:"description"`
	Price       float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	ImageURL    string    `gorm:"column:image_url;type:varchar(255);null" json:"image_url"`
	IsAvailable bool      `gorm:"column:is_available;type:bool;not null;default:true" json:"is_available"`
	CategoryID  uint      `gorm:"column:category_id;not null" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"category,omitempty"`
}
