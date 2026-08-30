package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService quản lý nghiệp vụ logic đăng nhập và tạo token.
type AuthService struct {
	userRepository *UserRepository
	jwtSecret      string // jwtSecret là khóa bí mật dùng để ký JWT
}

// NewAuthService khởi tạo AuthService nhận vào con trỏ UserRepository và khóa bí mật.
func NewAuthService(repository *UserRepository, secret string) *AuthService {
	return &AuthService{
		userRepository: repository,
		jwtSecret:      secret,
	}
}

// Login kiểm tra thông tin tài khoản và so sánh mật khẩu.
func (s *AuthService) Login(email, password string) (*User, error) {
	// 1. Dùng UserRepository tìm user theo email
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("email hoặc mật khẩu không chính xác")
	}

	// 2. So sánh mật khẩu đã hash bằng bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("email hoặc mật khẩu không chính xác")
	}

	return user, nil
}

// GenerateTokenPair sinh cặp Access Token (hạn 15 phút) và Refresh Token (hạn 7 ngày),
// đồng thời lưu Refresh Token vào Database.
func (s *AuthService) GenerateTokenPair(user *User) (string, string, error) {
	// 1. Sinh Access Token (JWT)
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"role":     user.Role,
		"store_id": user.StoreID,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // Hạn ngắn: 15 phút
		"iat":      time.Now().Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// 2. Sinh Refresh Token (Chuỗi ngẫu nhiên bảo mật dài 32 bytes = 64 ký tự hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	refreshTokenStr := hex.EncodeToString(bytes)

	// 3. Lưu Refresh Token vào Database với thời gian hết hạn 7 ngày
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	rfTokenModel := &RefreshToken{
		Token:     refreshTokenStr,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepository.CreateRefreshToken(rfTokenModel); err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenStr, nil
}

// Logout thu hồi Refresh Token bằng cách xóa nó khỏi Database.
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepository.DeleteRefreshToken(refreshToken)
}

// Refresh xác thực Refresh Token cũ và cấp cặp Access Token & Refresh Token mới (Rotation).
func (s *AuthService) Refresh(refreshTokenStr string) (string, string, error) {
	// 1. Tìm Refresh Token trong DB
	rfToken, err := s.userRepository.FindRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", errors.New("refresh token không hợp lệ hoặc đã hết hạn")
	}

	// 2. Kiểm tra xem Refresh Token có bị hết hạn không
	if time.Now().After(rfToken.ExpiresAt) {
		// Xóa token hết hạn khỏi DB để dọn dẹp
		_ = s.userRepository.DeleteRefreshToken(refreshTokenStr)
		return "", "", errors.New("refresh token đã hết hạn, vui lòng đăng nhập lại")
	}

	// 3. Đảm bảo thông tin User được nạp thành công
	if rfToken.User == nil {
		return "", "", errors.New("không tìm thấy thông tin người dùng tương ứng")
	}

	// 4. Xóa Refresh Token cũ (Refresh Token Rotation để bảo mật)
	if err := s.userRepository.DeleteRefreshToken(refreshTokenStr); err != nil {
		return "", "", err
	}

	// 5. Tạo cặp token mới
	return s.GenerateTokenPair(rfToken.User)
}
