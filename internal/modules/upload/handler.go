package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/utils"
)

type UploadHandler struct{}

// NewUploadHandler khởi tạo instance UploadHandler mới.
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadProductImage xử lý tải ảnh lên cho Sản phẩm (lưu vào thư mục "products").
func (h *UploadHandler) UploadProductImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file ảnh tải lên với key 'image'"})
		return
	}

	path, err := utils.SaveUploadedFile(file, "products")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": path,
	})
}

// UploadProfileImage xử lý tải ảnh lên cho Profile (lưu vào thư mục "profiles").
func (h *UploadHandler) UploadProfileImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file ảnh tải lên với key 'image'"})
		return
	}

	path, err := utils.SaveUploadedFile(file, "profiles")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": path,
	})
}
