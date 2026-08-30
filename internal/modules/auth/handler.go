package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest định nghĩa cấu trúc dữ liệu yêu cầu đăng nhập từ client.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LogoutRequest định nghĩa cấu trúc dữ liệu yêu cầu đăng xuất từ client.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshRequest định nghĩa cấu trúc dữ liệu làm mới token từ client.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthHandler quản lý việc tiếp nhận các request HTTP cho Authentication.
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler khởi tạo AuthHandler nhận vào con trỏ AuthService.
func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

// Login xử lý yêu cầu đăng nhập từ người dùng.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// 1. Validate dữ liệu đầu vào (phải là JSON hợp lệ, email đúng định dạng)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email hoặc mật khẩu không hợp lệ",
		})
		return
	}

	// 2. Gọi Service để kiểm tra thông tin đăng nhập
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Tạo cặp Access Token và Refresh Token
	accessToken, refreshToken, err := h.authService.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lỗi hệ thống khi tạo token xác thực",
		})
		return
	}

	// 4. Trả về cặp token và thông tin cơ bản của User
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"role":     user.Role,
			"store_id": user.StoreID,
		},
	})
}

// Logout xử lý yêu cầu đăng xuất từ người dùng.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest

	// 1. Validate dữ liệu đầu vào
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yêu cầu đăng xuất thiếu refresh_token",
		})
		return
	}

	// 2. Gọi Service để thu hồi Refresh Token (xóa khỏi DB)
	if err := h.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lỗi hệ thống khi thực hiện đăng xuất",
		})
		return
	}

	// 3. Trả về kết quả thành công
	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng xuất thành công",
	})
}

// Refresh xử lý yêu cầu làm mới Access Token và Refresh Token mới.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest

	// 1. Validate dữ liệu đầu vào
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yêu cầu làm mới thiếu refresh_token",
		})
		return
	}

	// 2. Gọi Service để thực hiện xoay vòng token
	accessToken, refreshToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Trả về cặp token mới
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
