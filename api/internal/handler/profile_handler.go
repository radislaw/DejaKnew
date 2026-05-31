package handler

import (
	"errors"
	"net/http"

	"flashcards/api/internal/domain"
	"flashcards/api/internal/middleware"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	auth *service.AuthService
}

func NewProfileHandler(auth *service.AuthService) *ProfileHandler {
	return &ProfileHandler{auth: auth}
}

type updateProfileBody struct {
	Email  string `json:"email"`
	Locale string `json:"locale"`
}

type changePasswordBody struct {
	Current string `json:"current" binding:"required,min=8"`
	New     string `json:"new" binding:"required,min=8"`
}

func (h *ProfileHandler) Update(c *gin.Context) {
	var body updateProfileBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString(middleware.ContextUserID)
	user, err := h.auth.UpdateProfile(userID, body.Email, body.Locale)
	if errors.Is(err, domain.ErrUserExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	var body changePasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString(middleware.ContextUserID)
	err := h.auth.ChangePassword(userID, body.Current, body.New)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}
	if errors.Is(err, domain.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"error": "password change not available for OAuth accounts"})
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	if err := h.auth.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.Status(http.StatusNoContent)
}
