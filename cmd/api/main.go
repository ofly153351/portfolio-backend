package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"portfolio-backend/internal/app/http"
	"portfolio-backend/internal/config"
	"portfolio-backend/internal/database"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	bootCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	mongoClient, err := database.NewMongo(bootCtx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()
	if err := database.EnsureSchema(bootCtx, mongoClient, cfg.MongoDB); err != nil {
		log.Fatal(err)
	}
	if err := database.MigrateLegacyPortfolioContent(bootCtx, mongoClient, cfg.MongoDB); err != nil {
		log.Fatal(err)
	}

	app := http.NewApp(cfg)

	log.Printf("server starting on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
