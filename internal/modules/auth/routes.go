package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes khởi tạo các dependency và đăng ký các endpoints cho module Auth.
func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	userRepo := NewUserRepository(db)
	authService := NewAuthService(userRepo, jwtSecret)
	authHandler := NewAuthHandler(authService)

	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", middleware.AuthMiddleware(jwtSecret), authHandler.Logout)
	}
}
