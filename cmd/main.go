package main

import (
	"atlas/internal/app"
	"atlas/internal/shared/config"
	"atlas/internal/shared/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Carga .env o, si no existe, .env.prod (sin sobreescribir variables ya seteadas).
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.prod")

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Base de Datos: %v", err)
	}

	application := app.New(cfg, db)

	if err := application.Run(); err != nil {
		log.Fatalf("Servidor: %v", err)
	}

}
