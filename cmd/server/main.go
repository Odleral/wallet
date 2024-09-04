package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"wallet/internal/bootstrap"
	"wallet/internal/config"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.New(ctx)

	fmt.Printf("%+v\n", cfg.PostgresURL())

	go func() {
		<-quitSignal

		cancel()
	}()

	app := bootstrap.NewApp(cfg)

	app.Run(ctx)
}
