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

// Save thực hiện lưu (Create nếu ID=0, Update nếu ID>0) cho Danh mục.
func (r *CategoryRepository) Save(category *Category) error {
	return r.db.Save(category).Error
}

// Create lưu một danh mục (Category) mới vào database.
func (r *CategoryRepository) Create(category *Category) error {
	return r.Save(category)
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
	return r.Save(category)
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

// Save thực hiện lưu (Create nếu ID=0, Update nếu ID>0) cho Sản phẩm.
func (r *ProductRepository) Save(product *Product) error {
	return r.db.Save(product).Error
}

// Create lưu một sản phẩm (Product) mới vào database.
func (r *ProductRepository) Create(product *Product) error {
	return r.Save(product)
}

// FindAllByStoreID lấy danh sách sản phẩm thuộc storeID.
func (r *ProductRepository) FindAllByStoreID(storeID uint) ([]Product, error) {
	var products []Product
	err := r.db.
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("categories.store_id = ?", storeID).
		Preload("Category").
		Find(&products).
		Error
	return products, err
}

// FindByIDAndStoreID tìm 1 sản phẩm theo ID và storeID.
func (r *ProductRepository) FindByIDAndStoreID(id uint, storeID uint) (*Product, error) {
	var product Product
	err := r.db.
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.id = ? AND categories.store_id = ?", id, storeID).
		Preload("Category").
		First(&product).
		Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update cập nhật thông tin sản phẩm.
func (r *ProductRepository) Update(product *Product) error {
	return r.Save(product)
}

// Delete xóa sản phẩm theo ID.
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

