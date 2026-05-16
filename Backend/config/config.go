package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App           AppConfig
	Postgres      PostgresConfig
	Redis         RedisConfig
	JWT           JWTConfig
	RateLimit     RateLimitConfig
	AI            AIConfig
	Observability ObservabilityConfig
	CORS          CORSConfig
	Security      SecurityConfig
}

type AppConfig struct {
	Env  string
	Name string
	Port string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type RateLimitConfig struct {
	Authenticated int
	Anonymous     int
}

type AIConfig struct {
	Provider        string
	OpenAIAPIKey    string
	AnthropicAPIKey string
}

type ObservabilityConfig struct {
	SentryDSN string
	LogLevel  string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SecurityConfig struct {
	PasswordEncryptionPrivateKeyB64 string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		return nil, err
	}
	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := time.ParseDuration(getEnv("POSTGRES_CONN_MAX_LIFETIME", "1h"))
	if err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Name: getEnv("APP_NAME", "flowforge-api"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnv("POSTGRES_PORT", "5432"),
			User:            getEnv("POSTGRES_USER", "flowforge"),
			Password:        getEnv("POSTGRES_PASSWORD", "secret"),
			Name:            getEnv("POSTGRES_NAME", "flowforge_db"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("POSTGRES_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: connMaxLifetime,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "change-me-min-32-characters"),
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		},
		RateLimit: RateLimitConfig{
			Authenticated: getEnvInt("RATE_LIMIT_AUTHENTICATED", 100),
			Anonymous:     getEnvInt("RATE_LIMIT_ANONYMOUS", 20),
		},
		AI: AIConfig{
			Provider:        getEnv("AI_PROVIDER", ""),
			OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
			AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		},
		Observability: ObservabilityConfig{
			SentryDSN: getEnv("SENTRY_DSN", ""),
			LogLevel:  getEnv("LOG_LEVEL", "info"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		},
		Security: SecurityConfig{
			PasswordEncryptionPrivateKeyB64: getEnv("PASSWORD_ENCRYPTION_PRIVATE_KEY_B64", ""),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
