package entity

import (
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/store/entity"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/model"
)

type Table struct {
	model.Base
	Name string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Area string `gorm:"column:area;type:varchar(255);null" json:"area"`
	UUID string `gorm:"column:uuid;type:uuid;unique;not null" json:"uuid"`
	// ---- PHẦN THIẾT LẬP MỐI QUAN HỆ GORM ---
	// Khóa ngoại (Foreign Key) lưu trữ ID của Store mà bàn này thuộc về
	StoreID uint `gorm:"column:store_id;not null" json:"store_id"`

	// 2. Struct đại diện cho Store (Dùng con trỏ * để tối ưu bộ nhớ và cho phép null nếu cần)
	// Tag constraint giúp cấu hình Cascade Delete (Nếu Store bị xóa, các Table thuộc Store đó tự động bị xóa theo)
	Store *entity.Store `gorm:"foreignKey:StoreID;references:ID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"store,omitempty"`
}
