package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes khởi tạo các dependency và đăng ký các endpoints cho module Auth.
func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	userRepo := NewUserRepository(db)
	authService := NewAuthService(userRepo, jwtSecret)
	authHandler := NewAuthHandler(authService)

	authGroup := router.Group("/api/v1/auth")
	{
		// Các route công khai không cần token
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)

		// Các route yêu cầu kiểm tra JWT Access Token
		// protected := authGroup.Group("")
		// protected.Use(AuthMiddleware(jwtSecret))
		// {
		// 	protected.GET("/me", authHandler.Me)
		// }
	}
}
