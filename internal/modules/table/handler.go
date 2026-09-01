package table

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

type TableHandler struct {
	service *TableService
}

func NewTableHandler(service *TableService) *TableHandler {
	return &TableHandler{service: service}
}

// Save xử lý Lưu Bàn ăn (Tạo mới nếu ID=0, Cập nhật nếu ID>0). Kiểm tra giới hạn 10 bàn.
func (h *TableHandler) Save(c *gin.Context) {
	var dto SaveTableDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên bàn không được để trống"})
		return
	}

	if paramID, err := getParamID(c); err == nil && paramID > 0 {
		dto.ID = paramID
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	table, err := h.service.SaveTable(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.ID > 0 {
		c.JSON(http.StatusOK, table)
	} else {
		c.JSON(http.StatusCreated, table)
	}
}

// List trả về danh sách tất cả bàn ăn thuộc StoreID
func (h *TableHandler) List(c *gin.Context) {
	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tables, err := h.service.GetTables(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách bàn ăn"})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// Get trả về chi tiết 1 bàn ăn theo ID
func (h *TableHandler) Get(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã bàn không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	table, err := h.service.GetTableByID(id, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrTableNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// UpdateStatus cập nhật nhanh trạng thái bàn (available / occupied / reserved)
func (h *TableHandler) UpdateStatus(c *gin.Context) {
	var dto UpdateTableStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trạng thái status là bắt buộc"})
		return
	}

	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã bàn không hợp lệ"})
		return
	}
	dto.ID = id

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	table, err := h.service.UpdateStatus(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// Delete xóa bàn ăn
func (h *TableHandler) Delete(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã bàn không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteTable(id, storeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa bàn ăn thành công"})
}

// VerifyQRToken Public API cho khách vãng lai quét QR Code tại bàn
func (h *TableHandler) VerifyQRToken(c *gin.Context) {
	uuidStr := c.Query("token")
	if uuidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã token QR không được để trống"})
		return
	}

	table, err := h.service.VerifyQRToken(uuidStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}
