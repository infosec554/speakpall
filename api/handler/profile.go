// api/handler/profile.go
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
)

// ------------------------------------------------------------
// GET /user/me
// ------------------------------------------------------------

// @Summary      Get my profile
// @Tags         profile
// @Produce      json
// @Success      200 {object} models.Profile
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me [get]
// @Security     ApiKeyAuth
func (h Handler) GetMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	prof, err := h.services.Profile().GetProfile(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load profile", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, prof)
}

// ------------------------------------------------------------
// PATCH /user/me
// ------------------------------------------------------------

// @Summary      Update my profile (partial)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        data body models.UpdateProfileRequest true "Fields to update (partial)"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me [patch]
// @Security     ApiKeyAuth
func (h Handler) PatchMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Profile().UpdateProfile(ctx, userID.(string), req); err != nil {
		handleResponse(c, h.log, "failed to update profile", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "profile updated", http.StatusOK, nil)
}

// ------------------------------------------------------------
// GET /user/me/interests
// ------------------------------------------------------------

// @Summary      Get my interests
// @Tags         profile
// @Produce      json
// @Success      200 {object} models.InterestsResponse
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/interests [get]
// @Security     ApiKeyAuth
func (h Handler) GetMyInterests(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	ids, err := h.services.Profile().GetUserInterests(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load interests", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.InterestsResponse{InterestIDs: ids})
}

// ------------------------------------------------------------
// PUT /user/me/interests  (replace all)
// ------------------------------------------------------------

// @Summary      Replace my interests
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        data body models.UpdateInterestsRequest true "Full list to set (replaces existing)"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /user/me/interests [put]
// @Security     ApiKeyAuth
func (h Handler) PutMyInterests(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	var req models.UpdateInterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Profile().ReplaceUserInterests(ctx, userID.(string), req.InterestIDs); err != nil {
		handleResponse(c, h.log, "failed to replace interests", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "interests updated", http.StatusOK, nil)
}

// ------------------------------------------------------------
// GET /user/me/settings
// ------------------------------------------------------------

// @Summary      Get my settings
// @Tags         profile
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

	s, err := h.services.Profile().GetUserSettings(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load settings", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, s)
}

// ------------------------------------------------------------
// PATCH /user/me/settings
// ------------------------------------------------------------

// @Summary      Update my settings (partial)
// @Tags         profile
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

	if err := h.services.Profile().UpsertUserSettings(ctx, userID.(string), req); err != nil {
		handleResponse(c, h.log, "failed to update settings", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "settings updated", http.StatusOK, nil)
}

// ------------------------------------------------------------
// GET /user/me/match-prefs
// ------------------------------------------------------------

// @Summary      Get my match preferences
// @Tags         profile
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

	mp, err := h.services.Profile().GetMatchPrefs(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load match preferences", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, mp)
}

// ------------------------------------------------------------
// PATCH /user/me/match-prefs
// ------------------------------------------------------------

// @Summary      Update my match preferences (partial)
// @Tags         profile
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

	if err := h.services.Profile().UpsertMatchPrefs(ctx, userID.(string), req); err != nil {
		handleResponse(c, h.log, "failed to update match preferences", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "match preferences updated", http.StatusOK, nil)
}
