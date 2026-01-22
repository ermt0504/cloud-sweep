package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	AWS      AWSConfig
	Azure    AzureConfig
	GCP      GCPConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string
	Environment string
	Debug       bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// AWSConfig holds AWS configuration
type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// AzureConfig holds Azure configuration
type AzureConfig struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
}

// GCPConfig holds GCP configuration
type GCPConfig struct {
	ProjectID       string
	CredentialsFile string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.debug", true)

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "cloudsweep")
	v.SetDefault("database.password", "cloudsweep_secret")
	v.SetDefault("database.name", "cloudsweep")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("aws.region", "us-east-1")

	// Config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/cloudsweep")

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind specific env vars
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.environment", "SERVER_ENV")
	v.BindEnv("server.debug", "SERVER_DEBUG")

	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.name", "DB_NAME")
	v.BindEnv("database.sslmode", "DB_SSLMODE")

	v.BindEnv("redis.addr", "REDIS_ADDR")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("redis.db", "REDIS_DB")

	v.BindEnv("aws.region", "AWS_REGION")
	v.BindEnv("aws.accesskeyid", "AWS_ACCESS_KEY_ID")
	v.BindEnv("aws.secretaccesskey", "AWS_SECRET_ACCESS_KEY")

	config := &Config{
		Server: ServerConfig{
			Port:        v.GetString("server.port"),
			Environment: v.GetString("server.environment"),
			Debug:       v.GetBool("server.debug"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetString("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			Name:     v.GetString("database.name"),
			SSLMode:  v.GetString("database.sslmode"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("redis.addr"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		AWS: AWSConfig{
			Region:          v.GetString("aws.region"),
			AccessKeyID:     v.GetString("aws.accesskeyid"),
			SecretAccessKey: v.GetString("aws.secretaccesskey"),
		},
		Azure: AzureConfig{
			TenantID:       v.GetString("azure.tenantid"),
			ClientID:       v.GetString("azure.clientid"),
			ClientSecret:   v.GetString("azure.clientsecret"),
			SubscriptionID: v.GetString("azure.subscriptionid"),
		},
		GCP: GCPConfig{
			ProjectID:       v.GetString("gcp.projectid"),
			CredentialsFile: v.GetString("gcp.credentialsfile"),
		},
	}

	return config, nil
}
