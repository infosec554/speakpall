package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PostFriend godoc
// @Summary      Add a friend
// @Description  Send a friend request or add friend by ID
// @Tags         friends
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Friend ID"
// @Security     ApiKeyAuth
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /friends/{id} [post]
func (h Handler) PostFriend(c *gin.Context) {
	// JWT middleware orqali user_id olinadi
	uid, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	userID := uid.(string)
	friendID := c.Param("id") // path param

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Friend().AddFriend(ctx, userID, friendID); err != nil {
		handleResponse(c, h.log, "failed to add friend", http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(c, h.log, "friend added", http.StatusOK, nil)
}

// DeleteFriend godoc
// @Summary      Delete friend
// @Description  Remove friend by ID
// @Tags         friends
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Friend ID"
// @Security     ApiKeyAuth
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Router       /user/friends/{id} [delete]
func (h Handler) DeleteFriend(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	userID := uid.(string)
	friendID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.services.Friend().RemoveFriend(ctx, userID, friendID); err != nil {
		handleResponse(c, h.log, "failed to delete friend", http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(c, h.log, "friend deleted", http.StatusOK, nil)
}

// GetFriends godoc
// @Summary      Get friends list
// @Description  Returns all friends of the logged-in user
// @Tags         friends
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} models.Response{data=[]models.User}
// @Failure      401 {object} models.Response
// @Router       /user/friends [get]
func (h Handler) GetFriends(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	userID := uid.(string)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	friends, err := h.services.Friend().ListFriends(ctx, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to get friends", http.StatusBadRequest, err.Error())
		return
	}

	handleResponse(c, h.log, "friends list", http.StatusOK, friends)
}
