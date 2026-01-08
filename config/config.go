package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

// Config holds NHN Cloud configuration
// Priority: credentials file > environment variables > interactive (runtime)
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
	NKSTenantID      string
	OBSTenantID      string

	// Track source of each credential for debugging
	sources map[string]string
	mu      sync.RWMutex
}

// CredentialSource indicates where a credential was loaded from
type CredentialSource string

const (
	SourceFile        CredentialSource = "file"
	SourceEnv         CredentialSource = "env"
	SourceInteractive CredentialSource = "interactive"
	SourceNone        CredentialSource = "none"
)

// Load creates a new Config with priority: file > env > interactive (empty initially)
func Load() *Config {
	cfg := &Config{
		sources: make(map[string]string),
	}

	// 1. Load from credentials file first
	cfg.loadFromFile()

	// 2. Override with environment variables (if set)
	cfg.loadFromEnv()

	return cfg
}

// loadFromFile loads credentials from ~/.nhncloud/credentials
func (c *Config) loadFromFile() {
	configPath := filepath.Join(os.Getenv("HOME"), ".nhncloud", "credentials")
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDefaultProfile := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			inDefaultProfile = (profile == "default")
			continue
		}

		if !inDefaultProfile {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "access_key_id":
			c.setIfEmpty("AccessKeyID", &c.AccessKeyID, value, string(SourceFile))
		case "secret_access_key":
			c.setIfEmpty("SecretAccessKey", &c.SecretAccessKey, value, string(SourceFile))
		case "region":
			c.setIfEmpty("Region", &c.Region, value, string(SourceFile))
		case "username":
			c.setIfEmpty("Username", &c.Username, value, string(SourceFile))
		case "api_password":
			c.setIfEmpty("Password", &c.Password, value, string(SourceFile))
		case "tenant_id":
			c.setIfEmpty("TenantID", &c.TenantID, value, string(SourceFile))
		case "nks_tenant_id":
			c.setIfEmpty("NKSTenantID", &c.NKSTenantID, value, string(SourceFile))
		case "obs_tenant_id":
			c.setIfEmpty("OBSTenantID", &c.OBSTenantID, value, string(SourceFile))
		case "rds_app_key":
			c.setIfEmpty("MySQLAppKey", &c.MySQLAppKey, value, string(SourceFile))
		case "rds_mariadb_app_key":
			c.setIfEmpty("MariaDBAppKey", &c.MariaDBAppKey, value, string(SourceFile))
		case "rds_postgresql_app_key":
			c.setIfEmpty("PostgreSQLAppKey", &c.PostgreSQLAppKey, value, string(SourceFile))
		}
	}
}

// loadFromEnv overrides with environment variables if set
func (c *Config) loadFromEnv() {
	c.setFromEnv("Region", &c.Region, "NHN_CLOUD_REGION")
	c.setFromEnv("AccessKeyID", &c.AccessKeyID, "NHN_CLOUD_ACCESS_KEY_ID")
	c.setFromEnv("SecretAccessKey", &c.SecretAccessKey, "NHN_CLOUD_SECRET_ACCESS_KEY")
	c.setFromEnv("MySQLAppKey", &c.MySQLAppKey, "NHN_CLOUD_MYSQL_APPKEY")
	c.setFromEnv("MariaDBAppKey", &c.MariaDBAppKey, "NHN_CLOUD_MARIADB_APPKEY")
	c.setFromEnv("PostgreSQLAppKey", &c.PostgreSQLAppKey, "NHN_CLOUD_POSTGRESQL_APPKEY")
	c.setFromEnv("Username", &c.Username, "NHN_CLOUD_USERNAME")
	c.setFromEnv("Password", &c.Password, "NHN_CLOUD_PASSWORD")
	c.setFromEnv("TenantID", &c.TenantID, "NHN_CLOUD_TENANT_ID")
	c.setFromEnv("NKSTenantID", &c.NKSTenantID, "NHN_CLOUD_NKS_TENANT_ID")
	c.setFromEnv("OBSTenantID", &c.OBSTenantID, "NHN_CLOUD_OBS_TENANT_ID")

	// Default region if not set
	if c.Region == "" {
		c.Region = "kr1"
		c.sources["Region"] = "default"
	}
}

// SetInteractive sets credentials from interactive auth (runtime)
func (c *Config) SetInteractive(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch key {
	case "access_key_id":
		c.AccessKeyID = value
		c.sources["AccessKeyID"] = string(SourceInteractive)
	case "secret_access_key":
		c.SecretAccessKey = value
		c.sources["SecretAccessKey"] = string(SourceInteractive)
	case "region":
		c.Region = value
		c.sources["Region"] = string(SourceInteractive)
	case "mysql_appkey":
		c.MySQLAppKey = value
		c.sources["MySQLAppKey"] = string(SourceInteractive)
	case "mariadb_appkey":
		c.MariaDBAppKey = value
		c.sources["MariaDBAppKey"] = string(SourceInteractive)
	case "postgresql_appkey":
		c.PostgreSQLAppKey = value
		c.sources["PostgreSQLAppKey"] = string(SourceInteractive)
	case "username":
		c.Username = value
		c.sources["Username"] = string(SourceInteractive)
	case "password":
		c.Password = value
		c.sources["Password"] = string(SourceInteractive)
	case "tenant_id":
		c.TenantID = value
		c.sources["TenantID"] = string(SourceInteractive)
	}
}

// GetSource returns where a credential was loaded from
func (c *Config) GetSource(key string) CredentialSource {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if src, ok := c.sources[key]; ok {
		return CredentialSource(src)
	}
	return SourceNone
}

// GetStatus returns a map of credential status (configured or not, source)
func (c *Config) GetStatus() map[string]map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := make(map[string]map[string]string)

	check := func(name, value string) {
		configured := "no"
		if value != "" {
			configured = "yes"
		}
		source := c.sources[name]
		if source == "" {
			source = string(SourceNone)
		}
		status[name] = map[string]string{
			"configured": configured,
			"source":     source,
		}
	}

	check("AccessKeyID", c.AccessKeyID)
	check("SecretAccessKey", c.SecretAccessKey)
	check("Region", c.Region)
	check("MySQLAppKey", c.MySQLAppKey)
	check("MariaDBAppKey", c.MariaDBAppKey)
	check("PostgreSQLAppKey", c.PostgreSQLAppKey)
	check("Username", c.Username)
	check("Password", c.Password)
	check("TenantID", c.TenantID)

	return status
}

// NewNHNCloudClient creates a new NHN Cloud SDK client
func (c *Config) NewNHNCloudClient() (*nhncloud.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

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
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AccessKeyID != "" && c.SecretAccessKey != ""
}

// HasComputeCredentials checks if Compute/Network credentials are configured
func (c *Config) HasComputeCredentials() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Username != "" && c.Password != "" && c.TenantID != ""
}

// helper functions

func (c *Config) setIfEmpty(name string, field *string, value, source string) {
	if *field == "" && value != "" {
		*field = value
		c.sources[name] = source
	}
}

func (c *Config) setFromEnv(name string, field *string, envKey string) {
	if v := os.Getenv(envKey); v != "" {
		*field = v
		c.sources[name] = string(SourceEnv)
	}
}

// LoadFromEnv is deprecated, use Load() instead
// Kept for backward compatibility
func LoadFromEnv() *Config {
	return Load()
}
