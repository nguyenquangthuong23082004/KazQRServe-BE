package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes khởi tạo các dependency và đăng ký các routes cho module Menu.
func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	// 1. Categories
	repo := NewCategoryRepository(db)
	service := NewCategoryService(repo)
	handler := NewCategoryHandler(service)

	categoryGroup := router.Group("/api/v1/categories")
	{
		// Cần bảo vệ bằng xác thực đăng nhập (AuthMiddleware)
		categoryGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Admin & Staff có quyền xem danh sách và chi tiết danh mục
			categoryGroup.GET("", middleware.RequireRoles("admin", "staff"), handler.List)
			categoryGroup.GET("/:id", middleware.RequireRoles("admin", "staff"), handler.Get)

			// Chỉ tài khoản có quyền admin mới được phép Thêm mới/Cập nhật (Save) hoặc Xóa
			categoryGroup.POST("", middleware.RequireRoles("admin"), handler.Save)
			categoryGroup.PUT("/:id", middleware.RequireRoles("admin"), handler.Save)
			categoryGroup.DELETE("/:id", middleware.RequireRoles("admin"), handler.Delete)
		}
	}

	// 2. Products
	productRepo := NewProductRepository(db)
	productService := NewProductService(productRepo, repo)
	productHandler := NewProductHandler(productService)

	productGroup := router.Group("/api/v1/products")
	{
		productGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Admin & Staff có quyền xem danh sách và chi tiết sản phẩm
			productGroup.GET("", middleware.RequireRoles("admin", "staff"), productHandler.List)
			productGroup.GET("/:id", middleware.RequireRoles("admin", "staff"), productHandler.Get)

			// Cả Admin và Staff đều có quyền cập nhật nhanh trạng thái Còn/Hết món (is_available)
			productGroup.PATCH("/:id/availability", middleware.RequireRoles("admin", "staff"), productHandler.UpdateAvailability)

			// Chỉ tài khoản có quyền admin mới được phép Thêm mới/Cập nhật (Save) hoặc Xóa
			productGroup.POST("", middleware.RequireRoles("admin"), productHandler.Save)
			productGroup.PUT("/:id", middleware.RequireRoles("admin"), productHandler.Save)
			productGroup.DELETE("/:id", middleware.RequireRoles("admin"), productHandler.Delete)
		}
	}
}
