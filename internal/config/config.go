package config

import (
	"os"
	"time"
)

type Config struct {
	AppEnv        string
	AppPort       string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	MongoURI string
	MongoDB  string

	TwitchClientID     string
	TwitchClientSecret string
	TwitchRedirectURI  string

	StripeSecretKey    string
	StripeWebhookSecret string

	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Bucket          string
	S3Region          string
}

func Load() *Config {
	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		JWTSecret:     getEnv("APP_JWT_SECRET", "super-secret-change-me-in-production"),
		JWTAccessTTL:  parseDurationEnv("APP_JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL: parseDurationEnv("APP_JWT_REFRESH_TTL", 7*24*time.Hour),

		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGODB_DB", "raffle_saas"),

		TwitchClientID:     getEnv("TWITCH_CLIENT_ID", ""),
		TwitchClientSecret: getEnv("TWITCH_CLIENT_SECRET", ""),
		TwitchRedirectURI:  getEnv("TWITCH_REDIRECT_URI", ""),

		StripeSecretKey:    getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),

		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
		S3Bucket:          getEnv("S3_BUCKET", ""),
		S3Region:          getEnv("S3_REGION", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
