package handler

import (
	"errors"
	"net/http"

	"flashcards/api/internal/domain"
	"flashcards/api/internal/middleware"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

// cardLister is the minimal card dependency for loading a deck's cards.
type cardLister interface {
	ListByDeck(deckID string) ([]domain.Card, error)
}

type DeckHandler struct {
	decks *service.DeckService
	cards cardLister
}

func NewDeckHandler(decks *service.DeckService, cards cardLister) *DeckHandler {
	return &DeckHandler{decks: decks, cards: cards}
}

type deckDetail struct {
	*domain.Deck
	Cards []domain.Card `json:"cards"`
}

type deckBody struct {
	Title       string   `json:"title" binding:"required,min=1,max=200"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"isPublic"`
}

func (b deckBody) toInput() service.DeckInput {
	return service.DeckInput{
		Title:       b.Title,
		Description: b.Description,
		Tags:        b.Tags,
		IsPublic:    b.IsPublic,
	}
}

func (h *DeckHandler) Create(c *gin.Context) {
	var body deckBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deck, err := h.decks.Create(c.GetString(middleware.ContextUserID), body.toInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create deck"})
		return
	}
	c.JSON(http.StatusCreated, deck)
}

func (h *DeckHandler) List(c *gin.Context) {
	decks, err := h.decks.List(c.GetString(middleware.ContextUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list decks"})
		return
	}
	c.JSON(http.StatusOK, decks)
}

func (h *DeckHandler) Get(c *gin.Context) {
	deck, err := h.decks.Get(c.GetString(middleware.ContextUserID), c.Param("id"))
	if status, failed := deckErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	cards, err := h.cards.ListByDeck(deck.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load cards"})
		return
	}
	if cards == nil {
		cards = []domain.Card{}
	}
	c.JSON(http.StatusOK, deckDetail{Deck: deck, Cards: cards})
}

func (h *DeckHandler) Update(c *gin.Context) {
	var body deckBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deck, err := h.decks.Update(c.GetString(middleware.ContextUserID), c.Param("id"), body.toInput())
	if status, failed := deckErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deck)
}

func (h *DeckHandler) Delete(c *gin.Context) {
	err := h.decks.Delete(c.GetString(middleware.ContextUserID), c.Param("id"))
	if status, failed := deckErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// deckErrorStatus maps a service error to an HTTP status.
// The bool is false when err is nil — nothing went wrong.
func deckErrorStatus(err error) (int, bool) {
	switch {
	case err == nil:
		return 0, false
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, true
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, true
	default:
		return http.StatusInternalServerError, true
	}
}
