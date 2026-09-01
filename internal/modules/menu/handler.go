package menu

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper private lấy store_id từ Context (sau AuthMiddleware)
func getStoreIDFromContext(c *gin.Context) (uint, error) {
	val, exists := c.Get("store_id")
	if !exists {
		return 0, errors.New("không tìm thấy thông tin cửa hàng của người dùng")
	}

	switch v := val.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	default:
		return 0, errors.New("định dạng mã cửa hàng không hợp lệ")
	}
}

// Helper private đọc id từ URL param (ví dụ: /:id)
func getParamID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("mã định danh không hợp lệ")
	}
	return uint(id), nil
}

// ==========================================
// CATEGORY HANDLER
// ==========================================

type CategoryHandler struct {
	service *CategoryService
}

func NewCategoryHandler(service *CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) Save(c *gin.Context) {
	var dto SaveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên danh mục và thứ tự hiển thị không hợp lệ"})
		return
	}

	// Lấy ID từ URL Param nếu có (ví dụ: PUT /categories/:id)
	if paramID, err := getParamID(c); err == nil && paramID > 0 {
		dto.ID = paramID
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	category, err := h.service.SaveCategory(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.ID > 0 {
		c.JSON(http.StatusOK, category)
	} else {
		c.JSON(http.StatusCreated, category)
	}
}

func (h *CategoryHandler) List(c *gin.Context) {
	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	categories, err := h.service.GetCategories(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách danh mục món ăn"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã danh mục không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	category, err := h.service.GetCategoryByID(id, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrCategoryNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã danh mục không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteCategory(id, storeID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa danh mục thành công"})
}

// ==========================================
// PRODUCT HANDLER
// ==========================================

type ProductHandler struct {
	service *ProductService
}

func NewProductHandler(service *ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Save(c *gin.Context) {
	var dto SaveProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Thông tin sản phẩm không hợp lệ (tên, giá >= 0, category_id và trạng thái là bắt buộc)",
		})
		return
	}

	// Lấy ID từ URL Param nếu có (ví dụ: PUT /products/:id)
	if paramID, err := getParamID(c); err == nil && paramID > 0 {
		dto.ID = paramID
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	product, err := h.service.SaveProduct(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.ID > 0 {
		c.JSON(http.StatusOK, product)
	} else {
		c.JSON(http.StatusCreated, product)
	}
}

func (h *ProductHandler) List(c *gin.Context) {
	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	products, err := h.service.GetProducts(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách sản phẩm"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã sản phẩm không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.GetProductByID(id, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrProductNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã sản phẩm không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteProduct(id, storeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa sản phẩm thành công"})
}

func (h *ProductHandler) UpdateAvailability(c *gin.Context) {
	var dto UpdateProductAvailabilityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trạng thái is_available (true/false) là bắt buộc"})
		return
	}

	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã sản phẩm không hợp lệ"})
		return
	}
	dto.ID = id

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	product, err := h.service.UpdateProductAvailability(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

