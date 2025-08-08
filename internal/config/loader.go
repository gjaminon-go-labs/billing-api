// Configuration Loading System
//
// This file implements configuration loading with YAML files and environment variable overrides.
// Provides: Kubernetes-ready configuration management, environment-specific settings
// Pattern: Base configuration + environment overrides + environment variable substitution
// Used by: Production main.go, development environments, Kubernetes deployments
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	Storage           StorageConfig   `yaml:"storage"`
	Migration         MigrationConfig `yaml:"migration"`
	Server            ServerConfig    `yaml:"server"`
	Database          DatabaseConfig  `yaml:"database"`
	MigrationDatabase DatabaseConfig  `yaml:"migration_database"`
	Logging           LoggingConfig   `yaml:"logging"`
	API               APIConfig       `yaml:"api"`
	RateLimit         RateLimitConfig `yaml:"rate_limit"`
	Health            HealthConfig    `yaml:"health"`
	Metrics           MetricsConfig   `yaml:"metrics"`
	Tracing           TracingConfig   `yaml:"tracing"`
}

// StorageConfig defines storage configuration
type StorageConfig struct {
	Type string `yaml:"type"` // memory, postgres
}

// MigrationConfig defines database migration configuration
type MigrationConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Path        string `yaml:"path"`
	AutoMigrate bool   `yaml:"auto_migrate"`
	TableName   string `yaml:"table_name"`
}

