package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func New() Config {
	var c Config
	err := envconfig.Process("WALLET", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}

type Config struct {
	AppLabel string   `env:"APP_LABEL" envDefault:"wallet"`
	Host     string   `env:"HOST" envDefault:""`
	Port     string   `env:"PORT" envDefault:"8080"`
	LogLevel string   `env:"LOG_LEVEL" envDefault:"debug"`
	Postgres Postgres `,prefix:"POSTGRES_"`
	Redis    Redis    `,prefix:"REDIS_"`
}

type Postgres struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Database string `env:"DATABASE" envDefault:"postgres"`
}

func (c Config) PostgresURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database,
	)
}

type Redis struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"6379"`
	Password string `env:"PASSWORD" envDefault:""`
}

func (c Config) RedisURL() string {
	return fmt.Sprintf(
		"%s:%d",
		c.Redis.Host, c.Redis.Port,
	)
}
