package bootstrap

import (
	"go.uber.org/zap"
	"wallet/internal/config"
)

type App struct {
	config   config.Config
	teardown []func()
}

func (a App) Run() {

}

func NewApp(cfg config.Config) *App {
	//teardown := make([]func(), 0)

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	return &App{
		config: cfg,
	}
}
