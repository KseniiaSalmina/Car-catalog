package main

import (
	"embed"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	_ "github.com/KseniiaSalmina/Car-catalog/docs"
	app "github.com/KseniiaSalmina/Car-catalog/internal"
	"github.com/KseniiaSalmina/Car-catalog/internal/config"
)

var cfg config.Application

//go:embed schema/20240414213620_init.sql
var embedMigrations embed.FS

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.Postgres.Migration = embedMigrations
}

// @title Car catalog
// @version 1.0.0
// @description microservice for storing cars info
// @host localhost:8088
// @BasePath /
func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}
