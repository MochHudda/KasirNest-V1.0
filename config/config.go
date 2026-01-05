package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"kasirnest/firebase"

	"gopkg.in/ini.v1"
)

// Config holds all application configuration
type Config struct {
	Firebase *firebase.FirebaseConfig
	App      *AppConfig
	Security *SecurityConfig
	Database *DatabaseConfig
	filePath string
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name         string
	Version      string
	Debug        bool
	WindowWidth  int
	WindowHeight int
	Theme        string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EncryptionKey  string
	SessionTimeout int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	AutoBackup     bool
	BackupInterval int
}

// Load loads configuration from app.ini file
func Load() (*Config, error) {
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	// Try different possible locations for config file
	configPaths := []string{
		filepath.Join(exeDir, "config", "app.ini"),
		filepath.Join("config", "app.ini"),
		"app.ini",
	}

	var configFile string
	var found bool

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("configuration file not found. Please create app.ini from app.ini.example")
	}

	cfg, err := ini.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	config := &Config{
		filePath: configFile,
	}

	// Load Firebase configuration
	config.Firebase = &firebase.FirebaseConfig{
		ProjectID:        cfg.Section("firebase").Key("project_id").String(),
		PrivateKeyID:     cfg.Section("firebase").Key("private_key_id").String(),
		PrivateKey:       strings.ReplaceAll(cfg.Section("firebase").Key("private_key").String(), "\\n", "\n"),
		ClientEmail:      cfg.Section("firebase").Key("client_email").String(),
		ClientID:         cfg.Section("firebase").Key("client_id").String(),
		AuthURI:          cfg.Section("firebase").Key("auth_uri").String(),
		TokenURI:         cfg.Section("firebase").Key("token_uri").String(),
		AuthProviderX509: cfg.Section("firebase").Key("auth_provider_x509_cert_url").String(),
		ClientX509:       cfg.Section("firebase").Key("client_x509_cert_url").String(),
		StorageBucket:    cfg.Section("firebase").Key("storage_bucket").String(),
	}

	// Load App configuration
	config.App = &AppConfig{
		Name:         cfg.Section("app").Key("name").MustString("KasirNest"),
		Version:      cfg.Section("app").Key("version").MustString("1.0.0"),
		Debug:        cfg.Section("app").Key("debug").MustBool(false),
		WindowWidth:  cfg.Section("app").Key("window_width").MustInt(1200),
		WindowHeight: cfg.Section("app").Key("window_height").MustInt(800),
		Theme:        cfg.Section("app").Key("theme").MustString("light"),
	}

	// Load Security configuration
	config.Security = &SecurityConfig{
		EncryptionKey:  cfg.Section("security").Key("encryption_key").String(),
		SessionTimeout: cfg.Section("security").Key("session_timeout").MustInt(3600),
	}

	// Load Database configuration
	config.Database = &DatabaseConfig{
		AutoBackup:     cfg.Section("database").Key("auto_backup").MustBool(true),
		BackupInterval: cfg.Section("database").Key("backup_interval").MustInt(24),
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate Firebase config
	if c.Firebase.ProjectID == "" || strings.Contains(c.Firebase.ProjectID, "your-") {
		return fmt.Errorf("firebase project_id is not configured")
	}

	if c.Firebase.ClientEmail == "" || strings.Contains(c.Firebase.ClientEmail, "your-") {
		return fmt.Errorf("firebase client_email is not configured")
	}

	if c.Firebase.PrivateKey == "" || strings.Contains(c.Firebase.PrivateKey, "Your-Private-Key") {
		return fmt.Errorf("firebase private_key is not configured")
	}

	if c.Security.EncryptionKey == "" || strings.Contains(c.Security.EncryptionKey, "your-") {
		return fmt.Errorf("security encryption_key is not configured")
	}

	return nil
}

// Save saves the current configuration to file
func (c *Config) Save() error {
	cfg := ini.Empty()

	// Firebase section
	firebaseSection, _ := cfg.NewSection("firebase")
	firebaseSection.NewKey("project_id", c.Firebase.ProjectID)
	firebaseSection.NewKey("private_key_id", c.Firebase.PrivateKeyID)
	firebaseSection.NewKey("private_key", strings.ReplaceAll(c.Firebase.PrivateKey, "\n", "\\n"))
	firebaseSection.NewKey("client_email", c.Firebase.ClientEmail)
	firebaseSection.NewKey("client_id", c.Firebase.ClientID)
	firebaseSection.NewKey("auth_uri", c.Firebase.AuthURI)
	firebaseSection.NewKey("token_uri", c.Firebase.TokenURI)
	firebaseSection.NewKey("auth_provider_x509_cert_url", c.Firebase.AuthProviderX509)
	firebaseSection.NewKey("client_x509_cert_url", c.Firebase.ClientX509)
	firebaseSection.NewKey("storage_bucket", c.Firebase.StorageBucket)

	// App section
	appSection, _ := cfg.NewSection("app")
	appSection.NewKey("name", c.App.Name)
	appSection.NewKey("version", c.App.Version)
	appSection.NewKey("debug", strconv.FormatBool(c.App.Debug))
	appSection.NewKey("window_width", strconv.Itoa(c.App.WindowWidth))
	appSection.NewKey("window_height", strconv.Itoa(c.App.WindowHeight))
	appSection.NewKey("theme", c.App.Theme)

	// Security section
	securitySection, _ := cfg.NewSection("security")
	securitySection.NewKey("encryption_key", c.Security.EncryptionKey)
	securitySection.NewKey("session_timeout", strconv.Itoa(c.Security.SessionTimeout))

	// Database section
	databaseSection, _ := cfg.NewSection("database")
	databaseSection.NewKey("auto_backup", strconv.FormatBool(c.Database.AutoBackup))
	databaseSection.NewKey("backup_interval", strconv.Itoa(c.Database.BackupInterval))

	return cfg.SaveTo(c.filePath)
}

// GetConfigPath returns the path to the configuration file
func (c *Config) GetConfigPath() string {
	return c.filePath
}

// IsDebug returns true if debug mode is enabled
func (c *Config) IsDebug() bool {
	return c.App.Debug
}

// GetWindowSize returns window width and height
func (c *Config) GetWindowSize() (int, int) {
	return c.App.WindowWidth, c.App.WindowHeight
}
