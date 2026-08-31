package auth

import "gorm.io/gorm"

// UserRepository quản lý các thao tác truy vấn dữ liệu User trong Database.
type UserRepository struct {
	// Con trỏ tới gorm.DB để dùng chung database connection.
	db *gorm.DB
}

// NewUserRepository khởi tạo UserRepository mới nhận kết nối DB.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindUserByEmail tìm User theo email.
func (r *UserRepository) FindUserByEmail(email string) (*User, error) {
	var user User

	// GORM ánh xạ kết quả query từ bảng users vào struct User.
	if err := r.db.
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	// Trả về con trỏ tới User để Service sử dụng.
	return &user, nil
}

// CreateRefreshToken lưu Refresh Token mới vào Database.
func (r *UserRepository) CreateRefreshToken(token *RefreshToken) error {
	return r.db.Create(token).Error
}

// RevokeRefreshToken cập nhật trạng thái Refresh Token thành đã thu hồi (is_revoked = true).
func (r *UserRepository) RevokeRefreshToken(token string) error {
	return r.db.Model(&RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}

// FindRefreshToken tìm kiếm Refresh Token trong Database (nạp sẵn thông tin User).
func (r *UserRepository) FindRefreshToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := r.db.Preload("User").Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}
