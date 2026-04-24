package config

import (
	"os"
	"strconv"
	"strings"
)

// Config stores runtime configuration.
type Config struct {
	Port                 string
	AppName              string
	MongoURI             string
	MongoDB              string
	MinIORootUser        string
	MinIORootPassword    string
	MinIOEndpoint        string
	MinIOBucket          string
	MinIOUseSSL          bool
	MinIOPublicBaseURL   string
	AdminUsername        string
	AdminPassword        string
	AdminPasswordHash    string
	AuthTokenTTLMinutes  int
	AIServiceURL         string
	AIServiceTimeout     int
	PublicTokenSecret    string
	PublicTokenTTL       int
	OpenAIAPIKey         string
	OpenAIBaseURL        string
	OpenAIChatModel      string
	OpenRouterRef        string
	OpenRouterTitle      string
	CORSAllowOrigins     string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool
}

// Load reads config from environment with safe defaults.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "portfolio-backend"
	}

	allowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000"
	}

	allowMethods := os.Getenv("CORS_ALLOW_METHODS")
	if allowMethods == "" {
		allowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	}

	allowHeaders := os.Getenv("CORS_ALLOW_HEADERS")
	if allowHeaders == "" {
		allowHeaders = "Origin,Content-Type,Accept,Authorization,X-Public-Token"
	}
	allowHeaders = ensureHeaderAllowed(allowHeaders, "X-Public-Token")
	allowCredentials := true
	if raw := os.Getenv("CORS_ALLOW_CREDENTIALS"); raw != "" {
		if parsed, err := strconv.ParseBool(raw); err == nil {
			allowCredentials = parsed
		}
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "portfolio"
	}

	minioRootUser := os.Getenv("MINIO_ROOT_USER")
	if minioRootUser == "" {
		minioRootUser = "admin"
	}
	minioRootPassword := os.Getenv("MINIO_ROOT_PASSWORD")
	if minioRootPassword == "" {
		minioRootPassword = "password"
	}
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}
	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "portfolio"
	}
	minioUseSSL := false
	if raw := os.Getenv("MINIO_USE_SSL"); raw != "" {
		if parsed, err := strconv.ParseBool(raw); err == nil {
			minioUseSSL = parsed
		}
	}
	minioPublicBaseURL := os.Getenv("MINIO_PUBLIC_BASE_URL")
	if minioPublicBaseURL == "" {
		scheme := "http"
		if minioUseSSL {
			scheme = "https"
		}
		minioPublicBaseURL = scheme + "://" + minioEndpoint
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}
	adminPasswordHash := os.Getenv("ADMIN_PASSWORD_HASH")
	authTokenTTLMinutes := 720
	if raw := os.Getenv("AUTH_TOKEN_TTL_MINUTES"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			authTokenTTLMinutes = parsed
		}
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}
	aiServiceTimeout := 180
	if raw := os.Getenv("AI_SERVICE_TIMEOUT_SECONDS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			aiServiceTimeout = parsed
		}
	}
	publicTokenSecret := os.Getenv("PUBLIC_TOKEN_SECRET")
	if publicTokenSecret == "" {
		publicTokenSecret = "dev-public-token-secret-change-me"
	}
	publicTokenTTL := 300
	if raw := os.Getenv("PUBLIC_TOKEN_TTL_SECONDS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			publicTokenTTL = parsed
		}
	}

	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	if openAIBaseURL == "" {
		openAIBaseURL = "https://api.openai.com/v1"
	}
	openAIChatModel := os.Getenv("OPENAI_CHAT_MODEL")
	if openAIChatModel == "" {
		openAIChatModel = "gpt-4o-mini"
	}
	openRouterRef := os.Getenv("OPENROUTER_HTTP_REFERER")
	openRouterTitle := os.Getenv("OPENROUTER_X_TITLE")
	if openRouterTitle == "" {
		openRouterTitle = appName
	}

	return Config{
		Port:                 port,
		AppName:              appName,
		MongoURI:             mongoURI,
		MongoDB:              mongoDB,
		MinIORootUser:        minioRootUser,
		MinIORootPassword:    minioRootPassword,
		MinIOEndpoint:        minioEndpoint,
		MinIOBucket:          minioBucket,
		MinIOUseSSL:          minioUseSSL,
		MinIOPublicBaseURL:   minioPublicBaseURL,
		AdminUsername:        adminUsername,
		AdminPassword:        adminPassword,
		AdminPasswordHash:    adminPasswordHash,
		AuthTokenTTLMinutes:  authTokenTTLMinutes,
		AIServiceURL:         aiServiceURL,
		AIServiceTimeout:     aiServiceTimeout,
		PublicTokenSecret:    publicTokenSecret,
		PublicTokenTTL:       publicTokenTTL,
		OpenAIAPIKey:         openAIAPIKey,
		OpenAIBaseURL:        openAIBaseURL,
		OpenAIChatModel:      openAIChatModel,
		OpenRouterRef:        openRouterRef,
		OpenRouterTitle:      openRouterTitle,
		CORSAllowOrigins:     allowOrigins,
		CORSAllowMethods:     allowMethods,
		CORSAllowHeaders:     allowHeaders,
		CORSAllowCredentials: allowCredentials,
	}
}

func ensureHeaderAllowed(existing, required string) string {
	required = strings.TrimSpace(required)
	if required == "" {
		return existing
	}
	parts := strings.Split(existing, ",")
	for _, part := range parts {
		if strings.EqualFold(strings.TrimSpace(part), required) {
			return existing
		}
	}
	if strings.TrimSpace(existing) == "" {
		return required
	}
	return existing + "," + required
}
