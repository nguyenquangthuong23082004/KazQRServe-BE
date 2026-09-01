package table

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/utils"
)

var (
	ErrTableNotFound     = errors.New("bàn ăn không tồn tại hoặc không thuộc cửa hàng của bạn")
	ErrTableLimitReached = errors.New("cửa hàng của bạn đã đạt giới hạn tối đa 10 bàn ăn theo gói đăng ký. Vui lòng nâng cấp gói dịch vụ để thêm bàn mới")
	ErrInvalidStatus     = errors.New("trạng thái bàn ăn không hợp lệ (chỉ chấp nhận: available, occupied, reserved)")
)

type TableService struct {
	repo *TableRepository
}

func NewTableService(repo *TableRepository) *TableService {
	return &TableService{repo: repo}
}

// SaveTable hợp nhất luồng Save (ID > 0: Update, ID = 0: Create)
func (s *TableService) SaveTable(dto SaveTableDTO) (*Table, error) {
	if dto.ID > 0 {
		return s.updateTable(dto)
	}
	return s.createTable(dto)
}

func (s *TableService) createTable(dto SaveTableDTO) (*Table, error) {
	// 1. Kiểm tra Quota giới hạn tối đa 10 bàn cho store theo gói đăng ký
	count, err := s.repo.CountByStoreID(dto.StoreID)
	if err != nil {
		return nil, err
	}
	if count >= MaxTablesPerStore {
		return nil, ErrTableLimitReached
	}

	capacity := dto.Capacity
	if capacity <= 0 {
		capacity = 4 // Mặc định 4 ghế
	}

	// 2. Sinh UUID cho mã QR
	tableUUID := uuid.New().String()

	// 3. Xây dựng nội dung URL quét Order
	feBaseURL := os.Getenv("FE_BASE_URL")
	if feBaseURL == "" {
		feBaseURL = "http://localhost:3000"
	}
	qrContent := fmt.Sprintf("%s/order?token=%s", feBaseURL, tableUUID)

	// 4. Sinh file ảnh PNG QR Code sử dụng skip2/go-qrcode
	qrImageName := fmt.Sprintf("table_%s", tableUUID)
	qrCodeURL, err := utils.GenerateQRCodeImage(qrContent, qrImageName)
	if err != nil {
		return nil, err
	}

	table := &Table{
		Name:      dto.Name,
		Area:      dto.Area,
		Capacity:  capacity,
		Status:    StatusAvailable,
		UUID:      tableUUID,
		QRCodeURL: qrCodeURL,
		StoreID:   dto.StoreID,
	}

	if err := s.repo.Save(table); err != nil {
		// Nếu lưu DB thất bại, dọn dẹp file ảnh QR vừa sinh
		_ = utils.DeleteFile(qrCodeURL)
		return nil, err
	}
	return table, nil
}

func (s *TableService) updateTable(dto SaveTableDTO) (*Table, error) {
	table, err := s.repo.FindByIDAndStoreID(dto.ID, dto.StoreID)
	if err != nil {
		return nil, ErrTableNotFound
	}

	table.Name = dto.Name
	table.Area = dto.Area
	if dto.Capacity > 0 {
		table.Capacity = dto.Capacity
	}

	if err := s.repo.Save(table); err != nil {
		return nil, err
	}
	return table, nil
}

func (s *TableService) GetTables(storeID uint) ([]Table, error) {
	return s.repo.FindAllByStoreID(storeID)
}

func (s *TableService) GetTableByID(id uint, storeID uint) (*Table, error) {
	return s.repo.FindByIDAndStoreID(id, storeID)
}

func (s *TableService) UpdateStatus(dto UpdateTableStatusDTO) (*Table, error) {
	if dto.Status != StatusAvailable && dto.Status != StatusOccupied && dto.Status != StatusReserved {
		return nil, ErrInvalidStatus
	}

	table, err := s.repo.FindByIDAndStoreID(dto.ID, dto.StoreID)
	if err != nil {
		return nil, ErrTableNotFound
	}

	table.Status = dto.Status
	if err := s.repo.Save(table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *TableService) DeleteTable(id uint, storeID uint) error {
	table, err := s.repo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return ErrTableNotFound
	}

	if err := s.repo.Delete(id, storeID); err != nil {
		return err
	}

	// Xóa file ảnh QR Code vật lý trên server khi bàn bị xóa
	if table.QRCodeURL != "" {
		_ = utils.DeleteFile(table.QRCodeURL)
	}

	return nil
}


func (s *TableService) VerifyQRToken(uuidStr string) (*Table, error) {
	table, err := s.repo.FindByUUID(uuidStr)
	if err != nil {
		return nil, errors.New("mã QR bàn không hợp lệ hoặc không tồn tại")
	}
	return table, nil
}
