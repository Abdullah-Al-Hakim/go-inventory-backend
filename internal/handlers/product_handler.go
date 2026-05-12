package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-inventory-backend/internal/models"
	"go-inventory-backend/internal/services"
	"net/http"
	"strconv"
	"strings"
)

func ProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("METHOD", r.Method)

		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			products, err := services.GetProducts(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(products)

		case http.MethodPost:
			var product models.Product
			if err :=
				json.NewDecoder(r.Body).Decode(&product); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			if product.Name == "" || product.Stock < 0 {
				http.Error(w, "Invalid product data", http.StatusBadRequest)
				return
			}

			if err := services.CreateProduct(db, product); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "product created succesfully",
			})

		case http.MethodPut:
			idStr := strings.TrimPrefix(r.URL.Path, "/products/")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}
			var payload struct {
				Stock int `json:"stock"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			if payload.Stock < 0 {
				http.Error(w, "Stock must be >=0", http.StatusBadRequest)
				return
			}
			if err := services.UpdateProductStock(db, id, payload.Stock); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{
				"message": "stock updated successfully",
			})

		case http.MethodDelete:
			idStr := strings.TrimPrefix(r.URL.Path, "/products/")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}
			if err := services.DeleteProduct(db, id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted successfully",
			})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
