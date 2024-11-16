package config


import (
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Port     string
    Database DatabaseConfig
    Email    EmailConfig
    JWT      JWTConfig
    Frontend FrontendConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

type JWTConfig struct {
    Secret string
}

type FrontendConfig struct {
    URL string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    smtpPort, err := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
    if err != nil {
        return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
    }

    return &Config{
        Port: getEnvOrDefault("PORT", "8080"),
        Database: DatabaseConfig{
            Host:     getEnvOrDefault("DB_HOST", "localhost"),
            Port:     dbPort,
            User:     getEnvOrDefault("DB_USER", "postgres"),
            Password: os.Getenv("DB_PASSWORD"),
            DBName:   getEnvOrDefault("DB_NAME", "userdb"),
            SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
        },
        Email: EmailConfig{
            Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
            Port:     smtpPort,
            Username: os.Getenv("GMAIL_USER"),
            Password: os.Getenv("GMAIL_APP_PASSWORD"),
            From:     os.Getenv("GMAIL_USER"),
        },
        JWT: JWTConfig{
            Secret: getEnvOrDefault("JWT_SECRET", "your-default-secret"),
        },
        Frontend: FrontendConfig{
            URL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
        },
    }, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}