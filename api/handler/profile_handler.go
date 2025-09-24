package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
)

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
