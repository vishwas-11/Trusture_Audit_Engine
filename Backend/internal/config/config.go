package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	Port            string
	MongoURI        string
	PolygonRPC      string
	ContractAddress string
	JWTSecret       string
	GeminiAPIKey    string
	OpenAIAPIKey    string
	AuditPassThreshold float64
	DistanceWarnKM     float64
	DistanceFailKM     float64
	GSTVerificationEnabled bool
	GSTAPITimeoutMS    int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, relying on environment variables")
	}

	cfg := &Config{
		AppName:         os.Getenv("APP_NAME"),
		Port:            os.Getenv("PORT"),
		MongoURI:        os.Getenv("MONGO_URI"),
		PolygonRPC:      os.Getenv("POLYGON_RPC"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		GeminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		AuditPassThreshold:   getFloatEnv("AUDIT_PASS_THRESHOLD", 0.70),
		DistanceWarnKM:       getFloatEnv("AUDIT_DISTANCE_WARN_KM", 100),
		DistanceFailKM:       getFloatEnv("AUDIT_DISTANCE_FAIL_KM", 600),
		GSTVerificationEnabled: getBoolEnv("GST_VERIFICATION_ENABLED", true),
		GSTAPITimeoutMS:      getIntEnv("GST_API_TIMEOUT_MS", 1500),
	}

	// 🔒 Fail fast (important for blockchain apps)
	if cfg.MongoURI == "" || cfg.PolygonRPC == "" || cfg.ContractAddress == "" {
		log.Fatal("❌ Missing required environment variables")
	}

	return cfg
}

func getFloatEnv(key string, fallback float64) float64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return val
}

func getBoolEnv(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return val
}

func getIntEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return val
}
