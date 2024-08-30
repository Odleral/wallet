package bootstrap

import (
	"wallet/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initDB(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.PostgresURL())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
