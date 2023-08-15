package main

import (
	"github.com/dportaluppi/journeys-bff/pkg/product"
	"github.com/dportaluppi/journeys-bff/pkg/segment"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	const headlessURL = "https://storefront-admin.yalochat.dev/v3/admin/storefronts"
	productRepo := product.NewGraphQLRepository(headlessURL)

	handler := product.NewHandler(productRepo)
	router := gin.Default()
	router.GET("/products", handler.SearchProducts)

	const audienceToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJZYWxvQXBpS2V5IiwieWlkIjoiMjUxIiwieXQiOiJjb25zdW1lciJ9.PKgWMQnr9oJVvtVOvsTPK95HzfPVA6pFLMqtL8cDJEk"
	const audienceURL = "https://api2-staging.yalochat.com/customers-api/v1/filter-ruleset"

	audienceRepo := segment.NewHTTPAudienceRepo(audienceURL, audienceToken)
	audienceHandler := segment.NewHandler(audienceRepo)
	router.GET("/segments", audienceHandler.GetSegments)

	log.Fatal(router.Run(":8061"))
}
