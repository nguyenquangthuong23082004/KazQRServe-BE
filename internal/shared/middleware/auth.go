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
