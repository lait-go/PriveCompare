package repositories

import (
	"pricecompare/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *CategoryRepository) Create(cat *models.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *models.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) FindAll() ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Order("name ASC").Find(&stores).Error
	return stores, err
}

func (r *StoreRepository) FindByID(id uint) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, id).Error
	return &store, err
}

func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) Delete(id uint) error {
	return r.db.Delete(&models.Store{}, id).Error
}

type ProductFilter struct {
	Query      string
	CategoryID *uint
	StoreID    *uint
	Brand      string
	MinPrice   *float64
	MaxPrice   *float64
	HasDiscount *bool
	SortBy     string
	Page       int
	PerPage    int
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Prices.Store").First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) FindAll(filter ProductFilter) ([]models.Product, int64, error) {
	query := r.db.Model(&models.Product{}).Preload("Category")

	if filter.Query != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+filter.Query+"%")
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Brand != "" {
		query = query.Where("LOWER(brand) LIKE LOWER(?)", "%"+filter.Brand+"%")
	}

	if filter.StoreID != nil || filter.MinPrice != nil || filter.MaxPrice != nil || filter.HasDiscount != nil {
		subQuery := r.db.Model(&models.ProductPrice{}).Select("DISTINCT product_id")
		if filter.StoreID != nil {
			subQuery = subQuery.Where("store_id = ?", *filter.StoreID)
		}
		if filter.MinPrice != nil {
			subQuery = subQuery.Where("COALESCE(discount_price, price) >= ?", *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			subQuery = subQuery.Where("COALESCE(discount_price, price) <= ?", *filter.MaxPrice)
		}
		if filter.HasDiscount != nil && *filter.HasDiscount {
			subQuery = subQuery.Where("discount_price IS NOT NULL")
		}
		query = query.Where("id IN (?)", subQuery)
	}

	var total int64
	query.Count(&total)

	switch filter.SortBy {
	case "price_asc":
		query = query.Joins("LEFT JOIN product_prices pp ON pp.product_id = products.id AND pp.deleted_at IS NULL").
			Group("products.id").
			Order("MIN(COALESCE(pp.discount_price, pp.price)) ASC")
	case "price_desc":
		query = query.Joins("LEFT JOIN product_prices pp ON pp.product_id = products.id AND pp.deleted_at IS NULL").
			Group("products.id").
			Order("MIN(COALESCE(pp.discount_price, pp.price)) DESC")
	case "name":
		query = query.Order("name ASC")
	case "popularity":
		query = query.Joins("LEFT JOIN product_prices pp ON pp.product_id = products.id AND pp.deleted_at IS NULL").
			Group("products.id").
			Order("MAX(pp.popularity) DESC")
	case "rating":
		query = query.Joins("LEFT JOIN product_prices pp ON pp.product_id = products.id AND pp.deleted_at IS NULL").
			Joins("LEFT JOIN stores s ON s.id = pp.store_id AND s.deleted_at IS NULL").
			Group("products.id").
			Order("MAX(s.rating) DESC")
	default:
		query = query.Order("created_at DESC")
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	var products []models.Product
	err := query.Preload("Prices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Store")
	}).Offset(offset).Limit(filter.PerPage).Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) GetPopular(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Preload("Prices.Store").
		Joins("JOIN product_prices pp ON pp.product_id = products.id AND pp.deleted_at IS NULL").
		Group("products.id").
		Order("MAX(pp.popularity) DESC").
		Limit(limit).Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetBestDeals(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Preload("Prices", func(db *gorm.DB) *gorm.DB {
		return db.Where("discount_price IS NOT NULL").Preload("Store")
	}).
		Joins("JOIN product_prices pp ON pp.product_id = products.id AND pp.discount_price IS NOT NULL AND pp.deleted_at IS NULL").
		Group("products.id").
		Order("MAX(pp.discount_percent) DESC").
		Limit(limit).Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetRecentPriceChanges(limit int) ([]models.ProductPrice, error) {
	var prices []models.ProductPrice
	err := r.db.Preload("Product").Preload("Store").
		Order("updated_at DESC").
		Limit(limit).Find(&prices).Error
	return prices, err
}

type PriceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

func (r *PriceRepository) FindByID(id uint) (*models.ProductPrice, error) {
	var price models.ProductPrice
	err := r.db.Preload("Store").Preload("Product").First(&price, id).Error
	return &price, err
}

func (r *PriceRepository) FindByProductAndStore(productID, storeID uint) (*models.ProductPrice, error) {
	var price models.ProductPrice
	err := r.db.Where("product_id = ? AND store_id = ?", productID, storeID).First(&price).Error
	return &price, err
}

func (r *PriceRepository) Create(price *models.ProductPrice) error {
	return r.db.Create(price).Error
}

func (r *PriceRepository) Update(price *models.ProductPrice) error {
	return r.db.Save(price).Error
}

func (r *PriceRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductPrice{}, id).Error
}

func (r *PriceRepository) FindByProductID(productID uint) ([]models.ProductPrice, error) {
	var prices []models.ProductPrice
	err := r.db.Preload("Store").Where("product_id = ?", productID).Find(&prices).Error
	return prices, err
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func cartPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("Items.Product.Prices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Store")
	}).Preload("Items.Store")
}

func (r *CartRepository) FindByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := cartPreloads(r.db).Where("user_id = ?", userID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{UserID: &userID}
		if createErr := r.db.Create(&cart).Error; createErr != nil {
			return nil, createErr
		}
		return &cart, nil
	}
	return &cart, err
}

func (r *CartRepository) FindBySessionID(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := cartPreloads(r.db).Where("session_id = ?", sessionID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{SessionID: sessionID}
		if createErr := r.db.Create(&cart).Error; createErr != nil {
			return nil, createErr
		}
		return &cart, nil
	}
	return &cart, err
}

func (r *CartRepository) AddItem(item *models.CartItem) error {
	var existing models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ? AND store_id = ?",
		item.CartID, item.ProductID, item.StoreID).First(&existing).Error
	if err == nil {
		existing.Quantity += item.Quantity
		return r.db.Save(&existing).Error
	}
	return r.db.Create(item).Error
}

func (r *CartRepository) UpdateItemQuantity(itemID uint, quantity int) error {
	if quantity <= 0 {
		return r.db.Delete(&models.CartItem{}, itemID).Error
	}
	return r.db.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(itemID uint) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *CartRepository) GetCartItem(itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.First(&item, itemID).Error
	return &item, err
}

func (r *CartRepository) MergeCarts(sessionID string, userID uint) error {
	sessionCart, err := r.FindBySessionID(sessionID)
	if err != nil {
		return err
	}
	userCart, err := r.FindByUserID(userID)
	if err != nil {
		return err
	}

	for _, item := range sessionCart.Items {
		item.CartID = userCart.ID
		item.ID = 0
		if addErr := r.AddItem(&item); addErr != nil {
			return addErr
		}
	}

	return r.db.Where("id = ?", sessionCart.ID).Delete(&models.Cart{}).Error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items.Product").Preload("Items.Store").
		Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items.Product").Preload("Items.Store").Preload("User").
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items.Product").Preload("Items.Store").Preload("User").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
