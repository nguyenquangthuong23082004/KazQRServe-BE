package menu

import "gorm.io/gorm"

// CategoryRepository chịu trách nhiệm truy vấn bảng categories trong database.
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository khởi tạo CategoryRepository mới nhận kết nối DB.
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// Create lưu một danh mục (Category) mới vào database.
func (r *CategoryRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}

// FindAllByStoreID lấy danh sách danh mục thuộc storeID, sắp xếp theo rank tăng dần.
func (r *CategoryRepository) FindAllByStoreID(storeID uint) ([]Category, error) {
	var categories []Category
	if err := r.db.
		Where("store_id = ?", storeID).
		Order("rank asc").
		Find(&categories).
		Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByIDAndStoreID tìm 1 danh mục theo ID và storeID.
func (r *CategoryRepository) FindByIDAndStoreID(id uint, storeID uint) (*Category, error) {
	var category Category
	if err := r.db.
		Where("id = ? AND store_id = ?", id, storeID).
		First(&category).
		Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Update cập nhật thông tin danh mục.
func (r *CategoryRepository) Update(category *Category) error {
	return r.db.Save(category).Error
}

// Delete xóa danh mục theo ID và storeID.
func (r *CategoryRepository) Delete(id uint, storeID uint) error {
	return r.db.
		Where("id = ? AND store_id = ?", id, storeID).
		Delete(&Category{}).
		Error
}

// ProductRepository chịu trách nhiệm truy vấn bảng products trong database.
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository khởi tạo ProductRepository mới nhận kết nối DB.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create lưu một sản phẩm (Product) mới vào database.
func (r *ProductRepository) Create(product *Product) error {
	return r.db.Create(product).Error
}
