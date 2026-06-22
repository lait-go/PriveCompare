package services

import (
	"errors"
	"time"

	"pricecompare/internal/config"
	"pricecompare/internal/models"
	"pricecompare/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JWTClaims struct {
	UserID uint            `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(cfg.JWTSecret),
		jwtExpiry: time.Duration(cfg.JWTExpiryHours) * time.Hour,
	}
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(email, password, name string) (*AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing.ID > 0 {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Login(email, password string) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type ProductService struct {
	productRepo *repositories.ProductRepository
	priceRepo   *repositories.PriceRepository
}

func NewProductService(productRepo *repositories.ProductRepository, priceRepo *repositories.PriceRepository) *ProductService {
	return &ProductService{productRepo: productRepo, priceRepo: priceRepo}
}

func (s *ProductService) GetProducts(filter repositories.ProductFilter) ([]models.Product, int64, error) {
	return s.productRepo.FindAll(filter)
}

func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) GetPopular(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 8
	}
	return s.productRepo.GetPopular(limit)
}

func (s *ProductService) GetBestDeals(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 8
	}
	return s.productRepo.GetBestDeals(limit)
}

func (s *ProductService) GetRecentPriceChanges(limit int) ([]models.ProductPrice, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.productRepo.GetRecentPriceChanges(limit)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(cat *models.Category) error {
	return s.repo.Create(cat)
}

func (s *CategoryService) Update(cat *models.Category) error {
	return s.repo.Update(cat)
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

type StoreService struct {
	repo *repositories.StoreRepository
}

func NewStoreService(repo *repositories.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}

func (s *StoreService) GetAll() ([]models.Store, error) {
	return s.repo.FindAll()
}

func (s *StoreService) GetByID(id uint) (*models.Store, error) {
	return s.repo.FindByID(id)
}

func (s *StoreService) Create(store *models.Store) error {
	return s.repo.Create(store)
}

func (s *StoreService) Update(store *models.Store) error {
	return s.repo.Update(store)
}

func (s *StoreService) Delete(id uint) error {
	return s.repo.Delete(id)
}

type PriceService struct {
	repo *repositories.PriceRepository
}

func NewPriceService(repo *repositories.PriceRepository) *PriceService {
	return &PriceService{repo: repo}
}

func (s *PriceService) GetByProductID(productID uint) ([]models.ProductPrice, error) {
	return s.repo.FindByProductID(productID)
}

func (s *PriceService) Create(price *models.ProductPrice) error {
	return s.repo.Create(price)
}

func (s *PriceService) Update(price *models.ProductPrice) error {
	return s.repo.Update(price)
}

func (s *PriceService) Delete(id uint) error {
	return s.repo.Delete(id)
}

type CartService struct {
	cartRepo  *repositories.CartRepository
	priceRepo *repositories.PriceRepository
}

func NewCartService(cartRepo *repositories.CartRepository, priceRepo *repositories.PriceRepository) *CartService {
	return &CartService{cartRepo: cartRepo, priceRepo: priceRepo}
}

func (s *CartService) GetCart(userID *uint, sessionID string) (*models.Cart, error) {
	if userID != nil {
		return s.cartRepo.FindByUserID(*userID)
	}
	return s.cartRepo.FindBySessionID(sessionID)
}

func (s *CartService) AddItem(userID *uint, sessionID string, productID, storeID uint, quantity int) (*models.Cart, error) {
	cart, err := s.GetCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	_, err = s.priceRepo.FindByProductAndStore(productID, storeID)
	if err != nil {
		return nil, errors.New("price not found for this product and store")
	}

	item := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		StoreID:   storeID,
		Quantity:  quantity,
	}

	if err := s.cartRepo.AddItem(item); err != nil {
		return nil, err
	}

	return s.GetCart(userID, sessionID)
}

func (s *CartService) UpdateQuantity(userID *uint, sessionID string, itemID uint, quantity int) (*models.Cart, error) {
	cart, err := s.GetCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	item, err := s.cartRepo.GetCartItem(itemID)
	if err != nil {
		return nil, err
	}
	if item.CartID != cart.ID {
		return nil, errors.New("item not in cart")
	}

	if err := s.cartRepo.UpdateItemQuantity(itemID, quantity); err != nil {
		return nil, err
	}

	return s.GetCart(userID, sessionID)
}

func (s *CartService) RemoveItem(userID *uint, sessionID string, itemID uint) (*models.Cart, error) {
	cart, err := s.GetCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	item, err := s.cartRepo.GetCartItem(itemID)
	if err != nil {
		return nil, err
	}
	if item.CartID != cart.ID {
		return nil, errors.New("item not in cart")
	}

	if err := s.cartRepo.RemoveItem(itemID); err != nil {
		return nil, err
	}

	return s.GetCart(userID, sessionID)
}

func (s *CartService) ClearCart(userID *uint, sessionID string) error {
	cart, err := s.GetCart(userID, sessionID)
	if err != nil {
		return err
	}
	return s.cartRepo.ClearCart(cart.ID)
}

type OrderService struct {
	orderRepo *repositories.OrderRepository
	cartRepo  *repositories.CartRepository
	priceRepo *repositories.PriceRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, cartRepo *repositories.CartRepository, priceRepo *repositories.PriceRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, cartRepo: cartRepo, priceRepo: priceRepo}
}

type CreateOrderRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	Comment string `json:"comment"`
}

func (s *OrderService) CreateOrder(userID uint, req CreateOrderRequest) (*models.Order, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var total float64
	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		price, err := s.priceRepo.FindByProductAndStore(item.ProductID, item.StoreID)
		if err != nil {
			return nil, err
		}

		effectivePrice := price.Price
		if price.DiscountPrice != nil {
			effectivePrice = *price.DiscountPrice
		}

		total += effectivePrice * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			StoreID:   item.StoreID,
			Quantity:  item.Quantity,
			Price:     effectivePrice,
		})
	}

	order := &models.Order{
		UserID:  userID,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
		Comment: req.Comment,
		Status:  models.OrderStatusPending,
		Total:   total,
		Items:   orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	if err := s.cartRepo.ClearCart(cart.ID); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(order.ID)
}

func (s *OrderService) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *OrderService) UpdateOrderStatus(id uint, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func EffectivePrice(price models.ProductPrice) float64 {
	if price.DiscountPrice != nil {
		return *price.DiscountPrice
	}
	return price.Price
}

func MinPrice(prices []models.ProductPrice) float64 {
	if len(prices) == 0 {
		return 0
	}
	min := EffectivePrice(prices[0])
	for _, p := range prices[1:] {
		ep := EffectivePrice(p)
		if ep < min {
			min = ep
		}
	}
	return min
}
