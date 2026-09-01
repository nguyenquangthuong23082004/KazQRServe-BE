package order

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

// Helper private đọc id từ URL param
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

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateCustomerOrder (Public API) Khách vãng lai quét QR tại bàn đặt món
func (h *OrderHandler) CreateCustomerOrder(c *gin.Context) {
	var dto CreateOrderDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đặt món không hợp lệ", "details": err.Error()})
		return
	}

	order, err := h.service.CreateCustomerOrder(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// List (Protected API) Lấy danh sách đơn hàng cho Staff / Admin
func (h *OrderHandler) List(c *gin.Context) {
	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	statusFilter := c.Query("status")
	orders, err := h.service.GetOrders(storeID, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Get (Protected API) Lấy chi tiết 1 đơn hàng theo ID
func (h *OrderHandler) Get(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã đơn hàng không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.GetOrderByID(id, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrOrderNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateStatus (Protected API) Duyệt đơn (confirmed) hoặc Hủy đơn (cancelled)
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var dto UpdateOrderStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trạng thái status là bắt buộc"})
		return
	}

	id, err := getParamID(c)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã đơn hàng không hợp lệ"})
		return
	}
	dto.OrderID = id

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	order, err := h.service.UpdateOrderStatus(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetTableSession (Protected API) Xem chi tiết Session và tổng Bill hiện tại của Bàn
func (h *OrderHandler) GetTableSession(c *gin.Context) {
	tableID, err := getParamID(c)
	if err != nil || tableID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã bàn không hợp lệ"})
		return
	}

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := h.service.GetActiveSessionByTable(tableID, storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrSessionNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// Checkout (Protected API) Thu ngân tính tiền và đóng phiên bàn
func (h *OrderHandler) Checkout(c *gin.Context) {
	var dto CheckoutSessionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phương thức thanh toán payment_method là bắt buộc"})
		return
	}

	tableID, err := getParamID(c)
	if err != nil || tableID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã bàn không hợp lệ"})
		return
	}
	dto.TableID = tableID

	storeID, err := getStoreIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dto.StoreID = storeID

	session, err := h.service.CheckoutSession(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Thanh toán và đóng phiên bàn thành công",
		"session": session,
	})
}

// GetCustomerSession (Public API) cho phép khách hàng xem chi tiết phiên đặt món của bàn mình qua mã QR token
func (h *OrderHandler) GetCustomerSession(c *gin.Context) {
	tableToken := c.Query("token")
	if tableToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã token bàn không được để trống"})
		return
	}

	session, err := h.service.GetCustomerSessionByToken(tableToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

