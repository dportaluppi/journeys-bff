package main

import (
	"github.com/dportaluppi/journeys-bff/internal/product"
	domain "github.com/dportaluppi/journeys-bff/pkg/recommendations"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	graphQLURL := "https://storefront-admin.yalochat.dev/v3/admin/storefronts"
	productRepo := product.NewProductRepository(graphQLURL)
	productService := domain.NewSearcher(productRepo)

	handler := product.NewHandler(productService)
	router := gin.Default()
	router.GET("/products", handler.SearchProducts)

	log.Fatal(router.Run(":8061"))
}
