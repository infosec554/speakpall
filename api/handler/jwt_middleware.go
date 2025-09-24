package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"speakpall/pkg/jwt"
)


func (h Handler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
			handleResponse(c, h.log, "missing bearer token", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		token := strings.TrimSpace(authHeader[7:])
		if token == "" {
			handleResponse(c, h.log, "invalid bearer token", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		claims, err := jwt.ExtractClaims(token)
		if err != nil {
			handleResponse(c, h.log, "invalid or expired token", http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if t, _ := claims["typ"].(string); t != "access" {
			handleResponse(c, h.log, "invalid token type", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)
		if userID == "" {
			handleResponse(c, h.log, "invalid token claims", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
