package config

import (
	"os"
	"strconv"
	"time"
)

// Config aggregates application settings loaded from environment variables.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Storage  StorageConfig
	OAuth    OAuthConfig
	Session  SessionConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type RedisConfig struct {
	Addr     string
	DB       int
	Password string
}

type StorageConfig struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
	CDNBaseURL      string
	SignedURLTTL    time.Duration
	DefaultImageURL string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	VKClientID         string
	VKClientSecret     string
	VKRedirectURL      string
	TelegramBotToken   string
}

type SessionConfig struct {
	Key string
}

var loaded Config

func LoadConfig() Config {
	addr := getenv("SERVER_ADDR", "")
	if addr == "" {
		if port := os.Getenv("PORT"); port != "" {
			addr = ":" + port
		} else {
			addr = ":8080"
		}
	}

	cfg := Config{
		Server: ServerConfig{
			Addr: addr,
		},
		Database: DatabaseConfig{
			Host:            getenv("DB_HOST", "localhost"),
			Port:            getenv("DB_PORT", "5432"),
			User:            getenv("DB_USER", "postgres"),
			Password:        getenv("DB_PASSWORD", ""),
			Name:            getenv("DB_NAME", "app"),
			SSLMode:         getenv("DB_SSL_MODE", "disable"),
			MaxConns:        int32(getenvInt("DB_MAX_CONNS", 10)),
			MinConns:        int32(getenvInt("DB_MIN_CONNS", 2)),
			MaxConnLifetime: getenvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			MaxConnIdleTime: getenvDuration("DB_MAX_CONN_IDLE", 15*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getenv("REDIS_ADDR", "localhost:6379"),
			DB:       getenvInt("REDIS_DB", 0),
			Password: getenv("REDIS_PASSWORD", ""),
		},
		Storage: StorageConfig{
			Endpoint:        getenv("S3_ENDPOINT", "storage.yandexcloud.net"),
			Region:          getenv("S3_REGION", "ru-central1"),
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			Bucket:          os.Getenv("S3_BUCKET"),
			UseSSL:          getenvBool("S3_USE_SSL", true),
			CDNBaseURL:      os.Getenv("S3_CDN_BASE_URL"),
			SignedURLTTL:    getenvDuration("S3_SIGNED_URL_TTL", 15*time.Minute),
			DefaultImageURL: getenv("DEFAULT_IMAGE_URL", ""),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  getenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			VKClientID:         os.Getenv("VK_CLIENT_ID"),
			VKClientSecret:     os.Getenv("VK_CLIENT_SECRET"),
			VKRedirectURL:      getenv("VK_REDIRECT_URL", "http://localhost:8080/auth/vk/callback"),
			TelegramBotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
		Session: SessionConfig{
			Key: getenv("SESSION_KEY", ""),
		},
	}

	loaded = cfg
	return cfg
}

func Get() Config {
	return loaded
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getenvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getenvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getenvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
