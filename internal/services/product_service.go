package services

import (
	"database/sql"
	"errors"
	"fmt"
	"go-inventory-backend/internal/models"
)

func GetProducts(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query("SELECT id,name,stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func CreateProduct(db *sql.DB, product models.Product) error {
	_, err := db.Exec("INSERT INTO products (name,stock) Values (?,?)",
		product.Name, product.Stock)
	return err
}

func UpdateProductStock(db *sql.DB, id int, stock int) error {
	query := "UPDATE products SET stock = ? WHERE id = ?"
	result, err := db.Exec(query, stock, id)

	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func DeleteProduct(db *sql.DB, id int) error {
	query := "DELETE FROM products WHERE id =?"
	result, err := db.Exec(query, id)

	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func GetProductsPagination(db *sql.DB, limit, offset int) ([]models.Product, error) {
	query := fmt.Sprintf("SELECT id, name, stock FROM products LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil

}
