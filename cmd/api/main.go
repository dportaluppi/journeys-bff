package main

import (
	"github.com/dportaluppi/journeys-bff/pkg/journey"
	"github.com/dportaluppi/journeys-bff/pkg/product"
	"github.com/dportaluppi/journeys-bff/pkg/segment"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	const headlessURL = "https://storefront-admin.yalochat.dev/v3/admin/storefronts"
	productRepo := product.NewGraphQLRepo(headlessURL)
	productHandler := product.NewHandler(productRepo)
	router := gin.Default()
	router.GET("/products", productHandler.SearchProducts)

	const audienceToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJZYWxvQXBpS2V5IiwieWlkIjoiMjUxIiwieXQiOiJjb25zdW1lciJ9.PKgWMQnr9oJVvtVOvsTPK95HzfPVA6pFLMqtL8cDJEk"
	const audienceURL = "https://api2-staging.yalochat.com/customers-api/v1/filter-ruleset"
	audienceRepo := segment.NewHTTPRepo(audienceURL, audienceToken)
	audienceHandler := segment.NewHandler(audienceRepo)
	router.GET("/segments", audienceHandler.GetSegments)

	const journeysURL = "http://localhost:8060/accounts/"
	const journeysToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0VGk5amthVldha0J3Q2hvYnRHOTRmWEpHUjZ1NjRhTiJ9.YWuWwltEdyqtpTor_AQjC_-YmdcqTgw2WvQcKpEOfwo"
	journeyRepo := journey.NewHTTPRepo(journeysURL, journeysToken)
	journeyService := journey.NewService(journeyRepo, productRepo, audienceRepo)
	journeyHandler := journey.NewHandler(journeyService)
	router.POST("/accounts/:accountId/journeys", journeyHandler.CreateJourney)
	router.GET("/accounts/:accountId/journeys", journeyHandler.GetJourneys)
	router.GET("/accounts/:accountId/journeys/:journeyId", journeyHandler.GetJourneyByID)

	log.Fatal(router.Run(":8061"))
}
