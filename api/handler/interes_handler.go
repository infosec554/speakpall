package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
)

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

	ids, err := h.services.Interes().GetUserInterests(ctx, userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to load interests", http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.InterestsResponse{InterestIDs: ids})
}

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

	if err := h.services.Interes().ReplaceUserInterests(ctx, userID.(string), req.InterestIDs); err != nil {
		handleResponse(c, h.log, "failed to replace interests", http.StatusBadRequest, err.Error())
		return
	}
	handleResponse(c, h.log, "interests updated", http.StatusOK, nil)
}
