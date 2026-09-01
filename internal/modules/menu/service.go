package menu

import (
	"errors"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/shared/utils"
)

var (
	ErrCategoryNotFound = errors.New("danh mục không tồn tại hoặc không thuộc cửa hàng của bạn")
	ErrProductNotFound  = errors.New("sản phẩm không tồn tại hoặc không thuộc cửa hàng của bạn")
)

// ==========================================
// CATEGORY SERVICE
// ==========================================

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// SaveCategory hợp nhất luồng Save (ID > 0: Update, ID = 0: Create)
func (s *CategoryService) SaveCategory(dto SaveCategoryDTO) (*Category, error) {
	if dto.ID > 0 {
		return s.updateCategory(dto)
	}
	return s.createCategory(dto)
}

func (s *CategoryService) createCategory(dto SaveCategoryDTO) (*Category, error) {
	category := &Category{
		Name:    dto.Name,
		Rank:    dto.Rank,
		StoreID: dto.StoreID,
	}

	if err := s.repo.Save(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) updateCategory(dto SaveCategoryDTO) (*Category, error) {
	category, err := s.repo.FindByIDAndStoreID(dto.ID, dto.StoreID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	category.Name = dto.Name
	category.Rank = dto.Rank

	if err := s.repo.Save(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetCategories(storeID uint) ([]Category, error) {
	return s.repo.FindAllByStoreID(storeID)
}

func (s *CategoryService) GetCategoryByID(id uint, storeID uint) (*Category, error) {
	return s.repo.FindByIDAndStoreID(id, storeID)
}

func (s *CategoryService) DeleteCategory(id uint, storeID uint) error {
	if _, err := s.repo.FindByIDAndStoreID(id, storeID); err != nil {
		return ErrCategoryNotFound
	}
	return s.repo.Delete(id, storeID)
}

// ==========================================
// PRODUCT SERVICE
// ==========================================

type ProductService struct {
	productRepo  *ProductRepository
	categoryRepo *CategoryRepository
}

func NewProductService(productRepo *ProductRepository, categoryRepo *CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// SaveProduct hợp nhất luồng Save (ID > 0: Update, ID = 0: Create)
func (s *ProductService) SaveProduct(dto SaveProductDTO) (*Product, error) {
	// 1. Kiểm tra CategoryID thuộc StoreID
	category, err := s.categoryRepo.FindByIDAndStoreID(dto.CategoryID, dto.StoreID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	if dto.ID > 0 {
		return s.updateProduct(dto, category)
	}
	return s.createProduct(dto, category)
}

func (s *ProductService) createProduct(dto SaveProductDTO, category *Category) (*Product, error) {
	product := &Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		ImageURL:    dto.ImageURL,
		IsAvailable: *dto.IsAvailable,
		CategoryID:  dto.CategoryID,
	}

	if err := s.productRepo.Save(product); err != nil {
		return nil, err
	}

	product.Category = category
	return product, nil
}

func (s *ProductService) updateProduct(dto SaveProductDTO, category *Category) (*Product, error) {
	product, err := s.productRepo.FindByIDAndStoreID(dto.ID, dto.StoreID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Xử lý thay đổi file ảnh trên đĩa server
	oldImageURL := product.ImageURL
	shouldDeleteOldImage := false

	if dto.ImageURL != "" && dto.ImageURL != oldImageURL {
		product.ImageURL = dto.ImageURL
		shouldDeleteOldImage = true
	}

	product.Name = dto.Name
	product.Description = dto.Description
	product.Price = dto.Price
	product.IsAvailable = *dto.IsAvailable
	product.CategoryID = dto.CategoryID
	product.Category = category

	if err := s.productRepo.Save(product); err != nil {
		return nil, err
	}

	// Đổi ảnh thành công mới tiến hành xóa file ảnh cũ
	if shouldDeleteOldImage {
		_ = utils.DeleteFile(oldImageURL)
	}

	return product, nil
}

func (s *ProductService) GetProducts(storeID uint) ([]Product, error) {
	return s.productRepo.FindAllByStoreID(storeID)
}

func (s *ProductService) GetProductByID(id uint, storeID uint) (*Product, error) {
	return s.productRepo.FindByIDAndStoreID(id, storeID)
}

func (s *ProductService) DeleteProduct(id uint, storeID uint) error {
	product, err := s.productRepo.FindByIDAndStoreID(id, storeID)
	if err != nil {
		return ErrProductNotFound
	}

	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	if product.ImageURL != "" {
		_ = utils.DeleteFile(product.ImageURL)
	}

	return nil
}

// UpdateProductAvailability cập nhật nhanh trạng thái Còn món / Hết món cho sản phẩm.
func (s *ProductService) UpdateProductAvailability(dto UpdateProductAvailabilityDTO) (*Product, error) {
	product, err := s.productRepo.FindByIDAndStoreID(dto.ID, dto.StoreID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	product.IsAvailable = *dto.IsAvailable
	if err := s.productRepo.Save(product); err != nil {
		return nil, err
	}

	return product, nil
}

