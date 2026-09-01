package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/menu"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/table"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound    = errors.New("đơn hàng không tồn tại hoặc không thuộc cửa hàng của bạn")
	ErrSessionNotFound  = errors.New("không tìm thấy phiên bàn đang hoạt động")
	ErrProductNotAvail  = errors.New("sản phẩm tạm hết hoặc không khả dụng")
	ErrInvalidStatus    = errors.New("trạng thái đơn hàng không hợp lệ")
	ErrTableNotVerified = errors.New("mã token bàn không hợp lệ")
)

type OrderService struct {
	orderRepo   *OrderRepository
	tableRepo   *table.TableRepository
	productRepo *menu.ProductRepository
}

func NewOrderService(
	orderRepo *OrderRepository,
	tableRepo *table.TableRepository,
	productRepo *menu.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		tableRepo:   tableRepo,
		productRepo: productRepo,
	}
}

// CreateCustomerOrder xử lý Khách vãng lai đặt món qua QR Code (Public API)
func (s *OrderService) CreateCustomerOrder(dto CreateOrderDTO) (*Order, error) {
	// 1. Xác thực mã QR Token bàn
	tbl, err := s.tableRepo.FindByUUID(dto.TableToken)
	if err != nil || tbl == nil {
		return nil, ErrTableNotVerified
	}

	storeID := tbl.StoreID
	tableID := tbl.ID

	// 2. Kiểm tra danh sách món & tính tổng tiền lượt order
	var orderItems []OrderItem
	var totalOrderAmount float64

	for _, itemDTO := range dto.Items {
		prod, err := s.productRepo.FindByIDAndStoreID(itemDTO.ProductID, storeID)
		if err != nil || prod == nil {
			return nil, fmt.Errorf("món ăn mã %d không tồn tại trong cửa hàng", itemDTO.ProductID)
		}
		if !prod.IsAvailable {
			return nil, fmt.Errorf("món '%s' hiện đã hết, vui lòng chọn món khác", prod.Name)
		}

		itemAmount := prod.Price * float64(itemDTO.Quantity)
		totalOrderAmount += itemAmount

		orderItems = append(orderItems, OrderItem{
			ProductID:           prod.ID,
			Quantity:            itemDTO.Quantity,
			PriceSnapshot:       prod.Price,
			ProductNameSnapshot: prod.Name,
			Note:                itemDTO.Note,
		})
	}

	// 3. Tìm hoặc Tạo mới Session đang active cho bàn
	session, err := s.orderRepo.FindActiveSessionByTableID(tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Tạo Session mới nếu bàn chưa có Session active
			session = &Session{
				TableID:     tableID,
				StoreID:     storeID,
				Status:      SessionStatusActive,
				TotalAmount: 0,
			}
			if err := s.orderRepo.SaveSession(session); err != nil {
				return nil, fmt.Errorf("không thể tạo phiên làm việc mới cho bàn: %v", err)
			}

			// Cập nhật trạng thái bàn sang occupied (Đang có khách)
			tbl.Status = table.StatusOccupied
			_ = s.tableRepo.Save(tbl)
		} else {
			return nil, err
		}
	}

	// 4. Tạo đơn hàng mới ở trạng thái 'pending' (Chờ nhân viên duyệt)
	newOrder := &Order{
		Status:      OrderStatusPending,
		TotalAmount: totalOrderAmount,
		Note:        dto.Note,
		SessionID:   session.ID,
		TableID:     tableID,
		StoreID:     storeID,
		Items:       orderItems,
	}

	if err := s.orderRepo.SaveOrder(newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}

