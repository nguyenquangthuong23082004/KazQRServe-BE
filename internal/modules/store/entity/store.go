package entity

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type Store struct {
	model.Base
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Address  string `gorm:"column:address;type:varchar(255);not null" json:"address"`
	Phone    string `gorm:"column:phone;type:varchar(255);not null" json:"phone"`
	Email    string `gorm:"column:email;type:varchar(255);not null" json:"email"`
	Logo     string `gorm:"column:logo;type:varchar(255);not null" json:"logo"`
	IsActive bool   `gorm:"column:is_active;type:boolean;not null" json:"is_active"`
}
