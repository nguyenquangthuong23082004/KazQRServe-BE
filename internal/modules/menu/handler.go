package menu

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCategoryRequest định nghĩa dữ liệu đầu vào khi tạo Danh mục.
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Rank int    `json:"rank" binding:"required"`
}

// CategoryHandler quản lý các request HTTP liên quan đến Danh mục.
type CategoryHandler struct {
	service *CategoryService
}

// NewCategoryHandler khởi tạo CategoryHandler nhận vào service tương ứng.
func NewCategoryHandler(service *CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// Create xử lý yêu cầu thêm mới một Danh mục.
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest

	// 1. Validate dữ liệu JSON gửi lên
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tên danh mục và thứ tự hiển thị không hợp lệ",
		})
		return
	}

	// 2. Lấy store_id từ Context (đã được lưu ở AuthMiddleware)
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 3. Gọi Service xử lý nghiệp vụ
	category, err := h.service.CreateCategory(req.Name, req.Rank, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo danh mục món ăn",
		})
		return
	}

	// 4. Trả về kết quả
	c.JSON(http.StatusCreated, category)
}

// List xử lý yêu cầu lấy danh sách toàn bộ danh mục của cửa hàng.
func (h *CategoryHandler) List(c *gin.Context) {
	// 1. Lấy store_id từ Context (đã được lưu ở AuthMiddleware)
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 2. Gọi Service để lấy danh sách
	categories, err := h.service.GetCategories(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách danh mục món ăn",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Get xử lý yêu cầu lấy thông tin chi tiết của một danh mục.
func (h *CategoryHandler) Get(c *gin.Context) {
	// 1. Lấy id từ URL param
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mã danh mục không hợp lệ",
		})
		return
	}

	// 2. Lấy store_id từ Context
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 3. Gọi Service để lấy thông tin chi tiết
	category, err := h.service.GetCategoryByID(uint(id), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy danh mục hoặc danh mục không thuộc cửa hàng của bạn",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategoryRequest định nghĩa dữ liệu cập nhật cho danh mục.
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Rank int    `json:"rank" binding:"required"`
}

// Update xử lý yêu cầu cập nhật thông tin danh mục.
func (h *CategoryHandler) Update(c *gin.Context) {
	// 1. Lấy id từ URL param
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mã danh mục không hợp lệ",
		})
		return
	}

	// 2. Validate dữ liệu cập nhật
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu cập nhật danh mục không hợp lệ",
		})
		return
	}

	// 3. Lấy store_id từ Context
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 4. Gọi Service thực hiện cập nhật
	category, err := h.service.UpdateCategory(uint(id), storeID, req.Name, req.Rank)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không thể cập nhật danh mục (không tìm thấy hoặc không thuộc cửa hàng của bạn)",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete xử lý yêu cầu xóa một danh mục.
func (h *CategoryHandler) Delete(c *gin.Context) {
	// 1. Lấy id từ URL param
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mã danh mục không hợp lệ",
		})
		return
	}

	// 2. Lấy store_id từ Context
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 3. Gọi Service thực hiện xóa
	if err := h.service.DeleteCategory(uint(id), storeID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không thể xóa danh mục (không tìm thấy hoặc không thuộc cửa hàng của bạn)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa danh mục thành công",
	})
}

// CreateProductRequest định nghĩa dữ liệu đầu vào khi tạo Sản phẩm.
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	ImageURL    string  `json:"image_url"`
	IsAvailable *bool   `json:"is_available" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

// ProductHandler quản lý các request HTTP liên quan đến Sản phẩm.
type ProductHandler struct {
	service *ProductService
}

// NewProductHandler khởi tạo ProductHandler nhận vào service tương ứng.
func NewProductHandler(service *ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// Create xử lý yêu cầu thêm mới một Sản phẩm.
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest

	// 1. Validate dữ liệu JSON gửi lên
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Thông tin sản phẩm không hợp lệ (tên, giá >= 0, category_id và trạng thái là bắt buộc)",
		})
		return
	}

	// 2. Lấy store_id từ Context (đã được lưu ở AuthMiddleware)
	storeIDVal, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy thông tin cửa hàng của người dùng",
		})
		return
	}

	var storeID uint
	switch v := storeIDVal.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng mã cửa hàng không hợp lệ",
		})
		return
	}

	// 3. Gọi Service xử lý nghiệp vụ
	product, err := h.service.CreateProduct(
		req.Name,
		req.Description,
		req.Price,
		req.ImageURL,
		*req.IsAvailable,
		req.CategoryID,
		storeID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 4. Trả về kết quả
	c.JSON(http.StatusCreated, product)
}

