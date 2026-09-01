package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/auth"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/upload"
	"gorm.io/gorm"
)

type App struct {
	// Dùng con trỏ (*) để dùng chung 1 connection pool DB và router duy nhất,
	// tránh copy toàn bộ struct gây tốn RAM.
	db     *gorm.DB
	router *gin.Engine
}

// Trả về con trỏ (*App) để chia sẻ 1 instance duy nhất và tối ưu bộ nhớ
func NewApp(db *gorm.DB) *App {
	app := &App{
		db:     db,
		router: gin.Default(),
	}

	// Đăng ký routes khi khởi tạo App
	app.setUpRoutes()

	return app
}

// setupRoutes định nghĩa các endpoint và gắn handler/middleware
func (a *App) setUpRoutes() {
	// Endpoint kiểm tra server nhanh (Health check)
	a.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Cấu hình phục vụ file tĩnh từ thư mục uploads
	a.router.Static("/uploads", "./uploads")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret"
	}

	// Đăng ký routes cho các modules
	auth.RegisterRoutes(a.router, a.db, jwtSecret)
	menu.RegisterRoutes(a.router, a.db, jwtSecret)
	table.RegisterRoutes(a.router, a.db, jwtSecret)
	upload.RegisterRoutes(a.router, jwtSecret)
}


// Run là hàm kích hoạt server (được gọi từ main.go)
func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
