package main

import (
	"log"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/config"
	"flashcards/api/internal/handler"
	"flashcards/api/internal/repository"
	"flashcards/api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("migration: %v", err)
	}

	tokens := auth.NewTokenManager(cfg.JWTSecret)

	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo, tokens)
	authHandler := handler.NewAuthHandler(authService)

	deckRepo := repository.NewDeckRepo(db)
	deckService := service.NewDeckService(deckRepo)
	deckHandler := handler.NewDeckHandler(deckService)

	r := handler.Router(tokens, authHandler, deckHandler)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
