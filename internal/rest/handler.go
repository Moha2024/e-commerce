package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"e-commerce/internal/domain"
)



func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок, что отправляем JSON
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(domain.Products)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, product := range domain.Products {
		if product.ID == id{
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input ProductRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	defer r.Body.Close()

	for i, product := range domain.Products {
		if product.ID == id {
			domain.Products[i].Name = input.Name
			domain.Products[i].Price = input.Price

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(domain.Products[i])
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request){
	var input ProductRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newProduct := domain.Product{
		ID: domain.Products[len(domain.Products)-1].ID + 1,
		Name: input.Name,
		Price: input.Price,
	}

	domain.Products = append(domain.Products, newProduct)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}