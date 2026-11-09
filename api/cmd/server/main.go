package main

import (
	"log"

	"github.com/joho/godotenv"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/config"
	"flashcards/api/internal/handler"
	"flashcards/api/internal/i18n"
	"flashcards/api/internal/repository"
	"flashcards/api/internal/service"
)

func main() {
	_ = godotenv.Load()

	i18n.Init()

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
	oauthHandler := handler.NewOAuthHandler(authService, cfg.AppURL, cfg.FrontendURL, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.YandexClientID, cfg.YandexClientSecret)
	profileHandler := handler.NewProfileHandler(authService)

	deckRepo := repository.NewDeckRepo(db)
	deckService := service.NewDeckService(deckRepo)

	cardRepo := repository.NewCardRepo(db)
	cardService := service.NewCardService(cardRepo, deckRepo)
	cardHandler := handler.NewCardHandler(cardService)

	deckHandler := handler.NewDeckHandler(deckService, cardService)

	studyProgressRepo := repository.NewStudyProgressRepo(db)
	studyService := service.NewStudyService(studyProgressRepo, cardRepo, deckRepo)
	studyHandler := handler.NewStudyHandler(studyService)

	r := handler.Router(tokens, authHandler, oauthHandler, deckHandler, cardHandler, studyHandler, profileHandler, cfg.AllowedOrigins)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