// UpdateOrderStatus dành cho Staff/Admin duyệt đơn (confirmed) hoặc Hủy đơn (cancelled)
func (s *OrderService) UpdateOrderStatus(dto UpdateOrderStatusDTO) (*Order, error) {
	if dto.Status != OrderStatusConfirmed && dto.Status != OrderStatusCancelled && dto.Status != OrderStatusCompleted {
		return nil, ErrInvalidStatus
	}

	order, err := s.orderRepo.FindOrderByIDAndStoreID(dto.OrderID, dto.StoreID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	oldStatus := order.Status
	order.Status = dto.Status
	if dto.Status == OrderStatusCancelled && dto.CancelReason != "" {
		order.CancelReason = dto.CancelReason
	}

	if err := s.orderRepo.SaveOrder(order); err != nil {
		return nil, err
	}

	// Nếu đơn được duyệt sang 'confirmed' hoặc 'completed', cập nhật lại tổng tiền cho Session
	if dto.Status == OrderStatusConfirmed || dto.Status == OrderStatusCompleted {
		s.recalculateSessionTotal(order.SessionID)
	}

	// Nếu đơn bị hủy ('cancelled') và trước đó là đơn duy nhất trong Session
	if dto.Status == OrderStatusCancelled && oldStatus == OrderStatusPending {
		s.checkAndCleanupEmptySession(order.SessionID, dto.StoreID, order.TableID)
	}

	return order, nil
}

// recalculateSessionTotal tính lại tổng tiền tất cả các đơn đã duyệt trong Session
func (s *OrderService) recalculateSessionTotal(sessionID uint) {
	orders, err := s.orderRepo.FindOrdersBySessionID(sessionID)
	if err != nil {
		return
	}

	var sessionTotal float64
	for _, o := range orders {
		if o.Status == OrderStatusConfirmed || o.Status == OrderStatusCompleted {
			sessionTotal += o.TotalAmount
		}
	}

	session, err := s.orderRepo.FindSessionByIDAndStoreID(sessionID, orders[0].StoreID)
	if err == nil && session != nil {
		session.TotalAmount = sessionTotal
		_ = s.orderRepo.SaveSession(session)
	}
}

// checkAndCleanupEmptySession nếu mọi đơn trong session bị hủy, tự động đóng session và trả bàn về available
func (s *OrderService) checkAndCleanupEmptySession(sessionID uint, storeID uint, tableID uint) {
	orders, err := s.orderRepo.FindOrdersBySessionID(sessionID)
	if err != nil || len(orders) == 0 {
		return
	}

	hasValidOrder := false
	for _, o := range orders {
		if o.Status != OrderStatusCancelled {
			hasValidOrder = true
			break
		}
	}

	// Nếu tất cả các đơn đều bị hủy -> Đóng session & Đặt bàn về available
	if !hasValidOrder {
		session, err := s.orderRepo.FindSessionByIDAndStoreID(sessionID, storeID)
		if err == nil && session != nil {
			session.Status = SessionStatusClosed
			_ = s.orderRepo.SaveSession(session)
		}

		tbl, err := s.tableRepo.FindByIDAndStoreID(tableID, storeID)
		if err == nil && tbl != nil {
			tbl.Status = table.StatusAvailable
			_ = s.tableRepo.Save(tbl)
		}
	}
}

// GetOrders lấy danh sách tất cả các đơn hàng của store (lọc theo status)
func (s *OrderService) GetOrders(storeID uint, status string) ([]Order, error) {
	return s.orderRepo.FindOrdersByStoreID(storeID, status)
}

// GetOrderByID lấy chi tiết 1 đơn hàng
func (s *OrderService) GetOrderByID(orderID uint, storeID uint) (*Order, error) {
	return s.orderRepo.FindOrderByIDAndStoreID(orderID, storeID)
}

// GetActiveSessionByTable lấy thông tin phiên bàn hiện tại đang active
func (s *OrderService) GetActiveSessionByTable(tableID uint, storeID uint) (*Session, error) {
	session, err := s.orderRepo.FindActiveSessionByTableID(tableID)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	return s.orderRepo.FindSessionByIDAndStoreID(session.ID, storeID)
}

// CheckoutSession Thu ngân tính tiền, thanh toán và đóng phiên làm việc của bàn
func (s *OrderService) CheckoutSession(dto CheckoutSessionDTO) (*Session, error) {
	session, err := s.orderRepo.FindActiveSessionByTableID(dto.TableID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// 1. Tính lại tổng tiền của các đơn đã duyệt trong session
	s.recalculateSessionTotal(session.ID)

	// 2. Lấy chi tiết phiên để cập nhật thanh toán
	fullSession, err := s.orderRepo.FindSessionByIDAndStoreID(session.ID, dto.StoreID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	fullSession.PaidMethod = dto.PaymentMethod
	fullSession.PaidAt = &now
	fullSession.Status = SessionStatusClosed

	if err := s.orderRepo.SaveSession(fullSession); err != nil {
		return nil, err
	}

	// 3. Giải phóng Bàn về trạng thái available (Bàn trống)
	tbl, err := s.tableRepo.FindByIDAndStoreID(dto.TableID, dto.StoreID)
	if err == nil && tbl != nil {
		tbl.Status = table.StatusAvailable
		_ = s.tableRepo.Save(tbl)
	}

	return fullSession, nil
}
