package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
)

// @Summary      Get my settings
// @Tags         settings
// @Produce      json
// @Success      200 {object} models.UserSettings
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/settings [get]
// @Security     ApiKeyAuth
func (h Handler) GetMySettings(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	s, err := h.services.Settings().GetUserSettings(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load settings", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, s)
}

// @Summary      Update my settings (partial)
// @Tags         settings
// @Accept       json
// @Produce      json
// @Param        data body models.UpdateSettingsRequest true "Fields to update"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/settings [patch]
// @Security     ApiKeyAuth
func (h Handler) PatchMySettings(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Settings().UpsertUserSettings(ctx, userID.(string), req); err != nil {
		handleResponse(c, h.log, "failed to update settings", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "settings updated", http.StatusOK, nil)
}
