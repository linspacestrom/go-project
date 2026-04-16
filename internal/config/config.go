package config

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   HTTPConfig     `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	Postgres PostgresConfig `yaml:"postgres"`
	Auth     AuthConfig     `yaml:"auth"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Outbox   OutboxConfig   `yaml:"outbox"`
}

type AuthConfig struct {
	Secret          string        `env:"JWT_SECRET" env-default:"my-super-secret-key"`
	TokenTTL        time.Duration `env:"JWT_TTL" env-default:"24h"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"720h"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic   string   `env:"KAFKA_TOPIC" env-default:"notification-events"`
}

type OutboxConfig struct {
	DispatchInterval time.Duration `env:"OUTBOX_DISPATCH_INTERVAL" env-default:"2s"`
	BatchSize        uint64        `env:"OUTBOX_BATCH_SIZE" env-default:"100"`
}

type HTTPConfig struct {
	Host    string `env:"HTTP_HOST" env-default:"0.0.0.0" yaml:"host"`
	Port    string `env:"HTTP_PORT" env-default:"8080"    yaml:"port"`
	Mode    string `env:"GIN_MODE"  env-default:"release"`
	Timeout struct {
		Server time.Duration `yaml:"server"`
		Write  time.Duration `yaml:"write"`
		Read   time.Duration `yaml:"read"`
		Idle   time.Duration `yaml:"idle"`
	} `yaml:"timeout"`
}

type LoggerConfig struct {
	Path string `env:"LOGGER_CONFIG_PATH" env-default:"config/logger.json"`
}

type PostgresConfig struct {
	Host           string `env:"POSTGRES_HOST"       env-default:"localhost"`
	Port           string `env:"POSTGRES_PORT"       env-default:"5432"`
	Username       string `env:"POSTGRES_USER"       env-default:"timur"`
	Password       string `env:"POSTGRES_PASSWORD"   env-default:"Mars237s!"`
	Database       string `env:"POSTGRES_DB"         env-default:"tbank"`
	SSLMode        string `env:"POSTGRES_SSL_MODE"   env-default:"disable"`
	MigrationsPath string `env:"GOOSE_MIGRATION_DIR" env-default:"./migrations"`

	MaxOpenConns    int32         `env-default:"25"  yaml:"max_open_conns"`
	MaxIdleConns    int32         `env-default:"5"   yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `env-default:"1h"  yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `env-default:"30m" yaml:"conn_max_idle_time"`
	Timeout         time.Duration `env-default:"5s"  yaml:"timeout"`
}

const (
	defaultConfigPath = "config/config.yaml"
	Path              = "CONFIG_PATH"
)

func MustLoad() *Config {
	cfg, err := load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	return cfg
}

func load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if loadErr := godotenv.Load(); loadErr != nil {
			return nil, fmt.Errorf("failed to load .env: %w", loadErr)
		}
	}

	configPath := getConfigPath()
	if configPath != "" {
		fileInfo, err := os.Stat(configPath)

		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file %s does not exist", configPath)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to access config file: %w", err)
		}

		if fileInfo.IsDir() {
			return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
		}
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read env vars: %w", err)
	}

	return &cfg, nil
}

func getConfigPath() string {
	configPath := os.Getenv(Path)
	if configPath == "" {
		configPath = defaultConfigPath
	}

	return configPath
}

func (h *HTTPConfig) GetAddr() string {
	return net.JoinHostPort(h.Host, h.Port)
}
