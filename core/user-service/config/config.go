package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
    Port int `yaml:"port"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
    SSLMode  string `yaml:"sslmode"`
}

type LoggingConfig struct {
    Level string `yaml:"level"`
}

func LoadConfig(path string) (*Config, error) {
    config := &Config{}

    file, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    if err := yaml.Unmarshal(file, config); err != nil {
        return nil, err
    }

    // Переопределяем значения из переменных окружения, если они заданы
    if port := os.Getenv("SERVER_PORT"); port != "" {
        config.Server.Port = atoi(port)
    }
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        config.Database.Host = dbHost
    }
    // Аналогично для других параметров...

    return config, nil
}

func atoi(s string) int {
    i := 0
    for _, c := range s {
        i = i*10 + int(c-'0')
    }
    return i
}

