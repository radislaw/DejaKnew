package handler

import (
	"errors"
	"net/http"

	"flashcards/api/internal/domain"
	"flashcards/api/internal/middleware"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cards *service.CardService
}

func NewCardHandler(cards *service.CardService) *CardHandler {
	return &CardHandler{cards: cards}
}

type cardBody struct {
	Front     string   `json:"front" binding:"required"`
	Back      string   `json:"back" binding:"required"`
	Topic     *string  `json:"topic"`
	Level     *string  `json:"level"`
	SourceURL *string  `json:"sourceUrl"`
	Tags      []string `json:"tags"`
}

func (b cardBody) toInput() service.CardInput {
	return service.CardInput{
		Front:     b.Front,
		Back:      b.Back,
		Topic:     b.Topic,
		Level:     b.Level,
		SourceURL: b.SourceURL,
		Tags:      b.Tags,
	}
}

func (h *CardHandler) Create(c *gin.Context) {
	var body cardBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := h.cards.Create(c.GetString(middleware.ContextUserID), c.Param("id"), body.toInput())
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *CardHandler) Update(c *gin.Context) {
	var body cardBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := h.cards.Update(c.GetString(middleware.ContextUserID), c.Param("id"), body.toInput())
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) Delete(c *gin.Context) {
	err := h.cards.Delete(c.GetString(middleware.ContextUserID), c.Param("id"))
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func cardErrorStatus(err error) (int, bool) {
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