// ServerConfig defines HTTP server configuration
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DatabaseConfig defines database connection configuration
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	Schema          string        `yaml:"schema"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	LogLevel        string        `yaml:"log_level"`
}

// LoggingConfig defines logging configuration
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

// APIConfig defines API-specific configuration
type APIConfig struct {
	Prefix      string   `yaml:"prefix"`
	EnableCORS  bool     `yaml:"enable_cors"`
	CORSOrigins []string `yaml:"cors_origins"`
	CORSMethods []string `yaml:"cors_methods"`
	CORSHeaders []string `yaml:"cors_headers"`
}

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
	Burst             int  `yaml:"burst"`
}

// HealthConfig defines health check configuration
type HealthConfig struct {
	Endpoint      string `yaml:"endpoint"`
	DatabaseCheck bool   `yaml:"database_check"`
}

// MetricsConfig defines metrics configuration
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	Namespace string `yaml:"namespace"`
}

// TracingConfig defines tracing configuration
type TracingConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ServiceName    string `yaml:"service_name"`
	JaegerEndpoint string `yaml:"jaeger_endpoint"`
}

// LoadConfig loads configuration from YAML files with environment overrides
func LoadConfig(environment string) (*Config, error) {
	// Load base configuration
	config, err := loadBaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	// Load environment-specific overrides
	if environment != "" {
		err = loadEnvironmentConfig(config, environment)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s config: %w", environment, err)
		}
	}

	// Apply environment variable overrides (Kubernetes secrets/configmaps)
	applyEnvironmentVariables(config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// loadBaseConfig loads the base configuration file
func loadBaseConfig() (*Config, error) {
	configPath := getConfigPath("base.yaml")
	return loadConfigFile(configPath)
}

// loadEnvironmentConfig loads environment-specific configuration overrides
func loadEnvironmentConfig(config *Config, environment string) error {
	configPath := getConfigPath(environment + ".yaml")

	// Check if environment config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Environment config is optional
		return nil
	}

	envConfig, err := loadConfigFile(configPath)
	if err != nil {
		return err
	}

	// Merge environment config into base config
	mergeConfigs(config, envConfig)
	return nil
}

// loadConfigFile loads a YAML configuration file
func loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return &config, nil
}

// getConfigPath returns the full path to a configuration file
func getConfigPath(filename string) string {
	// Check for custom config directory from environment
	if configDir := os.Getenv("CONFIG_DIR"); configDir != "" {
		return filepath.Join(configDir, filename)
	}

	// Default to configs directory relative to project root
	return filepath.Join("configs", filename)
}

// applyEnvironmentVariables overrides configuration with environment variables
func applyEnvironmentVariables(config *Config) {
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}

	// Database configuration (Kubernetes secrets)
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = p
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	// Migration database configuration (Kubernetes secrets)
	if migrationDbHost := os.Getenv("MIGRATION_DB_HOST"); migrationDbHost != "" {
		config.MigrationDatabase.Host = migrationDbHost
	}
	if migrationDbPort := os.Getenv("MIGRATION_DB_PORT"); migrationDbPort != "" {
		if p, err := strconv.Atoi(migrationDbPort); err == nil {
			config.MigrationDatabase.Port = p
		}
	}
	if migrationDbUser := os.Getenv("MIGRATION_DB_USER"); migrationDbUser != "" {
		config.MigrationDatabase.User = migrationDbUser
	}
	if migrationDbPassword := os.Getenv("MIGRATION_DB_PASSWORD"); migrationDbPassword != "" {
		config.MigrationDatabase.Password = migrationDbPassword
	}
	if migrationDbName := os.Getenv("MIGRATION_DB_NAME"); migrationDbName != "" {
		config.MigrationDatabase.DBName = migrationDbName
	}

	// Storage configuration
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		config.Storage.Type = storageType
	}

	// Migration configuration
	if autoMigrate := os.Getenv("AUTO_MIGRATE"); autoMigrate != "" {
		config.Migration.AutoMigrate = autoMigrate == "true"
	}

	// Logging configuration
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
}

// mergeConfigs merges source configuration into target configuration
func mergeConfigs(target, source *Config) {
	// Note: This is a simplified merge - in production you might want
	// a more sophisticated merging strategy using reflection or a library

	// Storage config
	if source.Storage.Type != "" {
		target.Storage.Type = source.Storage.Type
	}

	// Migration config
	if source.Migration.Path != "" {
		target.Migration.Path = source.Migration.Path
	}
	if source.Migration.TableName != "" {
		target.Migration.TableName = source.Migration.TableName
	}
	// Note: bool fields are merged only if explicitly set in source
	target.Migration.Enabled = source.Migration.Enabled || target.Migration.Enabled
	target.Migration.AutoMigrate = source.Migration.AutoMigrate || target.Migration.AutoMigrate

	// Server config
	if source.Server.Port != 0 {
		target.Server.Port = source.Server.Port
	}
	if source.Server.Host != "" {
		target.Server.Host = source.Server.Host
	}

	// Database config
	if source.Database.Host != "" {
		target.Database.Host = source.Database.Host
	}
	if source.Database.Port != 0 {
		target.Database.Port = source.Database.Port
	}
	if source.Database.DBName != "" {
		target.Database.DBName = source.Database.DBName
	}
	if source.Database.User != "" {
		target.Database.User = source.Database.User
	}
	if source.Database.Password != "" {
		target.Database.Password = source.Database.Password
	}
	if source.Database.Schema != "" {
		target.Database.Schema = source.Database.Schema
	}

	// Migration database config
	if source.MigrationDatabase.Host != "" {
		target.MigrationDatabase.Host = source.MigrationDatabase.Host
	}
	if source.MigrationDatabase.Port != 0 {
		target.MigrationDatabase.Port = source.MigrationDatabase.Port
	}
	if source.MigrationDatabase.DBName != "" {
		target.MigrationDatabase.DBName = source.MigrationDatabase.DBName
	}
	if source.MigrationDatabase.User != "" {
		target.MigrationDatabase.User = source.MigrationDatabase.User
	}
	if source.MigrationDatabase.Password != "" {
		target.MigrationDatabase.Password = source.MigrationDatabase.Password
	}
	if source.MigrationDatabase.Schema != "" {
		target.MigrationDatabase.Schema = source.MigrationDatabase.Schema
	}
	if source.MigrationDatabase.SSLMode != "" {
		target.MigrationDatabase.SSLMode = source.MigrationDatabase.SSLMode
	}

	// Logging config
	if source.Logging.Level != "" {
		target.Logging.Level = source.Logging.Level
	}
	if source.Logging.Format != "" {
		target.Logging.Format = source.Logging.Format
	}
}

// validateConfig validates the loaded configuration
func validateConfig(config *Config) error {
	// Storage validation
	validStorageTypes := []string{"memory", "postgres"}
	if !contains(validStorageTypes, config.Storage.Type) {
		return fmt.Errorf("invalid storage type: %s (must be one of: %s)", config.Storage.Type, strings.Join(validStorageTypes, ", "))
	}

	// Server validation
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Database validation
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Logging validation
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, strings.ToLower(config.Logging.Level)) {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetEnvironment returns the current environment from ENV variable or default
func GetEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		return "development" // Default environment
	}
	return env
}
