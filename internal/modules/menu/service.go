package menu

// CategoryService chịu trách nhiệm xử lý các logic nghiệp vụ (business logic) liên quan đến Danh mục.
type CategoryService struct {
	repo *CategoryRepository
}

// NewCategoryService khởi tạo CategoryService mới nhận repo tương ứng.
func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// CreateCategory xử lý việc thêm mới một Danh mục cho cửa hàng.
func (s *CategoryService) CreateCategory(name string, rank int, storeID uint) (*Category, error) {
	category := &Category{
		Name:    name,
		Rank:    rank,
		StoreID: storeID,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategories trả về danh sách danh mục thuộc storeID.
func (s *CategoryService) GetCategories(storeID uint) ([]Category, error) {
	return s.repo.FindAllByStoreID(storeID)
}

// GetCategoryByID trả về danh mục cụ thể theo ID và storeID.
func (s *CategoryService) GetCategoryByID(id uint, storeID uint) (*Category, error) {
	return s.repo.FindByIDAndStoreID(id, storeID)
}

// UpdateCategory cập nhật thông tin của danh mục.
func (s *CategoryService) UpdateCategory(id uint, storeID uint, name string, rank int) (*Category, error) {
	// 1. Kiểm tra xem danh mục có tồn tại và thuộc storeID này hay không
	category, err := s.repo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return nil, err
	}

	// 2. Cập nhật thông tin mới
	category.Name = name
	category.Rank = rank

	// 3. Lưu vào Database
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory xóa danh mục cụ thể.
func (s *CategoryService) DeleteCategory(id uint, storeID uint) error {
	// 1. Kiểm tra xem danh mục có tồn tại và thuộc storeID hay không (để tránh xóa nhầm)
	_, err := s.repo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return err
	}

	// 2. Thực hiện xóa
	return s.repo.Delete(id, storeID)
}
