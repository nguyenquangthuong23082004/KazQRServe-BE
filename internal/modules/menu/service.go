package menu

import "errors"

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

// ProductService chịu trách nhiệm xử lý các logic nghiệp vụ (business logic) liên quan đến Sản phẩm.
type ProductService struct {
	productRepo  *ProductRepository
	categoryRepo *CategoryRepository
}

// NewProductService khởi tạo ProductService mới nhận các repo tương ứng.
func NewProductService(productRepo *ProductRepository, categoryRepo *CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct xử lý việc thêm mới một Sản phẩm cho cửa hàng.
func (s *ProductService) CreateProduct(name string, description string, price float64, imageURL string, isAvailable bool, categoryID uint, storeID uint) (*Product, error) {
	// 1. Xác thực xem CategoryID có tồn tại và thuộc storeID này hay không
	category, err := s.categoryRepo.FindByIDAndStoreID(categoryID, storeID)
	if err != nil {
		return nil, errors.New("danh mục không tồn tại hoặc không thuộc cửa hàng của bạn")
	}

	// 2. Tạo đối tượng Product
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    imageURL,
		IsAvailable: isAvailable,
		CategoryID:  categoryID,
	}

	// 3. Lưu vào Database
	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	product.Category = category
	return product, nil
}

// GetProducts trả về danh sách sản phẩm thuộc storeID.
func (s *ProductService) GetProducts(storeID uint) ([]Product, error) {
	return s.productRepo.FindAllByStoreID(storeID)
}

// GetProductByID trả về sản phẩm cụ thể theo ID và storeID.
func (s *ProductService) GetProductByID(id uint, storeID uint) (*Product, error) {
	return s.productRepo.FindByIDAndStoreID(id, storeID)
}

// UpdateProduct cập nhật thông tin sản phẩm.
func (s *ProductService) UpdateProduct(id uint, storeID uint, name string, description string, price float64, imageURL string, isAvailable bool, categoryID uint) (*Product, error) {
	// 1. Kiểm tra xem sản phẩm có tồn tại và thuộc storeID này hay không
	product, err := s.productRepo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return nil, errors.New("sản phẩm không tồn tại hoặc không thuộc cửa hàng của bạn")
	}

	// 2. Nếu đổi danh mục, kiểm tra xem danh mục mới có thuộc storeID này không
	var category *Category
	if product.CategoryID != categoryID {
		category, err = s.categoryRepo.FindByIDAndStoreID(categoryID, storeID)
		if err != nil {
			return nil, errors.New("danh mục mới không tồn tại hoặc không thuộc cửa hàng của bạn")
		}
		product.CategoryID = categoryID
		product.Category = category
	}

	// 3. Cập nhật thông tin khác
	product.Name = name
	product.Description = description
	product.Price = price
	product.ImageURL = imageURL
	product.IsAvailable = isAvailable

	// 4. Lưu thay đổi
	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct xóa sản phẩm.
func (s *ProductService) DeleteProduct(id uint, storeID uint) error {
	// 1. Kiểm tra xem sản phẩm có tồn tại và thuộc storeID này không
	_, err := s.productRepo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return errors.New("sản phẩm không tồn tại hoặc không thuộc cửa hàng của bạn")
	}

	// 2. Xóa
	return s.productRepo.Delete(id)
}
