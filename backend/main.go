package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/config"
	"github.com/jwd0526/nutrirecipe/db"
	"github.com/jwd0526/nutrirecipe/handlers"
	"github.com/jwd0526/nutrirecipe/services"
)

func main() {
	cfg := config.Load()

	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	pool, err := db.Connect(context.Background())
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	agentSvc := services.NewAgentService()
	usdaSvc := services.NewUSDAService(pool, cfg.USDAAPIKey)

	r := gin.Default()
	api := r.Group("/api")

	agentH := handlers.NewAgentHandler(agentSvc)
	api.POST("/agent/parse", agentH.Parse)

	usdaH := handlers.NewUSDAHandler(usdaSvc)
	api.POST("/usda/validate", usdaH.Validate)

	recipeH := handlers.NewRecipeHandler(pool)
	api.POST("/recipes", recipeH.Save)
	api.GET("/recipes", recipeH.List)
	api.GET("/recipes/:id", recipeH.Get)

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
