package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port int
	Env  string // development | production
}

type DatabaseConfig struct {
	Host         string
	Port         int
	Name         string
	User         string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	PrivateKeyPath      string
	PublicKeyPath       string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	User     string
	Password string
	From     string
}

type SecurityConfig struct {
	EncryptionKey  string
	AllowedOrigins []string
	FrontendURL    string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Bind all expected env vars with defaults
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRY", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRY", "168h")
	viper.SetDefault("SMTP_PORT", 587)

	_ = viper.ReadInConfig() // not required - env vars take precedence

	cfg := &Config{}

	// Server
	cfg.Server.Port = viper.GetInt("SERVER_PORT")
	cfg.Server.Env = viper.GetString("SERVER_ENV")

	// Database - required
	cfg.Database.Host = requireEnv("DB_HOST")
	cfg.Database.Port = viper.GetInt("DB_PORT")
	cfg.Database.Name = requireEnv("DB_NAME")
	cfg.Database.User = requireEnv("DB_USER")
	cfg.Database.Password = requireEnv("DB_PASSWORD")
	cfg.Database.MaxOpenConns = viper.GetInt("DB_MAX_OPEN_CONNS")
	cfg.Database.MaxIdleConns = viper.GetInt("DB_MAX_IDLE_CONNS")

	// Redis
	cfg.Redis.URL = requireEnv("REDIS_URL")

	// JWT - required
	cfg.JWT.PrivateKeyPath = requireEnv("JWT_PRIVATE_KEY_PATH")
	cfg.JWT.PublicKeyPath = requireEnv("JWT_PUBLIC_KEY_PATH")

	accessExpiry, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_EXPIRY: %w", err)
	}
	cfg.JWT.AccessTokenExpiry = accessExpiry

	refreshExpiry, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_EXPIRY: %w", err)
	}
	cfg.JWT.RefreshTokenExpiry = refreshExpiry

	// Email
	cfg.Email.SMTPHost = requireEnv("SMTP_HOST")
	cfg.Email.SMTPPort = viper.GetInt("SMTP_PORT")
	cfg.Email.User = requireEnv("SMTP_USER")
	cfg.Email.Password = requireEnv("SMTP_PASSWORD")
	cfg.Email.From = requireEnv("EMAIL_FROM")

	// Security
	cfg.Security.EncryptionKey = requireEnv("ENCRYPTION_KEY")
	cfg.Security.AllowedOrigins = viper.GetStringSlice("ALLOWED_ORIGINS")
	cfg.Security.FrontendURL = requireEnv("FRONTEND_URL")

	return cfg, nil
}

func requireEnv(key string) string {
	val := viper.GetString(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return val
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}
