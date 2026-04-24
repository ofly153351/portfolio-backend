package http

import (
	"context"
	"portfolio-backend/internal/config"
	"portfolio-backend/internal/database"
	authmodule "portfolio-backend/internal/modules/auth"
	chatmodule "portfolio-backend/internal/modules/chat"
	contentmodule "portfolio-backend/internal/modules/content"
	healthmodule "portfolio-backend/internal/modules/health"
	publicauthmodule "portfolio-backend/internal/modules/publicauth"
	uploadmodule "portfolio-backend/internal/modules/upload"
	"time"
)

type Dependencies struct {
	AuthHandler    *authmodule.Handler
	HealthHandler  *healthmodule.Handler
	ContentHandler *contentmodule.Handler
	UploadHandler  *uploadmodule.Handler
	ChatHandler    *chatmodule.Handler
	PublicHandler  *publicauthmodule.Handler
}

func NewDependencies(cfg config.Config) *Dependencies {
	authService, err := authmodule.NewService(cfg)
	if err != nil {
		panic(err)
	}
	authHandler := authmodule.NewHandler(authService)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()
	mongoClient, err := database.NewMongo(dbCtx, cfg.MongoURI)
	if err != nil {
		panic(err)
	}
	contentRepo := contentmodule.NewMongoRepository(mongoClient.Database(cfg.MongoDB))
	contentService := contentmodule.NewService(contentRepo)
	publicAuthService := publicauthmodule.NewService(
		cfg.PublicTokenSecret,
		time.Duration(cfg.PublicTokenTTL)*time.Second,
	)
	publicAuthHandler := publicauthmodule.NewHandler(publicAuthService)
	contentHandler := contentmodule.NewHandler(contentService, authHandler, publicAuthHandler)

	uploadService, err := uploadmodule.NewService(cfg)
	if err != nil {
		panic(err)
	}
	uploadHandler := uploadmodule.NewHandler(uploadService, authHandler)

	chatRepo := chatmodule.NewAIServiceRepositoryWithTimeout(
		cfg.AIServiceURL,
		time.Duration(cfg.AIServiceTimeout)*time.Second,
	)
	chatService := chatmodule.NewService(chatRepo)

	return &Dependencies{
		AuthHandler:    authHandler,
		HealthHandler:  healthmodule.NewHandler(),
		ContentHandler: contentHandler,
		UploadHandler:  uploadHandler,
		ChatHandler:    chatmodule.NewHandler(chatService, publicAuthHandler),
		PublicHandler:  publicAuthHandler,
	}
}
