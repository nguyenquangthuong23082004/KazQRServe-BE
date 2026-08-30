package auth

import (
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

// GenerateToken sinh chuỗi JWT Token từ thông tin User.
func (s *AuthService) GenerateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"role":     user.Role,
		"store_id": user.StoreID,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
