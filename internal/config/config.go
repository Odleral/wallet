package config

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"log"
)

func New(ctx context.Context) Config {
	var c Config
	err := envconfig.Process(ctx, &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}

type Config struct {
	AppLabel string `env:"APP_LABEL"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	LogLevel string `env:"LOG_LEVEL"`
	AMQP     string `env:"AMQP_URL"`

	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBDatabase string `env:"DB_DATABASE"`

	Redis  *Redis  `env:", prefix=REDIS_"`
	Jaeger *Jaeger `env:", prefix=JAEGER_"`

	TransactionSub *Subscribe `env:", prefix=TRANSACTION_SUB_"`
	TransactionPub *Publisher `env:", prefix=TRANSACTION_PUB_"`
}

type Redis struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Password string `env:"PASSWORD"`
}

func (c *Redis) RedisURL() string {
	return fmt.Sprintf(
		"%s:%s",
		c.Host, c.Port,
	)
}

type Subscribe struct {
	QueueName string `env:"QUEUE_NAME" `
}

type Publisher struct {
	Exchange   string `env:"EXCHANGE"`
	RoutingKey string `env:"ROUTING_KEY"`
	QueueName  string `env:"QUEUE_NAME"`
}

type Jaeger struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (c Config) PostgresURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBDatabase,
	)
}
