package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config - главная структура конфигурации
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	NATS     NATSConfig
	Cache    CacheConfig
	Logging  LoggingConfig
}

// ServerConfig - настройки HTTP сервера
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig - настройки PostgreSQL
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

// NATSConfig - настройки NATS Streaming
type NATSConfig struct {
	URL         string
	ClusterID   string
	ClientID    string
	Subject     string
	DurableName string
}

// CacheConfig - настройки кэша
type CacheConfig struct {
	Enabled bool
}

// LoggingConfig - настройки логирования
type LoggingConfig struct {
	Level  string
	Format string
}

// Load - загрузка конфигурации из переменных окружения
func Load() (*Config, error) {
	_ = loadEnvFile(".env")

	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvAsInt("DB_PORT", 5432),
			User:         getEnv("DB_USER", "orderuser"),
			Password:     getEnv("DB_PASSWORD", "orderpass"),
			DBName:       getEnv("DB_NAME", "orders_db"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		},
		NATS: NATSConfig{
			URL:         getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID:   getEnv("NATS_CLUSTER_ID", "test-cluster"),
			ClientID:    getEnv("NATS_CLIENT_ID", "order-service-1"),
			Subject:     getEnv("NATS_SUBJECT", "orders"),
			DurableName: getEnv("NATS_DURABLE_NAME", "order-service-durable"),
		},
		Cache: CacheConfig{
			Enabled: getEnvAsBool("CACHE_ENABLED", true),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return cfg, nil
}

// GetDSN - формирует строку подключения к PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetServerAddress - формирует адрес сервера
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// --- Вспомогательные функции ---

// getEnv - получить переменную окружения или вернуть дефолтное значение
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt - получить переменную как int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool - получить переменную как bool
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsDuration - получить переменную как time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// loadEnvFile - загружает .env файл
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем комментарии и пустые строки
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Разделяем по знаку =
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Устанавливаем переменную окружения
		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}
