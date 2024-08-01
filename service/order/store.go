package order

import (
	"database/sql"
	"time"

	"github.com/mohammadahmadkhader/golang-ecommerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateOrder(order types.Order) (types.Order, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?,?,?,?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return types.Order{}, err
	}

	orderId, err := res.LastInsertId()
	if err != nil {
		return types.Order{}, err
	}

	row := s.db.QueryRow("SELECT * FROM orders WHERE id = ?", orderId)
	newOrder, err := scanRowIntoOrder(row)
	if err != nil {
		return types.Order{}, err
	}

	return *newOrder, nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem, sd *sql.Tx) (types.OrderItem, error) {
	res, err := s.db.Exec("INSERT INTO orderItems (orderId, productId, quantity, price) VALUES (?,?,?,?)",
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return types.OrderItem{}, err
	}

	orderItemId, err := res.LastInsertId()
	if err != nil {
		return types.OrderItem{}, err
	}

	row := s.db.QueryRow("SELECT * FROM orderItems WHERE id = ?", orderItemId)
	newOrderItem, err := scanRowIntoOrderItem(row)

	return *newOrderItem, err
}

func scanRowIntoOrder(row *sql.Row) (*types.Order, error) {
	order := new(types.Order)
	err := row.Scan(orderAllFieldsScanner(order))
	if err != nil {
		return &types.Order{}, err
	}

	return order, nil
}

func scanRowIntoOrderItem(row *sql.Row) (*types.OrderItem, error) {
	orderItem := new(types.OrderItem)
	err := row.Scan(orderItemAllFieldsScanner(orderItem))
	if err != nil {
		return &types.OrderItem{}, err
	}

	return orderItem, nil
}

func orderAllFieldsScanner(order *types.Order) (*int, *int, *float64, *string, *string, *time.Time) {
	return &order.ID, &order.UserID, &order.Total, &order.Status, &order.Address, &order.CreatedAt
}

func orderItemAllFieldsScanner(orderItem *types.OrderItem) (*int, *int, *int, *int, *float64, *time.Time) {
	return &orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price, &orderItem.CreatedAt
}
