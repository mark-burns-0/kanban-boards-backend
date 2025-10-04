package config

import (
	"log/slog"
	"os"
	"strconv"
)

type dbConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

type jwtConfig struct {
	AccessTokenSecret string
	AccessTokenTTL    string
}

type Config struct {
	LogLevel    int8
	BcryptPower int8
	ServerPort  string
	dbConfig
	jwtConfig
}

func init() {
	mustLoadEnv()
}

func NewConfig() *Config {
	logLevel, err := strconv.ParseInt(os.Getenv("LOG_LEVEL"), 10, 8)
	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
	bcryptPower, err := strconv.ParseInt(os.Getenv("BCRYPT_POWER"), 10, 8)
	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
	serverConfig := &Config{
		ServerPort:  os.Getenv("HTTP_PORT"),
		LogLevel:    int8(logLevel),
		BcryptPower: int8(bcryptPower),
	}
	serverConfig.dbConfig = dbConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	serverConfig.jwtConfig = jwtConfig{
		AccessTokenSecret: os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		AccessTokenTTL:    os.Getenv("JWT_ACCESS_TOKEN_TTL"),
	}
	return serverConfig
}

func (c *Config) SetLogger() {
	setNewDefaultLogger(slog.Level(c.LogLevel))
}
