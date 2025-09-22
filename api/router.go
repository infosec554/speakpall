package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "speakpall/api/docs"
	"speakpall/api/handler"
	"speakpall/pkg/logger"
	"speakpall/service"
)

// @title           Auth + Profile API
// @version         1.0
// @description     Authentication, Profile, Settings, Match Prefs
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceManager, log logger.ILogger) *gin.Engine {
	h := handler.New(services, log)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// -------- AUTH --------
	auth := r.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/change-password", h.ChangePassword)
		auth.POST("/google", h.GoogleAuth)
		auth.POST("/logout", h.Logout)

		auth.POST("/request-password-reset", h.RequestPasswordReset)
		auth.POST("/reset-password", h.ResetPassword)
	}

	// -------- USER (JWT protected) --------
	user := r.Group("/user")
	user.Use(h.JWTMiddleware())
	{
		user.GET("/me", h.GetMe)
		user.PATCH("/me", h.PatchMe)

		user.GET("/me/interests", h.GetMyInterests)
		user.PUT("/me/interests", h.PutMyInterests)

		user.GET("/me/settings", h.GetMySettings)
		user.PATCH("/me/settings", h.PatchMySettings)

		user.GET("/me/match-prefs", h.GetMyMatchPrefs)
		user.PATCH("/me/match-prefs", h.PatchMyMatchPrefs)
	}

	return r
}
