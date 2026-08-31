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

			// Chỉ tài khoản có quyền admin mới được phép Thêm mới, Cập nhật, Xóa
			categoryGroup.POST("", middleware.RequireRoles("admin"), handler.Create)
			categoryGroup.PUT("/:id", middleware.RequireRoles("admin"), handler.Update)
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
			// Chỉ tài khoản có quyền admin mới được phép thêm sản phẩm mới
			productGroup.POST("", middleware.RequireRoles("admin"), productHandler.Create)
		}
	}
}
