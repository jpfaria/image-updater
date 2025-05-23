package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Docker   DockerConfig
	Git      GitConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// DockerConfig holds the Docker registry configuration
type DockerConfig struct {
	RegistryURL      string
	PollingInterval  int // in seconds
	CredentialsPath  string
	DefaultNamespace string
}

// GitConfig holds the Git repository configuration
type GitConfig struct {
	DefaultBranch string
	CommitMessage string
	AuthType      string // ssh or https
	Username      string
	Password      string
	SSHKeyPath    string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Host: getEnvStr("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Type:     getEnvStr("DB_TYPE", "sqlite"),
			Host:     getEnvStr("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnvStr("DB_USER", "postgres"),
			Password: getEnvStr("DB_PASSWORD", "postgres"),
			Name:     getEnvStr("DB_NAME", "image_updater"),
		},
		Docker: DockerConfig{
			RegistryURL:      getEnvStr("DOCKER_REGISTRY_URL", "docker.io"),
			PollingInterval:  getEnvInt("DOCKER_POLLING_INTERVAL", 300),
			CredentialsPath:  getEnvStr("DOCKER_CREDENTIALS_PATH", ""),
			DefaultNamespace: getEnvStr("DOCKER_DEFAULT_NAMESPACE", "library"),
		},
		Git: GitConfig{
			DefaultBranch: getEnvStr("GIT_DEFAULT_BRANCH", "main"),
			CommitMessage: getEnvStr("GIT_COMMIT_MESSAGE", "Update image version to %s"),
			AuthType:      getEnvStr("GIT_AUTH_TYPE", "https"),
			Username:      getEnvStr("GIT_USERNAME", ""),
			Password:      getEnvStr("GIT_PASSWORD", ""),
			SSHKeyPath:    getEnvStr("GIT_SSH_KEY_PATH", ""),
		},
	}

	return cfg, nil
}

// Helper functions to get environment variables with defaults
func getEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
