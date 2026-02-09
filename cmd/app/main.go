package main

import (
	"e-commerce/internal/rest"
	"net/http"
)

func main() {
	http.HandleFunc("GET /products", rest.GetProductsHandler)
	http.HandleFunc("GET /products/{id}", rest.GetProductHandler)
	http.HandleFunc("PUT /products/{id}", rest.UpdateProductsHandler)
	http.HandleFunc("POST /products", rest.CreateProductHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
