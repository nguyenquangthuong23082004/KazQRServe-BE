package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/middleware"
)

// RegisterRoutes cấu hình các route upload ảnh.
func RegisterRoutes(router *gin.Engine, jwtSecret string) {
	handler := NewUploadHandler()

	uploadGroup := router.Group("/api/v1/upload")
	{
		// Cần xác thực đăng nhập để thực hiện upload ảnh
		uploadGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Chỉ Admin và Staff mới được upload ảnh sản phẩm
			uploadGroup.POST("/product", middleware.RequireRoles("admin", "staff"), handler.UploadProductImage)
			// Bất kỳ user đã đăng nhập nào cũng được upload ảnh profile
			uploadGroup.POST("/profile", handler.UploadProfileImage)
		}
	}
}
