package handler

import (
	"net/http"

	"flashcards/api/internal/middleware"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

type StudyHandler struct {
	study *service.StudyService
}

func NewStudyHandler(study *service.StudyService) *StudyHandler {
	return &StudyHandler{study: study}
}

func (h *StudyHandler) GetQueue(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	cards, err := h.study.GetQueue(userID, c.Param("deckId"))
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cards)
}

func (h *StudyHandler) GetStats(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	stats, err := h.study.GetStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

type reviewBody struct {
	Quality int `json:"quality" binding:"min=0,max=3"`
}

func (h *StudyHandler) SubmitReview(c *gin.Context) {
	var body reviewBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString(middleware.ContextUserID)
	sp, err := h.study.SubmitReview(userID, c.Param("cardId"), body.Quality)
	if status, failed := cardErrorStatus(err); failed {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sp)
}
