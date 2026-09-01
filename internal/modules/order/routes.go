package order

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	orderRepo := NewOrderRepository(db)
	tableRepo := table.NewTableRepository(db)
	productRepo := menu.NewProductRepository(db)

	service := NewOrderService(orderRepo, tableRepo, productRepo)
	handler := NewOrderHandler(service)

	// 1. Public Routes (Dành cho Khách vãng lai quét QR Code tại bàn - Không cần login)
	publicGroup := router.Group("/api/v1/public/orders")
	{
		publicGroup.POST("", handler.CreateCustomerOrder)
	}

	// 2. Protected Routes (Dành cho Admin / Staff quản trị)
	orderGroup := router.Group("/api/v1/orders")
	{
		orderGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Admin & Staff xem danh sách và chi tiết đơn hàng
			orderGroup.GET("", middleware.RequireRoles("admin", "staff"), handler.List)
			orderGroup.GET("/:id", middleware.RequireRoles("admin", "staff"), handler.Get)

			// Admin & Staff duyệt đơn (confirmed) hoặc hủy đơn (cancelled)
			orderGroup.PATCH("/:id/status", middleware.RequireRoles("admin", "staff"), handler.UpdateStatus)
		}
	}

	// 3. Table Session & Checkout Routes (Dành cho Admin / Staff)
	tableSessionGroup := router.Group("/api/v1/tables")
	{
		tableSessionGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			tableSessionGroup.GET("/:id/session", middleware.RequireRoles("admin", "staff"), handler.GetTableSession)
			tableSessionGroup.POST("/:id/checkout", middleware.RequireRoles("admin", "staff"), handler.Checkout)
		}
	}
}
