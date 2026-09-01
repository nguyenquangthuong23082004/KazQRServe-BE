package table

// MaxTablesPerStore Giới hạn mặc định tối đa 10 bàn / store theo gói đăng ký ban đầu.
const MaxTablesPerStore = 10

// SaveTableDTO định nghĩa dữ liệu đầu vào cho thao tác Lưu (Tạo mới / Cập nhật) Bàn ăn.
type SaveTableDTO struct {
	ID       uint   `json:"id"`
	StoreID  uint   `json:"-"`
	Name     string `json:"name" binding:"required"`
	Area     string `json:"area"`
	Capacity int    `json:"capacity"`
}

// UpdateTableStatusDTO định nghĩa dữ liệu đầu vào khi đổi trạng thái bàn (available / occupied / reserved).
type UpdateTableStatusDTO struct {
	ID      uint        `json:"-"`
	StoreID uint        `json:"-"`
	Status  TableStatus `json:"status" binding:"required"`
}
