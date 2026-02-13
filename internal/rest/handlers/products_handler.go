package handlers

import (
	"e-commerce/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

func CreateProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateProductRequest

		if err := c.ShouldBindJSON(&input); err != nil { // подставляет совпадающие поля JSON в ProductRequest
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := repository.CreateProduct(pool, input.Name, input.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

// func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Устанавливаем заголовок, что отправляем JSON
// 	w.Header().Set("Content-Type", "application/json")

// 	json.NewEncoder(w).Encode(domain.Products)
// }

// func GetProductHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}
// 	for _, product := range domain.Products {
// 		if product.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(product)
// 			return
// 		}
// 	}
// 	http.Error(w, "Product not found", http.StatusNotFound)
// }

// func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var input ProductRequest
// 	err = json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	defer r.Body.Close()

// 	for i, product := range domain.Products {
// 		if product.ID == id {
// 			domain.Products[i].Name = input.Name
// 			domain.Products[i].Price = input.Price

// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(domain.Products[i])
// 			return
// 		}
// 	}

// 	http.Error(w, "Product not found", http.StatusNotFound)
// }

// // func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
// // 	var input ProductRequest
// // 	err := json.NewDecoder(r.Body).Decode(&input)
// // 	if err != nil {
// // 		http.Error(w, "Invalid json", http.StatusBadRequest)
// // 		return
// // 	}
// // 	defer r.Body.Close()

// // 	newProduct := domain.Product{
// // 		ID:    domain.Products[len(domain.Products)-1].ID + 1,
// // 		Name:  input.Name,
// // 		Price: input.Price,
// // 	}

// // 	domain.Products = append(domain.Products, newProduct)
// // 	w.Header().Set("Content-Type", "application/json")
// // 	w.WriteHeader(http.StatusCreated)
// // 	json.NewEncoder(w).Encode(newProduct)
// // }

// func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid id", http.StatusBadRequest)
// 		return
// 	}

// 	for i, product := range domain.Products {
// 		if product.ID == id {
// 			domain.Products = append(domain.Products[:i], domain.Products[i+1:]...)

// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 	}
// 	http.Error(w, "Product not found", http.StatusNotFound)
// }
