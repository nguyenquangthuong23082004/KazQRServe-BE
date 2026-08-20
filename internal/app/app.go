package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
}

// Run là hàm kích hoạt server (được gọi từ main.go)
func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
