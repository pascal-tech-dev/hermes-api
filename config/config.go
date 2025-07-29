package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Environment  string        `mapstructure:"environment"`
	LogLevel     string        `mapstructure:"log_level"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret   string   `mapstructure:"jwt_secret"`
	BcryptCost  int      `mapstructure:"bcrypt_cost"`
	CORSOrigins []string `mapstructure:"cors_origins"`
}

// Load loads configuration from multiple sources
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read from config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Try to read config file, but don't fail if it doesn't exist
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Parse into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.log_level", "info")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "120s")

	// Security defaults
	v.SetDefault("security.jwt_secret", "your-super-secret-jwt-key-change-in-production")
	v.SetDefault("security.bcrypt_cost", 12)

	// Database defaults
	v.SetDefault("database.host", "hermes-db")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.name", "hermes")
	v.SetDefault("database.user", "hermes")
	v.SetDefault("database.ssl_mode", "disable")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.db", 0)

	// Logging defaults
	v.SetDefault("logging.level", "debug")
	v.SetDefault("logging.format", "text")
	v.SetDefault("logging.output", "stdout")
}
