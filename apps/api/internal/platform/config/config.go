package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                string
	Environment            string
	Port                   string
	HTTPReadTimeout        time.Duration
	HTTPWriteTimeout       time.Duration
	HTTPShutdownTimeout    time.Duration
	DatabaseURL            string
	RedisURL               string
	LogLevel               string
	CORSAllowedOrigins     []string
	ClerkSecretKey         string
	AccessTokenSecret      string
	AccessTokenIssuer      string
	AccessTokenAudience    string
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	EncryptionKey          string
	OpenFinanceBaseURL     string
	OpenFinanceClientID    string
	OpenFinanceSecret      string
	OpenFinanceMTLSCert    string
	OpenFinanceMTLSKey     string
	CronSecret             string
	WorkerBatchSize        int
	SentryDSN              string
	SentryEnvironment      string
	SentryRelease          string
	SentryTracesSampleRate float64
	DatadogEnv             string
	DatadogService         string
	DatadogVersion         string
	DatadogTraceEnabled    bool
	OpenAIAPIKey           string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	readTimeout, err := getDurationFromSeconds("HTTP_READ_TIMEOUT", 15)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := getDurationFromSeconds("HTTP_WRITE_TIMEOUT", 15)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := getDurationFromSeconds("HTTP_SHUTDOWN_TIMEOUT", 10)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_SHUTDOWN_TIMEOUT: %w", err)
	}

	accessTTLMinutes, err := getInt("JWT_ACCESS_TTL_MINUTES", 15)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_ACCESS_TTL_MINUTES: %w", err)
	}

	refreshTTLDays, err := getInt("JWT_REFRESH_TTL_DAYS", 30)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_REFRESH_TTL_DAYS: %w", err)
	}

	workerBatchSize, err := getInt("WORKER_BATCH_SIZE", 25)
	if err != nil {
		return Config{}, fmt.Errorf("invalid WORKER_BATCH_SIZE: %w", err)
	}

	sentryTracesSampleRate, err := getFloat("SENTRY_TRACES_SAMPLE_RATE", 0)
	if err != nil {
		return Config{}, fmt.Errorf("invalid SENTRY_TRACES_SAMPLE_RATE: %w", err)
	}

	return Config{
		AppName:     getEnv("APP_NAME", "finance-api"),
		Environment: getEnv("APP_ENV", "development"),
		// Vercel Go backends rely on PORT at runtime; APP_PORT remains for local/dev overrides.
		Port:                   getEnvAny([]string{"PORT", "APP_PORT"}, "3000"),
		HTTPReadTimeout:        readTimeout,
		HTTPWriteTimeout:       writeTimeout,
		HTTPShutdownTimeout:    shutdownTimeout,
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/finance_app?sslmode=disable"),
		RedisURL:               getOptionalEnv("REDIS_URL"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigins:     getList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:8081", "http://localhost:19006"}),
		ClerkSecretKey:         getEnv("CLERK_SECRET_KEY", ""),
		AccessTokenSecret:      getEnv("JWT_ACCESS_SECRET", "change_me"),
		AccessTokenIssuer:      getEnv("JWT_ACCESS_ISSUER", "finance-api"),
		AccessTokenAudience:    getEnv("JWT_ACCESS_AUDIENCE", "finance-mobile"),
		AccessTokenTTL:         time.Duration(accessTTLMinutes) * time.Minute,
		RefreshTokenTTL:        time.Duration(refreshTTLDays) * 24 * time.Hour,
		EncryptionKey:          getEnv("APP_ENCRYPTION_KEY", "0123456789ABCDEF0123456789ABCDEF"),
		OpenFinanceBaseURL:     getEnv("OPENFINANCE_BASE_URL", ""),
		OpenFinanceClientID:    getEnv("OPENFINANCE_CLIENT_ID", ""),
		OpenFinanceSecret:      getEnv("OPENFINANCE_CLIENT_SECRET", ""),
		OpenFinanceMTLSCert:    getEnv("OPENFINANCE_MTLS_CERT_PATH", ""),
		OpenFinanceMTLSKey:     getEnv("OPENFINANCE_MTLS_KEY_PATH", ""),
		CronSecret:             getEnv("CRON_SECRET", ""),
		WorkerBatchSize:        workerBatchSize,
		SentryDSN:              getEnv("SENTRY_DSN", ""),
		SentryEnvironment:      getEnv("SENTRY_ENVIRONMENT", getEnv("APP_ENV", "development")),
		SentryRelease:          getEnv("SENTRY_RELEASE", ""),
		SentryTracesSampleRate: sentryTracesSampleRate,
		DatadogEnv:             getEnv("DD_ENV", getEnv("APP_ENV", "development")),
		DatadogService:         getEnv("DD_SERVICE", "finance-api"),
		DatadogVersion:         getEnv("DD_VERSION", ""),
		DatadogTraceEnabled:    getBool("DD_TRACE_ENABLED", false),
		OpenAIAPIKey:           getEnv("OPENAI_API_KEY", ""),
	}, nil
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getOptionalEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return fallback
}

func getInt(key string, fallback int) (int, error) {
	value := getEnv(key, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func getFloat(key string, fallback float64) (float64, error) {
	value := getEnv(key, fmt.Sprintf("%g", fallback))
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func getDurationFromSeconds(key string, fallback int) (time.Duration, error) {
	value, err := getInt(key, fallback)
	if err != nil {
		return 0, err
	}

	return time.Duration(value) * time.Second, nil
}

func getList(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	rawItems := strings.Split(value, ",")
	items := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		items = append(items, trimmed)
	}

	if len(items) == 0 {
		return fallback
	}

	return items
}

func getBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	switch strings.ToLower(value) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
