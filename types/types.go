package types

import (
	"time"
)

// Product types

type ProductStore interface {
	GetProductById(id int) (Product, error)
	GetProducts(page, offset int) ([]Product, int, error)
	GetProductsByID(productIDs []int) ([]Product, error)
	CreateProduct(payload ProductCreatePayload) (*Product, error)
	UpdateProduct(id int, payload ProductUpdatePayload) (*Product, error)
	DeleteProduct(id int) error
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProductCreatePayload struct {
	Name        string  `json:"name" validate:"required,min=3,max=256"`
	Description string  `json:"description" validate:"required,max=3000"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Quantity    int     `json:"quantity" validate:"required,gte=0"`
}

type ProductUpdatePayload struct {
	Name        string  `json:"name" validate:"min=3,max=256"`
	Description string  `json:"description" validate:"min=4,max=3000"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"gt=0"`
	Quantity    int     `json:"quantity" validate:"gte=0"`
}

// User types

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=64"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user User) error
}

// order types

type Order struct {
	ID        int `json:"id"`
	UserID    int `json:"userId"`
	Total     float64 `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderStore interface {
	CreateOrder(order Order) (Order ,error)
	CreateOrderItem(orderItem OrderItem) (OrderItem ,error)
}

// order items types

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderId"`
	ProductID int       `json:"productId"`
	Quantity  int       `json:"quantity" validate:"gte=0"`
	Price     float64   `json:"price" validate:"gte=0"`
	CreatedAt time.Time `json:"createdAt"`
}

// checkout type

type CartCheckoutItem struct {
	ProductID int `json:"productId"`
	Quantity int `json:"quantity" validate:"gte=0"`
}

type CartCheckoutItems struct {
	CartItems []CartCheckoutItem `json:"cartItems" validate:"required"`
}
