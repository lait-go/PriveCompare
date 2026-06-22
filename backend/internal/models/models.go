package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Role         UserRole       `gorm:"size:20;default:user" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Image       string         `gorm:"size:500" json:"image"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Store struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Logo        string         `gorm:"size:500" json:"logo"`
	Rating      float64        `gorm:"default:4.0" json:"rating"`
	Description string         `gorm:"type:text" json:"description"`
	DeliveryCost float64       `gorm:"default:0" json:"delivery_cost"`
	MinOrder    float64        `gorm:"default:0" json:"min_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:500;not null;index" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	CategoryID   uint           `gorm:"index;not null" json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Image        string         `gorm:"size:500" json:"image"`
	Unit         string         `gorm:"size:50" json:"unit"`
	Brand        string         `gorm:"size:255;index" json:"brand"`
	WeightVolume string         `gorm:"size:100" json:"weight_volume"`
	Barcode      string         `gorm:"size:50;index" json:"barcode,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Prices       []ProductPrice `gorm:"foreignKey:ProductID" json:"prices,omitempty"`
}

type ProductPrice struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ProductID       uint           `gorm:"index;not null" json:"product_id"`
	StoreID         uint           `gorm:"index;not null" json:"store_id"`
	Store           Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Price           float64        `gorm:"not null" json:"price"`
	DiscountPrice   *float64       `json:"discount_price,omitempty"`
	DiscountPercent *float64       `json:"discount_percent,omitempty"`
	Popularity      int            `gorm:"default:0" json:"popularity"`
	InStock         bool           `gorm:"default:true" json:"in_stock"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    *uint          `gorm:"index" json:"user_id,omitempty"`
	SessionID string         `gorm:"size:100;index" json:"session_id,omitempty"`
	Items     []CartItem     `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CartID    uint           `gorm:"index;not null" json:"cart_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	StoreID   uint           `gorm:"index;not null" json:"store_id"`
	Store     Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Quantity  int            `gorm:"default:1;not null" json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Phone     string         `gorm:"size:50;not null" json:"phone"`
	Address   string         `gorm:"type:text;not null" json:"address"`
	Comment   string         `gorm:"type:text" json:"comment"`
	Status    OrderStatus    `gorm:"size:50;default:pending" json:"status"`
	Total     float64        `gorm:"not null" json:"total"`
	Items     []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"index;not null" json:"order_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	StoreID   uint           `gorm:"index;not null" json:"store_id"`
	Store     Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"not null" json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
