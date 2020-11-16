package database

import (
	"database/sql"
	"fmt"
	"time"

	env "github.com/NuclearLouse/utilities-environ"
)

// Config ...
type Config struct {
	User           string
	Pass           string
	Host           string
	Port           string
	Database       string
	SSLMode        string
	MaxOpenConns   int
	MaxIdleConns   int
	ConMaxLifeTime time.Duration
	ConMaxIdleTime time.Duration
}

// DefaultConfigPostgres ...
func DefaultConfigPostgres() *Config {
	return &Config{
		User:           "postgres",
		Pass:           "postgres",
		Host:           "localhost",
		Port:           "5432",
		Database:       "postgres",
		SSLMode:        "disable",
		MaxIdleConns:   15,
		MaxOpenConns:   15,
		ConMaxIdleTime: 15,
		ConMaxLifeTime: 5,
	}
}

// Connect функция возвращающая соединение с базой данных. Если не передан укзатель на Config, настройки
// будут читаться из переменных окружения.
func Connect(driver string, config ...*Config) (*sql.DB, error) {
	var (
		cfg     *Config
		connStr string
	)

	if config == nil {
		cfg = &Config{
			User:           env.GetEnv("DB_USER", "postgres"),
			Pass:           env.GetEnv("DB_PASSWORD", "postgres"),
			Host:           env.GetEnv("DB_HOST", "localhost"),
			Port:           env.GetEnv("DB_PORT", "5432"),
			Database:       env.GetEnv("DB_DATABASE", "postgres"),
			SSLMode:        env.GetEnv("DB_SSL_MODE", "disable"),
			MaxIdleConns:   env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 15),
			MaxOpenConns:   env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 15),
			ConMaxIdleTime: time.Duration(env.GetEnvAsInt64("DB_CONN_MAX_IDLE_TIME", 15)),
			ConMaxLifeTime: time.Duration(env.GetEnvAsInt64("DB_CONN_MAX_LIFTIME", 5)),
		}
	} else {
		cfg = config[0]
	}
	switch driver {
	case "postgres":
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.SSLMode,
		)
	case "mysql":
	case "sqlite3":
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConMaxLifeTime * time.Minute)
	db.SetConnMaxIdleTime(cfg.ConMaxIdleTime * time.Second)

	return db, nil
}
