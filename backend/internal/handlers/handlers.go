package handlers

import (
	"net/http"
	"strconv"

	"pricecompare/internal/models"
	"pricecompare/internal/repositories"
	"pricecompare/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	auth     *services.AuthService
	product  *services.ProductService
	category *services.CategoryService
	store    *services.StoreService
	price    *services.PriceService
	cart     *services.CartService
	order    *services.OrderService
}

func NewHandlers(
	auth *services.AuthService,
	product *services.ProductService,
	category *services.CategoryService,
	store *services.StoreService,
	price *services.PriceService,
	cart *services.CartService,
	order *services.OrderService,
) *Handlers {
	return &Handlers{auth, product, category, store, price, cart, order}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body registerRequest true "Registration data"
// @Success 201 {object} services.AuthResponse
// @Router /auth/register [post]
func (h *Handlers) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.auth.Register(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Login credentials"
// @Success 200 {object} services.AuthResponse
// @Router /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Router /users/me [get]
func (h *Handlers) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")
	c.JSON(http.StatusOK, gin.H{
		"id":   userID,
		"role": role,
	})
}

// GetCategories godoc
// @Summary Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *Handlers) GetCategories(c *gin.Context) {
	cats, err := h.category.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// GetCategory godoc
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Router /categories/{id} [get]
func (h *Handlers) GetCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cat, err := h.category.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handlers) CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.category.Create(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *Handlers) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = uint(id)
	if err := h.category.Update(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handlers) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.category.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetStores godoc
// @Summary Get all stores
// @Tags stores
// @Produce json
// @Success 200 {array} models.Store
// @Router /stores [get]
func (h *Handlers) GetStores(c *gin.Context) {
	stores, err := h.store.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stores)
}

func (h *Handlers) GetStore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	store, err := h.store.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (h *Handlers) CreateStore(c *gin.Context) {
	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.Create(&store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, store)
}

func (h *Handlers) UpdateStore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store.ID = uint(id)
	if err := h.store.Update(&store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (h *Handlers) DeleteStore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.store.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetProducts godoc
// @Summary Search and filter products
// @Tags products
// @Produce json
// @Param q query string false "Search query"
// @Param category_id query int false "Category ID"
// @Param store_id query int false "Store ID"
// @Param brand query string false "Brand"
// @Param min_price query number false "Min price"
// @Param max_price query number false "Max price"
// @Param has_discount query bool false "Has discount"
// @Param sort query string false "Sort: price_asc, price_desc, name, popularity, rating"
// @Param page query int false "Page"
// @Param per_page query int false "Per page"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func (h *Handlers) GetProducts(c *gin.Context) {
	filter := repositories.ProductFilter{
		Query:   c.Query("q"),
		Brand:   c.Query("brand"),
		SortBy:  c.DefaultQuery("sort", ""),
		Page:    parseInt(c.Query("page"), 1),
		PerPage: parseInt(c.Query("per_page"), 20),
	}

	if v := c.Query("category_id"); v != "" {
		id := uint(parseInt(v, 0))
		filter.CategoryID = &id
	}
	if v := c.Query("store_id"); v != "" {
		id := uint(parseInt(v, 0))
		filter.StoreID = &id
	}
	if v := c.Query("min_price"); v != "" {
		f, _ := strconv.ParseFloat(v, 64)
		filter.MinPrice = &f
	}
	if v := c.Query("max_price"); v != "" {
		f, _ := strconv.ParseFloat(v, 64)
		filter.MaxPrice = &f
	}
	if v := c.Query("has_discount"); v == "true" {
		b := true
		filter.HasDiscount = &b
	}

	products, total, err := h.product.GetProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  filter.Page,
		"per_page": filter.PerPage,
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func (h *Handlers) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	product, err := h.product.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetHome godoc
// @Summary Get home page data
// @Tags home
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /home [get]
func (h *Handlers) GetHome(c *gin.Context) {
	popular, _ := h.product.GetPopular(8)
	deals, _ := h.product.GetBestDeals(8)
	recent, _ := h.product.GetRecentPriceChanges(10)
	categories, _ := h.category.GetAll()
	stores, _ := h.store.GetAll()

	c.JSON(http.StatusOK, gin.H{
		"popular_products":    popular,
		"best_deals":          deals,
		"recent_price_changes": recent,
		"categories":          categories,
		"stores":              stores,
	})
}

func (h *Handlers) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.product.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handlers) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = uint(id)
	if err := h.product.UpdateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handlers) DeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.product.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handlers) GetPrices(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 64)
	prices, err := h.price.GetByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

func (h *Handlers) CreatePrice(c *gin.Context) {
	var price models.ProductPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.price.Create(&price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, price)
}

func (h *Handlers) UpdatePrice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var price models.ProductPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	price.ID = uint(id)
	if err := h.price.Update(&price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, price)
}

func (h *Handlers) DeletePrice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.price.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type addCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	StoreID   uint `json:"store_id" binding:"required"`
	Quantity  int  `json:"quantity"`
}

// GetCart godoc
// @Summary Get cart
// @Tags cart
// @Produce json
// @Success 200 {object} models.Cart
// @Router /cart [get]
func (h *Handlers) GetCart(c *gin.Context) {
	userID, sessionID := getCartIdentity(c)
	cart, err := h.cart.GetCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// AddToCart godoc
// @Summary Add item to cart
// @Tags cart
// @Accept json
// @Produce json
// @Param body body addCartRequest true "Cart item"
// @Success 200 {object} models.Cart
// @Router /cart [post]
func (h *Handlers) AddToCart(c *gin.Context) {
	var req addCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	userID, sessionID := getCartIdentity(c)
	cart, err := h.cart.AddItem(userID, sessionID, req.ProductID, req.StoreID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handlers) UpdateCartItem(c *gin.Context) {
	itemID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	quantity := parseInt(c.Query("quantity"), 1)
	userID, sessionID := getCartIdentity(c)
	cart, err := h.cart.UpdateQuantity(userID, sessionID, uint(itemID), quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handlers) RemoveCartItem(c *gin.Context) {
	itemID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, sessionID := getCartIdentity(c)
	cart, err := h.cart.RemoveItem(userID, sessionID, uint(itemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handlers) ClearCart(c *gin.Context) {
	userID, sessionID := getCartIdentity(c)
	if err := h.cart.ClearCart(userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

// CreateOrder godoc
// @Summary Create order
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body services.CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Router /orders [post]
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)
	order, err := h.order.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"order":   order,
		"message": "Ваш заказ успешно оформлен и передан в доставку.",
	})
}

// GetMyOrders godoc
// @Summary Get user orders
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Router /orders/my [get]
func (h *Handlers) GetMyOrders(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)
	orders, err := h.order.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handlers) GetAllOrders(c *gin.Context) {
	orders, err := h.order.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

type updateStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

func (h *Handlers) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.order.UpdateOrderStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func getCartIdentity(c *gin.Context) (*uint, string) {
	if uid, exists := c.Get("userID"); exists {
		id := uid.(uint)
		return &id, ""
	}
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = c.ClientIP() + "-guest"
	}
	return nil, sessionID
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
