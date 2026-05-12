package main

import (
	"fmt"
	"go-inventory-backend/internal/database"
	"go-inventory-backend/internal/handlers"
	"net/http"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	http.HandleFunc("/products", handlers.ProductsHandler(db))
	http.HandleFunc("/products/", handlers.ProductsHandler(db))

	fmt.Println("Server Running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
