package product

import (
	"database/sql"
	"fmt"
	"strings"
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

func (s *Store) GetProductById(id int) (types.Product, error) {
	product := new(types.Product)
	err := s.db.QueryRow("SELECT * FROM products where id = ?", id).Scan(productAllFieldsScanner(product))

	if err != nil {
		return types.Product{}, err
	}

	return *product, nil
}

func (s *Store) GetProducts(limit, offset int) ([]types.Product, int, error) {
	rows, err := s.db.Query("SELECT * FROM products LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		prod, err := scanRowsIntoProducts(rows)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, *prod)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%v)", placeholders)

	args := make([]interface{}, len(productIDs))
	for i, val := range productIDs {
		args[i] = val
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		prod, err := scanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *prod)
	}

	return products, nil
}

func (s *Store) CreateProduct(payload types.ProductCreatePayload) (*types.Product, error) {
	var query = "INSERT INTO products (name,description,image,price,quantity) VALUES(?,?,?,?,?)"
	result, err := s.db.Exec(query, strings.TrimSpace(payload.Name), strings.TrimSpace(payload.Description), payload.Image, payload.Price, payload.Quantity)
	if err != nil {
		return nil, err
	}

	prodId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var getProdQuery = "SELECT * FROM products WHERE id = ?"
	row := s.db.QueryRow(getProdQuery, prodId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	createdProd, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}

	return createdProd, nil
}

func (s *Store) UpdateProduct(id int, payload types.ProductUpdatePayload) (*types.Product, error) {
	query := "UPDATE products SET"
	updates, args := handleProductFields(payload)

	query += " " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("SELECT * FROM products WHERE id = ?", id)
	prodAfterUpdate, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}

	return prodAfterUpdate, nil
}

func (s *Store) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := s.db.Exec(query, id)

	if err != nil {
		return err
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if RowsAffected == 0 {
		return fmt.Errorf("no product was found for id %v", id)
	}

	return nil
}

func scanRowsIntoProducts(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(productAllFieldsScanner(product))

	if err != nil {
		return &types.Product{}, err
	}

	return product, nil
}

func scanRowIntoProduct(row *sql.Row) (*types.Product, error) {
	product := new(types.Product)

	err := row.Scan(productAllFieldsScanner(product))

	if err != nil {
		return &types.Product{}, err
	}

	return product, nil
}

func productAllFieldsScanner(product *types.Product) (*int, *string, *string, *string, *float64, *int, *time.Time, *time.Time) {
	return &product.ID,
		&product.Name,
		&product.Description,
		&product.Image, &product.Price,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt
}

func handleProductFields(payload types.ProductUpdatePayload) ([]string, []any) {
	var updates []string
	var args []any

	if payload.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, strings.TrimSpace(payload.Name))
	}

	if payload.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, strings.TrimSpace(payload.Description))
	}

	if payload.Image != "" {
		updates = append(updates, "image = ?")
		args = append(args, payload.Image)
	}

	if payload.Price != 0 {
		updates = append(updates, "price = ?")
		args = append(args, payload.Price)
	}

	if payload.Quantity != 0 {
		updates = append(updates, "quantity = ?")
		args = append(args, payload.Quantity)
	}

	return updates, args
}
