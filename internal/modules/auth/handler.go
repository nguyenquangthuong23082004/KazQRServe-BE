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

	// 3. Tạo JWT Token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lỗi hệ thống khi tạo token xác thực",
		})
		return
	}

	// 4. Trả về token và thông tin cơ bản của User
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"role":     user.Role,
			"store_id": user.StoreID,
		},
	})
}
