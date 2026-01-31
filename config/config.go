package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	ERPDatabase DatabaseConfig `mapstructure:"erp_database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Excel       ExcelConfig    `mapstructure:"excel"`
	Logger      LoggerConfig   `mapstructure:"logger"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	DBName   string        `mapstructure:"name"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpiryHour int    `mapstructure:"expiry_hour"`
}

type ExcelConfig struct {
	DownloadPath    string `mapstructure:"download_path"`
	MaxSearchMonths int    `mapstructure:"max_search_months"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KANBAN")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func MustConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Fatal error loading config: %s", err)
	}
	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
	)
}

func (c *Config) GetERPDatabaseDSN() string {
	// Sử dụng url.URL để xây dựng connection string an toàn
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(c.ERPDatabase.User, c.ERPDatabase.Password),
		Host:   fmt.Sprintf("%s:%d", c.ERPDatabase.Host, c.ERPDatabase.Port),
	}

	query := u.Query()
	query.Add("database", c.ERPDatabase.DBName)
	query.Add("encrypt", "disable")
	query.Add("trustServerCertificate", "true")
	// Chuyển timeout từ duration sang mili-giây hoặc giây tùy driver,
	// driver mssql thường tính bằng giây trong connection string nhưng int64
	query.Add("connection timeout", fmt.Sprintf("%d", int(c.ERPDatabase.Timeout.Seconds())))

	u.RawQuery = query.Encode()

	return u.String()
}

// GetJWTExpiry returns JWT expiry duration
func (c *Config) GetJWTExpiry() time.Duration {
	return time.Duration(c.JWT.ExpiryHour) * time.Hour
}
