package main

import (
	"embed"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	app "github.com/KseniiaSalmina/Car-catalog/internal"
	"github.com/KseniiaSalmina/Car-catalog/internal/config"
)

var cfg config.Application

//go:embed schema/*.sql
var embedMigrations embed.FS

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.Postgres.Migration = embedMigrations
}

func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}
