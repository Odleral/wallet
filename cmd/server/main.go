package main

import (
	"fmt"
	"wallet/internal/bootstrap"
	"wallet/internal/config"
)

func main() {
	cfg := config.New()

	fmt.Println(cfg)

	app := bootstrap.NewApp(cfg)

	app.Run()
}
