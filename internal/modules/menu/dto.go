package menu

// SaveCategoryDTO định nghĩa dữ liệu đầu vào cho thao tác Lưu (Tạo/Sửa) Danh mục.
type SaveCategoryDTO struct {
	ID      uint   `json:"id"`
	StoreID uint   `json:"-"`
	Name    string `json:"name" binding:"required"`
	Rank    int    `json:"rank" binding:"required"`
}

// SaveProductDTO định nghĩa dữ liệu đầu vào cho thao tác Lưu (Tạo/Sửa) Sản phẩm.
type SaveProductDTO struct {
	ID          uint    `json:"id"`
	StoreID     uint    `json:"-"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	ImageURL    string  `json:"image_url"`
	IsAvailable *bool   `json:"is_available" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

// UpdateProductAvailabilityDTO định nghĩa dữ liệu đầu vào khi bật/tắt trạng thái Còn/Hết món.
type UpdateProductAvailabilityDTO struct {
	ID          uint  `json:"-"`
	StoreID     uint  `json:"-"`
	IsAvailable *bool `json:"is_available" binding:"required"`
}

