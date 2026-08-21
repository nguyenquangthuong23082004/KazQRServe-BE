package auth

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type User struct {
	model.Base
	Email    string       `gorm:"column:email;type:varchar(255);unique;not null" json:"email"`
	Password string       `gorm:"column:password;type:varchar(255);not null" json:"password"`
	Role     string       `gorm:"column:role;type:varchar(50);not null" json:"role"`
	StoreID  uint         `gorm:"column:store_id;not null" json:"store_id"`
	Store    *store.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
}
