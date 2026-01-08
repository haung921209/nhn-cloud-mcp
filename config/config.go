package config

import (
	"os"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

// Config holds NHN Cloud configuration
type Config struct {
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	MySQLAppKey      string
	MariaDBAppKey    string
	PostgreSQLAppKey string
	Username         string
	Password         string
	TenantID         string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	return &Config{
		Region:           getEnvOrDefault("NHN_CLOUD_REGION", "kr1"),
		AccessKeyID:      os.Getenv("NHN_CLOUD_ACCESS_KEY_ID"),
		SecretAccessKey:  os.Getenv("NHN_CLOUD_SECRET_ACCESS_KEY"),
		MySQLAppKey:      os.Getenv("NHN_CLOUD_MYSQL_APPKEY"),
		MariaDBAppKey:    os.Getenv("NHN_CLOUD_MARIADB_APPKEY"),
		PostgreSQLAppKey: os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY"),
		Username:         os.Getenv("NHN_CLOUD_USERNAME"),
		Password:         os.Getenv("NHN_CLOUD_PASSWORD"),
		TenantID:         os.Getenv("NHN_CLOUD_TENANT_ID"),
	}
}

// NewNHNCloudClient creates a new NHN Cloud SDK client
func (c *Config) NewNHNCloudClient() (*nhncloud.Client, error) {
	creds := credentials.NewStatic(c.AccessKeyID, c.SecretAccessKey)

	var identityCreds credentials.IdentityCredentials
	if c.Username != "" && c.Password != "" && c.TenantID != "" {
		identityCreds = credentials.NewStaticIdentity(c.Username, c.Password, c.TenantID)
	}

	cfg := &nhncloud.Config{
		Region:              c.Region,
		Credentials:         creds,
		IdentityCredentials: identityCreds,
		AppKeys: map[string]string{
			"rds-mysql":      c.MySQLAppKey,
			"rds-mariadb":    c.MariaDBAppKey,
			"rds-postgresql": c.PostgreSQLAppKey,
		},
	}

	return nhncloud.New(cfg)
}

// HasRDSCredentials checks if RDS credentials are configured
func (c *Config) HasRDSCredentials() bool {
	return c.AccessKeyID != "" && c.SecretAccessKey != ""
}

// HasComputeCredentials checks if Compute/Network credentials are configured
func (c *Config) HasComputeCredentials() bool {
	return c.Username != "" && c.Password != "" && c.TenantID != ""
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
