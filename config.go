package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// LDAPServer contains LDAP server configuration
type LDAPServer struct {
	URL               string        `json:"url" mapstructure:"url"`
	BindDN            string        `json:"bind_dn,omitempty" mapstructure:"bind_dn"`
	BindPassword      string        `json:"bind_password,omitempty" mapstructure:"bind_password"`
	ConnectionTimeout time.Duration `json:"connection_timeout" mapstructure:"connection_timeout"`
	SearchTimeout     time.Duration `json:"search_timeout" mapstructure:"search_timeout"`
	MaxRetries        int           `json:"max_retries" mapstructure:"max_retries"`
	RetryDelay        time.Duration `json:"retry_delay" mapstructure:"retry_delay"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	InsecureSkipVerify bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify"`
	CertFile           string `json:"cert_file,omitempty" mapstructure:"cert_file"`
	KeyFile            string `json:"key_file,omitempty" mapstructure:"key_file"`
	CAFile             string `json:"ca_file,omitempty" mapstructure:"ca_file"`
}

// Config contains the application configuration
type Config struct {
	LdapServer         LDAPServer `json:"ldap_server" mapstructure:"ldap_server"`
	BaseDN             string     `json:"base_dn" mapstructure:"base_dn"`
	PublicKeyAttribute string     `json:"public_key_attribute" mapstructure:"public_key_attribute"`
	UserAttribute      string     `json:"user_attribute" mapstructure:"user_attribute"`
	SearchFilter       string     `json:"search_filter" mapstructure:"search_filter"`
	TLS                TLSConfig  `json:"tls" mapstructure:"tls"`
	Debug              bool       `json:"debug" mapstructure:"debug"`
	ConfigPaths        []string   `json:"-"` // Not serialized, used internally
}

// setDefaults sets default configuration values
func setDefaults() {
	// LDAP Server defaults
	viper.SetDefault("ldap_server.connection_timeout", "10s")
	viper.SetDefault("ldap_server.search_timeout", "30s")
	viper.SetDefault("ldap_server.max_retries", 3)
	viper.SetDefault("ldap_server.retry_delay", "1s")

	// Search defaults
	viper.SetDefault("public_key_attribute", "sshPublicKey")
	viper.SetDefault("user_attribute", "uid")
	viper.SetDefault("search_filter", "(%s=%s)") // Will be formatted with user_attribute and username

	// TLS defaults
	viper.SetDefault("tls.insecure_skip_verify", false)

	// Application defaults
	viper.SetDefault("debug", false)
}

// setupViper configures viper with config paths and environment variable support
func setupViper(configPaths []string) {
	viper.SetConfigName("ldapPubKeyReader")
	viper.SetConfigType("json")

	// Set default config paths if none provided
	if len(configPaths) == 0 {
		configPaths = []string{
			"/etc/ssh",
			"/etc",
			".",
		}
	}

	// Add all config paths
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	// Support environment variables with LDAP_SSH_ prefix
	viper.SetEnvPrefix("LDAP_SSH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()
}

// validateConfig validates required configuration fields
func validateConfig(config *Config) error {
	if config.LdapServer.URL == "" {
		return fmt.Errorf("ldap_server.url is required")
	}

	if config.BaseDN == "" {
		return fmt.Errorf("base_dn is required")
	}

	if config.PublicKeyAttribute == "" {
		return fmt.Errorf("public_key_attribute cannot be empty")
	}

	if config.UserAttribute == "" {
		return fmt.Errorf("user_attribute cannot be empty")
	}

	if config.SearchFilter == "" {
		return fmt.Errorf("search_filter cannot be empty")
	}

	// Validate timeout values
	if config.LdapServer.ConnectionTimeout <= 0 {
		return fmt.Errorf("ldap_server.connection_timeout must be positive")
	}

	if config.LdapServer.SearchTimeout <= 0 {
		return fmt.Errorf("ldap_server.search_timeout must be positive")
	}

	if config.LdapServer.MaxRetries < 0 {
		return fmt.Errorf("ldap_server.max_retries cannot be negative")
	}

	return nil
}

// readConfig reads and validates the configuration
func readConfig(config *Config) error {
	setupViper(config.ConfigPaths)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal configuration into struct
	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	err = validateConfig(config)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

// GetConfigFile returns the path of the configuration file being used
func GetConfigFile() string {
	return viper.ConfigFileUsed()
}
