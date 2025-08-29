package config

import "os"

// AppConfig содержит всю конфигурацию приложения.
type AppConfig struct {
	Port     string
	DBFile   string
	Password string
}

var cfg AppConfig

// Load считывает конфигурацию из переменных окружения и устанавливает значения по умолчанию.
func Load() error {
	// godotenv.Load() должен быть вызван в main до этой функции.
	cfg = AppConfig{
		Port:     os.Getenv("TODO_PORT"),
		DBFile:   os.Getenv("TODO_DBFILE"),
		Password: os.Getenv("TODO_PASSWORD"),
	}

	if cfg.Port == "" {
		cfg.Port = "7540"
	}
	if cfg.DBFile == "" {
		cfg.DBFile = "scheduler.db"
	}
	return nil
}

// Get возвращает загруженную конфигурацию.
func Get() AppConfig {
	return cfg
}
