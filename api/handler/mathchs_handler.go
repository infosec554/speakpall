package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
)

// @Summary      Get my match preferences
// @Tags         match-prefs
// @Produce      json
// @Success      200 {object} models.MatchPreferences
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/match-prefs [get]
// @Security     ApiKeyAuth
func (h Handler) GetMyMatchPrefs(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	mp, err := h.services.Matchs().GetMatchPrefs(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load match preferences", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, mp)
}

// @Summary      Update my match preferences (partial)
// @Tags         match-prefs
// @Accept       json
// @Produce      json
// @Param        data body models.UpdateMatchPrefsRequest true "Fields to update"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/match-prefs [patch]
// @Security     ApiKeyAuth
func (h Handler) PatchMyMatchPrefs(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req models.UpdateMatchPrefsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Matchs().UpsertMatchPrefs(ctx, userID.(string), req); err != nil {
		handleResponse(c, h.log, "failed to update match preferences", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "match preferences updated", http.StatusOK, nil)
}
