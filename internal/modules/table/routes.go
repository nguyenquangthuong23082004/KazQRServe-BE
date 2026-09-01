package table

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	repo := NewTableRepository(db)
	service := NewTableService(repo)
	handler := NewTableHandler(service)

	// 1. Protected Routes (Dành cho Admin / Staff quản trị)
	tableGroup := router.Group("/api/v1/tables")
	{
		tableGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Xem danh sách & chi tiết bàn (Admin & Staff)
			tableGroup.GET("", middleware.RequireRoles("admin", "staff"), handler.List)
			tableGroup.GET("/:id", middleware.RequireRoles("admin", "staff"), handler.Get)

			// Cập nhật trạng thái bàn (available / occupied / reserved) cho cả Admin & Staff
			tableGroup.PATCH("/:id/status", middleware.RequireRoles("admin", "staff"), handler.UpdateStatus)

			// Tạo mới, Cập nhật thông tin, Xóa bàn (Chỉ Admin)
			tableGroup.POST("", middleware.RequireRoles("admin"), handler.Save)
			tableGroup.PUT("/:id", middleware.RequireRoles("admin"), handler.Save)
			tableGroup.DELETE("/:id", middleware.RequireRoles("admin"), handler.Delete)
		}
	}

	// 2. Public Route (Dành cho Khách vãng lai quét QR Code tại bàn - Không cần login)
	publicGroup := router.Group("/api/v1/public/tables")
	{
		publicGroup.GET("/verify", handler.VerifyQRToken)
	}
}
