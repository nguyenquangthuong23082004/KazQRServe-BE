package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware kiểm tra tính hợp lệ của JWT Token đính kèm trong Header.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy token từ Header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Thiếu token xác thực trong header",
			})
			c.Abort() // Dừng xử lý các handler tiếp theo
			return
		}

		// 2. Kiểm tra định dạng "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Định dạng token không hợp lệ (phải là Bearer <token>)",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Parse và kiểm tra chữ ký token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Đảm bảo thuật toán ký sử dụng là HMAC (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("thuật toán ký không hợp lệ")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token không hợp lệ hoặc đã hết hạn",
			})
			c.Abort()
			return
		}

		// 4. Lấy Claims từ token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Không thể giải mã dữ liệu token",
			})
			c.Abort()
			return
		}

		// 5. Lưu thông tin giải mã từ Token vào Context của Gin
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Set("store_id", claims["store_id"])

		c.Next() // Cho phép request đi tiếp vào handler xử lý tiếp theo
	}
}

// RequireRoles kiểm tra xem role của user hiện tại có nằm trong danh sách các role được phép truy cập hay không.
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Không tìm thấy thông tin quyền truy cập trong hệ thống",
			})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Thông tin quyền truy cập không hợp lệ",
			})
			c.Abort()
			return
		}

		// Kiểm tra xem role của user có nằm trong danh sách allowedRoles không
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Bạn không có quyền thực hiện hành động này",
		})
		c.Abort()
	}
}
